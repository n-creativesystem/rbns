package marshaler

import (
	"encoding/json"
	"io"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

var (
	DefaultContentType = "application/json"
)

type GatewayMarshaler struct{}

var _ runtime.Marshaler = (*GatewayMarshaler)(nil)

func (*GatewayMarshaler) ContentType(v interface{}) string {
	return DefaultContentType
}

func (m *GatewayMarshaler) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (m *GatewayMarshaler) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (m *GatewayMarshaler) NewDecoder(r io.Reader) runtime.Decoder {
	return runtime.DecoderFunc(
		func(v interface{}) error {
			buffer, err := io.ReadAll(r)
			if err != nil {
				return err
			}

			return m.Unmarshal(buffer, v)
		},
	)
}

func (m *GatewayMarshaler) NewEncoder(w io.Writer) runtime.Encoder {
	return json.NewEncoder(w)
}

func (m *GatewayMarshaler) Delimiter() []byte {
	return []byte("\n")
}
