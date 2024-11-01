package hotfix

import (
	"bufio"
	"bytes"
	"context"
	"strings"

	"github.com/google/yamlfmt"
)

type directiveKey string

var contextDirectivesKey directiveKey = "directives"

type Directive struct {
	line    int
	content string
}

func ContextWithDirectives(ctx context.Context, directives []Directive) context.Context {
	return context.WithValue(ctx, contextDirectivesKey, directives)
}

func DirectivesFromContext(ctx context.Context) []Directive {
	return ctx.Value(contextDirectivesKey).([]Directive)
}

func MakeFeatureStripDirectives(lineSepChar string) yamlfmt.Feature {
	return yamlfmt.Feature{
		Name:         "Strip Directives",
		BeforeAction: stripDirectivesFeature(lineSepChar),
		AfterAction:  restoreDirectivesFeature(lineSepChar),
	}
}

func stripDirectivesFeature(lineSepChar string) yamlfmt.FeatureFunc {
	return func(ctx context.Context, content []byte) (context.Context, []byte, error) {
		directives := []Directive{}
		reader := bytes.NewReader(content)
		scanner := bufio.NewScanner(reader)
		result := ""
		currLine := 1
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "%") {
				directives = append(directives, Directive{
					line:    currLine,
					content: line,
				})
			} else {
				result += line + lineSepChar
			}
			currLine++
		}
		return ContextWithDirectives(ctx, directives), []byte(result), nil
	}
}

func restoreDirectivesFeature(lineSepChar string) yamlfmt.FeatureFunc {
	return func(ctx context.Context, content []byte) (context.Context, []byte, error) {
		directives := DirectivesFromContext(ctx)
		directiveIdx := 0
		doneDirectives := directiveIdx == len(directives)
		reader := bytes.NewReader(content)
		scanner := bufio.NewScanner(reader)
		result := ""
		currLine := 1
		for scanner.Scan() {
			if !doneDirectives && currLine == directives[directiveIdx].line {
				result += directives[directiveIdx].content + lineSepChar
				currLine++
				directiveIdx++
				doneDirectives = directiveIdx == len(directives)
			}
			result += scanner.Text() + lineSepChar
			currLine++
		}
		// Edge case: There technically can be a directive as the final line. This would be
		// useless as far as I can tell so maybe yamlfmt should just remove it anyway LOL but
		// no we'll keep it.
		if !doneDirectives && currLine == directives[directiveIdx].line {
			result += directives[directiveIdx].content + lineSepChar
		}
		return ctx, []byte(result), nil
	}
}
