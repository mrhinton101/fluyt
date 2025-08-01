package utils

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/mrhinton101/fluyt/internal/app/core/logger"
)

type Reader interface {
	Read(path string) ([]byte, error)
}

type LocalReader struct{}

func (LocalReader) Read(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func GetSubFilesByExt(dirRoot string, filetypes []string) (map[string]string, error) {
	resultFiles := make(map[string]string)
	err := filepath.WalkDir(dirRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logger.SLogger(logger.LogEntry{
				Level:     slog.LevelError,
				Err:       err,
				Component: "utils",
				Action:    "load directory",
				Msg:       fmt.Sprintf("error while loading %s", dirRoot),
				Target:    dirRoot,
			})
			return err
		}
		if d.IsDir() {
			return nil
		}
		lowerName := strings.ToLower(d.Name())
		for _, fileType := range filetypes {
			if strings.HasSuffix(lowerName, "."+strings.ToLower(fileType)) {
				deviceName := strings.TrimSuffix(d.Name(), filepath.Ext(d.Name()))
				if existing, found := resultFiles[deviceName]; found {
					err := fmt.Errorf("device file %q is duplicated at %q and %q", deviceName, existing, path)
					logger.SLogger(logger.LogEntry{
						Level:     slog.LevelError,
						Err:       err,
						Component: "utils",
						Action:    "load directory",
						Msg:       "two files with the same name exist. each file name must be unique",
						Target:    dirRoot,
					})
					return err
				}
				resultFiles[deviceName] = path
				break
			}
		}
		return nil
	})
	return resultFiles, err
}
