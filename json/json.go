package json

import (
	"io"

	// "github.com/bytedance/sonic"

	"github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

var j = jsoniter.Config{
	MarshalFloatWith6Digits:       false,
	EscapeHTML:                    true,
	SortMapKeys:                   true,
	UseNumber:                     true,
	DisallowUnknownFields:         false,
	TagKey:                        "",
	OnlyTaggedField:               false,
	ValidateJsonRawMessage:        true,
	ObjectFieldMustBeSimpleString: false,
	CaseSensitive:                 false,
}.Froze()

func init() {
	extra.RegisterFuzzyDecoders()
}

func Marshal(v interface{}) ([]byte, error) {
	return j.Marshal(v)
}

func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return j.MarshalIndent(v, prefix, indent)
}

func Unmarshal(data []byte, v interface{}) error {
	return j.Unmarshal(data, v)
}

func NewEncoder(w io.Writer) *jsoniter.Encoder {
	return j.NewEncoder(w)
}

func NewDecoder(r io.Reader) *jsoniter.Decoder {
	return j.NewDecoder(r)
}
