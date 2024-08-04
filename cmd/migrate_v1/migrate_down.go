package migrate_v1

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io/fs"
	"path/filepath"
	"strings"
	"time"
)

func (m Migration) Down(ctx context.Context, db *gorm.DB) error {
	var mt []MigrateTable
	err := db.Order("id desc").Scan(&mt).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	var batch int8
	migrated := make(map[string]bool)
	if len(mt) == 0 {
		log.Print("no migration available to rollback")
		return nil
	} else {
		batch = mt[0].Batch
		for _, x := range mt {
			if x.Batch == batch {
				migrated[x.Name] = true
			}
		}
	}

	var newMigration []MigrateTable
	tx := db.Begin()
	err = filepath.Walk(m.OutputPath, func(path string, info fs.FileInfo, err error) error {
		f := strings.Split(info.Name(), ".")
		if len(f) == 1 {
			return nil
		}
		if _, ok := migrated[f[0]]; !ok {
			return nil
		}

		fpath := filepath.Join(m.OutputPath, info.Name())
		execute, err := findMigration(fpath)
		err = execute.MigrateDown(ctx, tx)
		if err != nil {
			return err
		}
		nm := MigrateTable{
			Name:      f[0],
			Batch:     batch,
			CreatedAt: time.Now(),
		}
		newMigration = append(newMigration, nm)

		return nil
	})
	if err != nil {
		rollErr := tx.Rollback().Error
		return errors.Join(err, rollErr)
	}
	if len(newMigration) == 0 {
		tx.Commit()
		fmt.Println("nothing to migrate")
		return nil
	}

	err = tx.Create(&newMigration).Error
	if err != nil {
		rollErr := tx.Rollback().Error
		return errors.Join(err, rollErr)
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}
