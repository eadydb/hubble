package kubernetes

import (
	"testing"

	"github.com/eadydb/hubble/pkg/testutil"
)

func TestSupportedKubernetesFormats(t *testing.T) {
	tests := []struct {
		description string
		in          string
		out         bool
	}{
		{
			description: "yaml",
			in:          "filename.yaml",
			out:         true,
		},
		{
			description: "yml",
			in:          "filename.yml",
			out:         true,
		},
		{
			description: "json",
			in:          "filename.json",
			out:         true,
		},
		{
			description: "txt",
			in:          "filename.txt",
			out:         false,
		},
	}
	for _, test := range tests {
		testutil.Run(t, test.description, func(t *testutil.T) {
			actual := HasKubernetesFileExtension(test.in)

			t.CheckDeepEqual(test.out, actual)
		})
	}
}
