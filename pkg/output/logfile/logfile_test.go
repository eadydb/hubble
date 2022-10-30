package logfile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/eadydb/hubble/pkg/testutil"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		description  string
		path         []string
		expectedName string
	}{
		{
			description:  "create file",
			path:         []string{"logs.txt"},
			expectedName: "logs.txt",
		},
		{
			description:  "create file in folder",
			path:         []string{"build", "logs.txt"},
			expectedName: filepath.Join("build", "logs.txt"),
		},
		{
			description:  "escape name",
			path:         []string{"a*name.txt"},
			expectedName: "a-name.txt",
		},
	}

	for _, test := range tests {
		testutil.Run(t, test.description, func(t *testutil.T) {
			file, err := Create(test.path...)
			defer func() {
				file.Close()
				os.Remove(file.Name())
			}()

			t.CheckNoError(err)
			t.CheckDeepEqual(filepath.Join(os.TempDir(), "hubble", test.expectedName), file.Name())
		})
	}
}

func TestEscape(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{name: "img", expected: "img"},
		{name: "log.txt", expected: "log.txt"},
		{name: "project/img", expected: "project-img"},
		{name: "gcr.io/project/img", expected: "gcr.io-project-img"},
	}
	for _, test := range tests {
		testutil.Run(t, test.name, func(t *testutil.T) {
			escaped := escape(test.name)

			t.CheckDeepEqual(test.expected, escaped)
		})
	}
}
