package yaml

import (
	"testing"

	"github.com/eadydb/hubble/pkg/testutil"
)

func TestMarshalWithSeparator(t *testing.T) {
	type Data struct {
		Foo string `yaml:"foo"`
	}

	tests := []struct {
		description string
		input       []Data
		expected    string
	}{
		{
			description: "single element slice",
			input: []Data{
				{Foo: "foo"},
			},
			expected: "foo: foo\n",
		},
		{
			description: "multi element slice",
			input: []Data{
				{Foo: "foo1"},
				{Foo: "foo2"},
			},
			expected: "foo: foo1\n---\nfoo: foo2\n",
		},
	}

	for _, test := range tests {
		testutil.Run(t, test.description, func(t *testutil.T) {
			output, err := MarshalWithSeparator(test.input)
			t.CheckNoError(err)
			t.CheckDeepEqual(string(output), test.expected)
		})
	}
}
