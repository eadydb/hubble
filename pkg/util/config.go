package util

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/eadydb/hubble/pkg/output/log"
	"github.com/spf13/afero"
)

var (
	// Fs is the underlying filesystem to use for reading hubble project files & configuration.  OS FS by default
	Fs    = afero.NewOsFs()
	stdin []byte
)

// ReadConfiguration reads a `hubble.yaml` configuration and
// returns its content.
func ReadConfiguration(filePath string) ([]byte, error) {
	switch {
	case filePath == "":
		return nil, errors.New("filename not specified")
	case filePath == "-":
		if len(stdin) == 0 {
			var err error
			stdin, err = io.ReadAll(os.Stdin)
			if err != nil {
				return []byte{}, err
			}
		}
		return stdin, nil
	case IsURL(filePath):
		return Download(filePath)
	default:
		if !filepath.IsAbs(filePath) {
			dir, err := os.Getwd()
			if err != nil {
				return []byte{}, err
			}
			filePath = filepath.Join(dir, filePath)
		}
		contents, err := afero.ReadFile(Fs, filePath)
		if err != nil {
			// If the config file is the default `hubble.yaml`,
			// then we also try to read `hubble.yml`.
			if filepath.Base(filePath) == "hubble.yaml" {
				log.Entry(context.TODO()).Infof("Could not open hubble.yaml: \"%s\"", err)
				log.Entry(context.TODO()).Info("Trying to read from hubble.yml instead")
				contents, errIgnored := afero.ReadFile(Fs, filepath.Join(filepath.Dir(filePath), "hubble.yml"))
				if errIgnored != nil {
					// Return original error because it's the one that matters
					return nil, err
				}

				return contents, nil
			}
		}

		return contents, err
	}
}

func ReadFile(filename string) ([]byte, error) {
	if !filepath.IsAbs(filename) {
		dir, err := os.Getwd()
		if err != nil {
			return []byte{}, err
		}
		filename = filepath.Join(dir, filename)
	}
	return afero.ReadFile(Fs, filename)
}
