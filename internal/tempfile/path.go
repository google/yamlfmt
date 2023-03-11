package tempfile

import (
	"os"
	"path/filepath"
)

type Path struct {
	BasePath string
	FileName string
	IsDir    bool
}

func (p *Path) Create() error {
	file := p.fullPath()
	var err error
	if p.IsDir {
		err = os.Mkdir(file, os.ModePerm)
	} else {
		err = os.WriteFile(file, []byte{}, os.ModePerm)
	}
	return err
}

func (p *Path) fullPath() string {
	return filepath.Join(p.BasePath, p.FileName)
}
