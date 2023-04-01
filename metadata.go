package yamlfmt

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/google/yamlfmt/internal/collections"
)

const MetadataIdentifier = "!yamlfmt!"

type MetadataType string

const (
	MetadataIgnore MetadataType = "ignore"
)

func IsMetadataType(mdValueStr string) bool {
	mdTypes := collections.Set[MetadataType]{}
	mdTypes.Add(MetadataIgnore)
	return mdTypes.Contains(MetadataType(mdValueStr))
}

type Metadata struct {
	Type    MetadataType
	LineNum int
}

var (
	ErrMalformedMetadata    = errors.New("metadata: malformed metadata string")
	ErrUnrecognizedMetadata = errors.New("metadata: unrecognized metadata type")
)

func ReadMetadata(content []byte) (collections.Set[Metadata], collections.Errors) {
	metadata := collections.Set[Metadata]{}
	mdErrs := collections.Errors{}
	// This could be `\r\n` but it won't affect the outcome of this operation.
	contentLines := strings.Split(string(content), "\n")
	for i, line := range contentLines {
		mdidIndex := strings.Index(line, MetadataIdentifier)
		if mdidIndex == -1 {
			continue
		}
		mdStr := scanMetadata(line, mdidIndex)
		mdComponents := strings.Split(line, ":")
		if len(mdComponents) != 2 {
			mdErrs = append(mdErrs, fmt.Errorf("%w: %s", ErrMalformedMetadata, mdStr))
			continue
		}
		if IsMetadataType(mdComponents[1]) {
			metadata.Add(Metadata{LineNum: i + 1, Type: MetadataType(mdComponents[1])})
		} else {
			mdErrs = append(mdErrs, fmt.Errorf("%w: %s", ErrUnrecognizedMetadata, mdComponents[1]))
		}
	}
	return metadata, mdErrs
}

func scanMetadata(line string, index int) string {
	mdBytes := []byte{}
	i := index
	for i < len(line) && !unicode.IsSpace(rune(line[i])) {
		mdBytes = append(mdBytes, line[i])
		i++
	}
	return string(mdBytes)
}
