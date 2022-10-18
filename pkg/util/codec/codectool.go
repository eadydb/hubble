package codec

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/invopop/yaml"
)

// NOTE: It's better to use functions of this package,
// since the vendor could be replaced at once if we need.

// MustUnmarshal wraps json.MustUnmarshal by panic instead of returning error.
// It converts yaml to json before decoding by leveraging Unmarshal.
// It panics if an error occurs.
func MustUnmarshal(data []byte, v interface{}) {
	err := Unmarshal(data, v)
	if err != nil {
		panic(fmt.Errorf("unmarshal %s to %#v failed: %v",
			data, v, err))
	}
}

// Unmarshal wraps json.Unmarshal.
// It will convert yaml to json before unmarshal.
// Since json is a subset of yaml, passing json through this method should be a no-op.
func Unmarshal(data []byte, v interface{}) error {
	data, err := yaml.YAMLToJSON(data)
	if err != nil {
		return fmt.Errorf("%s: convert yaml to json failed: %v", data, err)
	}
	json.Unmarshal(data, v)

	return json.Unmarshal(data, v)
}

// MustMarshalJSON wraps json.Marshal by panic instead of returning error.
func MustMarshalJSON(v interface{}) []byte {
	buff, err := MarshalJSON(v)
	if err != nil {
		panic(fmt.Errorf("marshal %#v to json failed: %v", v, err))
	}
	return buff
}

// MarshalJSON wraps json.Marshal.
func MarshalJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// UnmarshalJSON wraps json.Unmarshal.
func UnmarshalJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// MustUnmarshalJSON wraps json.Unmarshal.
// It panics if an error occurs.
func MustUnmarshalJSON(data []byte, v interface{}) {
	err := Unmarshal(data, v)
	if err != nil {
		panic(fmt.Errorf("unmarshal json %s to %#v failed: %v",
			data, v, err))
	}
}

// MustDecode decodes a json stream into a value of the given type.
// It converts yaml to json before decoding by leveraging Unmarshal.
// It panics if an error occurs.
func MustDecode(r io.Reader, v interface{}) {
	err := Decode(r, v)
	if err != nil {
		panic(err)
	}
}

// Decode decodes a json stream into a value of the given type.
// It converts yaml to json before decoding by leveraging Unmarshal.
func Decode(r io.Reader, v interface{}) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("read failed: %v", err)
	}
	return Unmarshal(data, v)
}

// MustDecodeJSON decodes a json stream into a value of the given type.
// It panics if an error occurs.
func MustDecodeJSON(r io.Reader, v interface{}) {
	err := DecodeJSON(r, v)
	if err != nil {
		panic(err)
	}
}

// DecodeJSON decodes a json stream into a value of the given type.
func DecodeJSON(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

// MustEncodeJSON encodes a value into a json stream.
// It panics if an error occurs.
func MustEncodeJSON(w io.Writer, v interface{}) {
	err := EncodeJSON(w, v)
	if err != nil {
		panic(err)
	}
}

// EncodeJSON encodes a value into a json stream.
func EncodeJSON(w io.Writer, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

// MustJSONToYAML converts a json stream into a yaml stream.
// It panics if an error occurs.
func MustJSONToYAML(in []byte) []byte {
	buff, err := JSONToYAML(in)
	if err != nil {
		panic(fmt.Errorf("json %s to yaml failed: %v", in, err))
	}
	return buff
}

// JSONToYAML converts a json stream into a yaml stream.
func JSONToYAML(in []byte) ([]byte, error) {
	return yaml.JSONToYAML(in)
}

// MustYAMLToJSON converts a json stream into a yaml stream.
// It panics if an error occurs.
func MustYAMLToJSON(in []byte) []byte {
	buff, err := YAMLToJSON(in)
	if err != nil {
		panic(fmt.Errorf("yaml %s to json failed: %v", in, err))
	}
	return buff
}

// YAMLToJSON converts a yaml stream into a json stream.
func YAMLToJSON(in []byte) ([]byte, error) {
	return yaml.YAMLToJSON(in)
}
