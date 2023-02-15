package yamlfmt

type Engine interface {
	FormatContent(content []byte) ([]byte, error)
	Format(paths []string) error
	Lint(paths []string) (string, error)
	DryRun(paths []string) (string, error)
}
