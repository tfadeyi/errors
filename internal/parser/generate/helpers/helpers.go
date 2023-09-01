package helpers

import (
	"io"
	"os"
	"path/filepath"

	"github.com/juju/errors"
)

// Clean removes all the files
func Clean(files ...string) error {
	for _, file := range files {
		if _, err := os.Stat(file); !errors.Is(err, os.ErrNotExist) {
			// delete spec file
			err = os.RemoveAll(file)
			if err != nil {
				return errors.Annotatef(err, "could not delete existing file %q", file)
			}
		}
	}
	return nil
}

// Write the files to the writer, the caller is in charge of closing the writer
func Write(w io.Writer, files map[string][]byte) error {
	for _, body := range files {
		var err error
		// write to writer, this must be closed by the caller
		_, err = w.Write(body)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteToFile writes the files to the specified file paths. The function handles its writers
func WriteToFile(files map[string][]byte) error {
	for path, body := range files {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		w, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		_, err = w.Write(body)
		if err != nil {
			return err
		}
	}
	return nil
}
