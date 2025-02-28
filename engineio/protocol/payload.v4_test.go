package protocol

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadPayloadV4(t *testing.T) {
	var opts []testoption

	runWithOptions := map[string]func(opts ...testoption) func(string, PayloadV4, bool, error) func(*testing.T){
		".Decode": func(opts ...testoption) func(string, PayloadV4, bool, error) func(*testing.T) {
			return func(data string, want PayloadV4, isXHR2 bool, xerr error) func(*testing.T) {
				return func(t *testing.T) {
					for _, opt := range opts {
						opt(t)
					}

					var have PayloadV4
					var dec = NewPayloadDecoderV4(strings.NewReader(data))
					dec.IsXHR2 = isXHR2
					var err = dec.Decode(&have)

					assert.ErrorIs(t, err, xerr)
					assert.Equal(t, want, have)
				}
			}
		},
		".ReadPayload": func(opts ...testoption) func(string, PayloadV4, bool, error) func(*testing.T) {
			return func(data string, want PayloadV4, isXHR2 bool, xerr error) func(*testing.T) {
				return func(t *testing.T) {
					for _, opt := range opts {
						opt(t)
					}

					var have PayloadV4
					var err = NewPayloadDecoderV4.From(strings.NewReader(data)).ReadPayload(&have)

					assert.ErrorIs(t, err, xerr)
					assert.Equal(t, want, have)
				}
			}
		},
	}

	spec := map[string]func() (string, PayloadV4, bool, error){
		"Without Binary": func() (string, PayloadV4, bool, error) {
			isBinary, isXHR2 := false, false
			data := "4hello\x1e4€"
			want := PayloadV4{
				{PacketV3{Packet{T: MessagePacket, D: "hello"}, isBinary}},
				{PacketV3{Packet{T: MessagePacket, D: "€"}, isBinary}},
			}
			return data, want, isXHR2, nil
		},
		"With Binary": func() (string, PayloadV4, bool, error) {
			isBinary, isXHR2 := true, false
			data := "4€\x1ebAQIDBA=="
			want := PayloadV4{
				{PacketV3{Packet{T: MessagePacket, D: "€"}, false}},
				{PacketV3{Packet{T: BinaryPacket, D: bytes.NewBuffer([]byte{0x01, 0x02, 0x03, 0x04})}, isBinary}},
			}
			return data, want, isXHR2, nil
		},
	}

	for name, testing := range spec {
		for suffix, runWithOption := range runWithOptions {
			t.Run(name+suffix, runWithOption(opts...)(testing()))
		}
	}
}

func TestWritePayloadV4(t *testing.T) {
	var opts []testoption

	runWithOptions := map[string]func(opts ...testoption) func(PayloadV4, string, bool, error) func(*testing.T){
		".Encode": func(opts ...testoption) func(PayloadV4, string, bool, error) func(*testing.T) {
			return func(data PayloadV4, want string, isXHR2 bool, xerr error) func(*testing.T) {
				return func(t *testing.T) {
					for _, opt := range opts {
						opt(t)
					}

					var have = new(bytes.Buffer)
					var enc = NewPayloadEncoderV4(have)
					var err = enc.Encode(data)

					assert.ErrorIs(t, err, xerr)
					assert.Equal(t, want, have.String())
				}
			}
		},
		".WritePayload": func(opts ...testoption) func(PayloadV4, string, bool, error) func(*testing.T) {
			return func(data PayloadV4, want string, isXHR2 bool, xerr error) func(*testing.T) {
				return func(t *testing.T) {
					for _, opt := range opts {
						opt(t)
					}

					var have = new(bytes.Buffer)
					var err = NewPayloadEncoderV4.To(have).WritePayload(data)

					assert.ErrorIs(t, err, xerr)
					assert.Equal(t, want, have.String())
				}
			}
		},
	}

	spec := map[string]func() (PayloadV4, string, bool, error){
		"Without Binary": func() (PayloadV4, string, bool, error) {
			isBinary, isXHR2 := false, false
			want := "4hello\x1e4€"
			data := PayloadV4{
				{PacketV3{Packet{T: MessagePacket, D: "hello"}, isBinary}},
				{PacketV3{Packet{T: MessagePacket, D: "€"}, isBinary}},
			}
			return data, want, isXHR2, nil
		},
		"With Binary": func() (PayloadV4, string, bool, error) {
			isBinary, isXHR2 := true, false
			want := "4€\x1ebAQIDBA=="
			data := PayloadV4{
				{PacketV3{Packet{T: MessagePacket, D: "€"}, false}},
				{PacketV3{Packet{T: BinaryPacket, D: bytes.NewBuffer([]byte{0x01, 0x02, 0x03, 0x04})}, isBinary}},
			}
			return data, want, isXHR2, nil
		},
	}

	for name, testing := range spec {
		for suffix, runWithOption := range runWithOptions {
			t.Run(name+suffix, runWithOption(opts...)(testing()))
		}
	}
}
