package migrate_v1

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func (m Migration) New(name string) error {
	fname := fmt.Sprintf("%s_%s.go", time.Now().Format("20060102150405"), name)

	fpath := filepath.Join(m.OutputPath, fname)
	err := os.MkdirAll(m.OutputPath, os.ModePerm)
	if err != nil {
		return err
	}

	tmp, err := migrateTemplate.ReadFile("template.tmpl")
	if err != nil {
		return err
	}
	err = os.WriteFile(fpath, tmp, 0o644)
	if err != nil {
		return err
	}
	return nil
}
