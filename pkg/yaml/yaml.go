package yaml

import (
	"bytes"
	"io"
	"reflect"

	yaml "gopkg.in/yaml.v3"
)

// UnmarshalStrict is like Unmarshal except that any fields that are found
// in the data that do not have corresponding struct members, or mapping
// keys that are duplicates, will result in an error.
// This is ensured by setting `KnownFields` true on `yaml.Decoder`.
func UnmarshalStrict(in []byte, out interface{}) error {
	b := bytes.NewReader(in)
	decoder := yaml.NewDecoder(b)
	decoder.KnownFields(true)
	// Ignore io.EOF which signals expected end of input.
	// This happens when input stream is empty or nil.
	if err := decoder.Decode(out); err != io.EOF {
		return err
	}
	return nil
}

// Unmarshal is wrapper around yaml.Unmarshal
func Unmarshal(in []byte, out interface{}) error {
	return yaml.Unmarshal(in, out)
}

// Marshal is same as yaml.Marshal except it creates a `yaml.Encoder` with
// indent space 2 for encoding.
func Marshal(in interface{}) (out []byte, err error) {
	var b bytes.Buffer
	encoder := yaml.NewEncoder(&b)
	encoder.SetIndent(2)
	if err := encoder.Encode(in); err != nil {
		return nil, err
	}
	if err = encoder.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// MarshalWithSeparator is same as Marshal except for slice or array types where each element is encoded individually and separated by "---".
func MarshalWithSeparator(in interface{}) (out []byte, err error) {
	var b bytes.Buffer
	encoder := yaml.NewEncoder(&b)
	encoder.SetIndent(2)

	switch reflect.TypeOf(in).Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		s := reflect.ValueOf(in)
		for i := 0; i < s.Len(); i++ {
			if err := encoder.Encode(s.Index(i).Interface()); err != nil {
				return nil, err
			}
		}
	default:
		if err := encoder.Encode(in); err != nil {
			return nil, err
		}
	}
	if err = encoder.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
