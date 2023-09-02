package tempfile

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/google/yamlfmt/internal/collections"
)

type Path struct {
	BasePath string
	FileName string
	IsDir    bool
	Content  []byte
}

func (p *Path) Create() error {
	file := p.fullPath()
	var err error
	if p.IsDir {
		err = os.Mkdir(file, os.ModePerm)
	} else {
		err = os.WriteFile(file, p.Content, os.ModePerm)
	}
	return err
}

func (p *Path) fullPath() string {
	return filepath.Join(p.BasePath, p.FileName)
}

type Paths []Path

func (ps Paths) CreateAll() error {
	errs := collections.Errors{}
	for _, path := range ps {
		errs = append(errs, path.Create())
	}
	return errs.Combine()
}

func ReplicateDirectory(dir string, newBase string) (Paths, error) {
	paths := Paths{}
	walkAllButCurrentDirectory := func(path string, info fs.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		// Skip the current directory (basically the . directory)
		if path == dir {
			return nil
		}
		content := []byte{}

		if !info.IsDir() {
			readContent, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			content = readContent
		}

		paths = append(paths, Path{
			BasePath: newBase,
			FileName: info.Name(),
			IsDir:    info.IsDir(),
			Content:  content,
		})
		return nil
	}
	err := filepath.Walk(dir, walkAllButCurrentDirectory)
	return paths, err
}
