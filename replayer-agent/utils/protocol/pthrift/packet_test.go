package pthrift

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/light-pan/sharingan/replayer-agent/utils/protocol/helper"
	"github.com/modern-go/parse"
	"github.com/modern-go/parse/model"
	"github.com/stretchr/testify/require"
)

func TestDecodeMessage(t *testing.T) {
	var testCase = []struct {
		raw       []byte
		expect    *messageBody
		shouldErr bool
	}{
		{
			raw: []byte{
				0x80, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x03, 0x53, 0x61, 0x79, 0x00,
				0x00, 0x00, 0x01,
			},
			expect: &messageBody{
				Type:       Call,
				Name:       "Say",
				SequenceID: 1,
			},
		},
		{
			raw: []byte{
				0x80, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x03, 0x53, 0x61, 0x79, 0x00,
				0x00, 0x00,
			},
			shouldErr: true,
		},
	}
	should := require.New(t)
	for idx, tc := range testCase {
		src, err := parse.NewSource(bytes.NewBuffer(tc.raw), 5)
		should.NoError(err)
		actual, err := decodeMessage(src)
		if tc.shouldErr {
			should.Error(err)
			continue
		}
		should.Equal(tc.expect, actual, "case #%d fail", idx)
	}
}

func TestDecodeMap(t *testing.T) {
	var testCase = []struct {
		raw       []byte
		expect    *mapVal
		shouldErr bool
	}{
		{
			raw: []byte{0x0a, 0x0b, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x03, 0x00, 0x00, 0x00, 0x09, 0x73, 0x6f, 0x6d, 0x65, 0x20, 0x65,
				0x6c, 0x73, 0x65, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a, 0x00,
				0x00, 0x00, 0x0d, 0x49, 0x20, 0x6c, 0x69, 0x6b, 0x65, 0x20, 0x73, 0x6f,
				0x63, 0x63, 0x65, 0x72, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
				0x00, 0x00, 0x00, 0x04, 0x68, 0x61, 0x68, 0x61},
			expect: &mapVal{
				KeyType: I64,
				ValType: String,
				Size:    3,
				Data: model.Map{
					3:  "some else",
					10: "I like soccer",
					1:  "haha",
				},
			},
		},
	}
	should := require.New(t)
	for idx, tc := range testCase {
		src, err := parse.NewSource(bytes.NewBuffer(tc.raw), 8)
		should.NoError(err)
		actual, err := decodeMap(src)
		if tc.shouldErr {
			should.Error(err)
			continue
		}
		should.Equal(tc.expect, actual, "case #%d fail", idx)
	}
}

func TestDecodeStruct(t *testing.T) {
	var testCase = []struct {
		raw       []byte
		expect    structVal
		shouldErr bool
	}{
		{
			raw: []byte{
				0x0b, 0x00, 0x00, 0x00, 0x00, 0x00, 0xa1, 0x48, 0x69,
				0x20, 0x63, 0x61, 0x69, 0x62, 0x69, 0x72, 0x64, 0x6d, 0x65, 0x2c, 0x20,
				0x49, 0x20, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x20, 0x79,
				0x6f, 0x75, 0x72, 0x20, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x3a,
				0x20, 0x49, 0x27, 0x6d, 0x20, 0x64, 0x65, 0x63, 0x6f, 0x64, 0x69, 0x6e,
				0x67, 0x20, 0x74, 0x68, 0x72, 0x69, 0x66, 0x74, 0x20, 0x70, 0x72, 0x6f,
				0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2c, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x74,
				0x68, 0x65, 0x20, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x20, 0x69, 0x73, 0x20,
				0x31, 0x30, 0x30, 0x2e, 0x38, 0x37, 0x35, 0x30, 0x30, 0x30, 0x2c, 0x20,
				0x61, 0x6e, 0x64, 0x20, 0x74, 0x68, 0x65, 0x20, 0x65, 0x78, 0x74, 0x72,
				0x61, 0x20, 0x69, 0x6e, 0x66, 0x6f, 0x20, 0x69, 0x73, 0x20, 0x6d, 0x61,
				0x70, 0x5b, 0x33, 0x3a, 0x73, 0x6f, 0x6d, 0x65, 0x20, 0x65, 0x6c, 0x73,
				0x65, 0x20, 0x31, 0x30, 0x3a, 0x49, 0x20, 0x6c, 0x69, 0x6b, 0x65, 0x20,
				0x73, 0x6f, 0x63, 0x63, 0x65, 0x72, 0x20, 0x31, 0x3a, 0x68, 0x61, 0x68,
				0x61, 0x5d, 0x0a, 0x00,
			},
			expect: structVal{
				0: "Hi caibirdme, I received your message: I'm decoding thrift protocol, and the score is 100.875000, and the extra info is map[3:some else 10:I like soccer 1:haha]\n",
			},
		},
	}
	should := require.New(t)
	for idx, tc := range testCase {
		src, err := parse.NewSource(bytes.NewBuffer(tc.raw), 10)
		should.NoError(err)
		actual, err := decodeStruct(src)
		if tc.shouldErr {
			should.Error(err)
			continue
		}
		should.Equal(tc.expect, actual, "case #%d fail", idx)
	}
}

func TestDecodeBinary(t *testing.T) {
	var testCase = []struct {
		raw       []byte
		expect    model.Map
		shouldErr bool
	}{
		{
			raw: []byte{
				0x80, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x03, 0x53, 0x61, 0x79, 0x00,
				0x00, 0x00, 0x01, 0x0c, 0x00, 0x01, 0x0b, 0x00, 0x01, 0x00, 0x00, 0x00,
				0x09, 0x63, 0x61, 0x69, 0x62, 0x69, 0x72, 0x64, 0x6d, 0x65, 0x0b, 0x00,
				0x02, 0x00, 0x00, 0x00, 0x1c, 0x49, 0x27, 0x6d, 0x20, 0x64, 0x65, 0x63,
				0x6f, 0x64, 0x69, 0x6e, 0x67, 0x20, 0x74, 0x68, 0x72, 0x69, 0x66, 0x74,
				0x20, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x00, 0x04, 0x00,
				0x02, 0x40, 0x59, 0x38, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0d, 0x00, 0x03,
				0x0a, 0x0b, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x03, 0x00, 0x00, 0x00, 0x09, 0x73, 0x6f, 0x6d, 0x65, 0x20, 0x65,
				0x6c, 0x73, 0x65, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a, 0x00,
				0x00, 0x00, 0x0d, 0x49, 0x20, 0x6c, 0x69, 0x6b, 0x65, 0x20, 0x73, 0x6f,
				0x63, 0x63, 0x65, 0x72, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
				0x00, 0x00, 0x00, 0x04, 0x68, 0x61, 0x68, 0x61, 0x00,
			},
			expect: model.Map{
				"type":        "call",
				"name":        "Say",
				"sequence_id": 1,
				"param": model.Map{
					getFieldIDKey(1): model.Map{
						getFieldIDKey(1): "caibirdme",
						getFieldIDKey(2): `I'm decoding thrift protocol`,
					},
					getFieldIDKey(2): 100.875,
					getFieldIDKey(3): model.Map{
						"key_type": "i64",
						"val_type": "string",
						"size":     3,
						"data": model.Map{
							3:  "some else",
							10: "I like soccer",
							1:  "haha",
						},
					},
				},
			},
		},
		{
			raw: []byte{
				0x80, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0x03, 0x53, 0x61, 0x79, 0x00,
				0x00, 0x00, 0x01, 0x0b, 0x00, 0x00, 0x00, 0x00, 0x00, 0xa1, 0x48, 0x69,
				0x20, 0x63, 0x61, 0x69, 0x62, 0x69, 0x72, 0x64, 0x6d, 0x65, 0x2c, 0x20,
				0x49, 0x20, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x20, 0x79,
				0x6f, 0x75, 0x72, 0x20, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x3a,
				0x20, 0x49, 0x27, 0x6d, 0x20, 0x64, 0x65, 0x63, 0x6f, 0x64, 0x69, 0x6e,
				0x67, 0x20, 0x74, 0x68, 0x72, 0x69, 0x66, 0x74, 0x20, 0x70, 0x72, 0x6f,
				0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2c, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x74,
				0x68, 0x65, 0x20, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x20, 0x69, 0x73, 0x20,
				0x31, 0x30, 0x30, 0x2e, 0x38, 0x37, 0x35, 0x30, 0x30, 0x30, 0x2c, 0x20,
				0x61, 0x6e, 0x64, 0x20, 0x74, 0x68, 0x65, 0x20, 0x65, 0x78, 0x74, 0x72,
				0x61, 0x20, 0x69, 0x6e, 0x66, 0x6f, 0x20, 0x69, 0x73, 0x20, 0x6d, 0x61,
				0x70, 0x5b, 0x33, 0x3a, 0x73, 0x6f, 0x6d, 0x65, 0x20, 0x65, 0x6c, 0x73,
				0x65, 0x20, 0x31, 0x30, 0x3a, 0x49, 0x20, 0x6c, 0x69, 0x6b, 0x65, 0x20,
				0x73, 0x6f, 0x63, 0x63, 0x65, 0x72, 0x20, 0x31, 0x3a, 0x68, 0x61, 0x68,
				0x61, 0x5d, 0x0a, 0x00,
			},
			expect: model.Map{
				"type":        "reply",
				"name":        "Say",
				"sequence_id": 1,
				"param": model.Map{
					getFieldIDKey(0): "Hi caibirdme, I received your message: I'm decoding thrift protocol, and the score is 100.875000, and the extra info is map[3:some else 10:I like soccer 1:haha]\n",
				},
			},
		},
	}
	should := require.New(t)
	for idx, tc := range testCase {
		actual, err := DecodeBinary(tc.raw)
		if tc.shouldErr {
			should.Error(err)
			continue
		}
		should.Equal(tc.expect, actual, "case #%d fail", idx)
	}
}

func TestDecodeCompact(t *testing.T) {
	var testCase = []struct {
		raw    []byte
		expect model.Map
	}{
		{
			raw: []byte{
				0x82, 0x21, 0x01, 0x03, 0x53, 0x61, 0x79, 0x1c, 0x18, 0x09, 0x63, 0x61,
				0x69, 0x62, 0x69, 0x72, 0x64, 0x6d, 0x65, 0x18, 0x1c, 0x49, 0x27, 0x6d,
				0x20, 0x64, 0x65, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x20, 0x74, 0x68,
				0x72, 0x69, 0x66, 0x74, 0x20, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f,
				0x6c, 0x29, 0x35, 0x80, 0x01, 0x09, 0xf6, 0x01, 0x00, 0x17, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x38, 0x59, 0x40, 0x1b, 0x03, 0x68, 0x02, 0x04, 0x68,
				0x61, 0x68, 0x61, 0x06, 0x09, 0x73, 0x6f, 0x6d, 0x65, 0x20, 0x65, 0x6c,
				0x73, 0x65, 0x14, 0x0d, 0x49, 0x20, 0x6c, 0x69, 0x6b, 0x65, 0x20, 0x73,
				0x6f, 0x63, 0x63, 0x65, 0x72, 0x00,
			},
			expect: model.Map{
				"type":        "call",
				"name":        "Say",
				"sequence_id": 1,
				"param": model.Map{
					getFieldIDKey(1): model.Map{
						getFieldIDKey(1): "caibirdme",
						getFieldIDKey(2): "I'm decoding thrift protocol",
						getFieldIDKey(4): model.Map{
							"val_type": "i32",
							"data":     model.List{64, -5, 123},
						},
					},
					getFieldIDKey(2): 100.875,
					getFieldIDKey(3): model.Map{
						"data": model.Map{
							1:  "haha",
							3:  "some else",
							10: "I like soccer",
						},
						"key_type": "i64",
						"val_type": "string",
						"size":     3,
					},
				},
			},
		},
	}
	should := require.New(t)
	for idx, tc := range testCase {
		actual, err := DecodeCompact(tc.raw)
		should.NoError(err)
		should.Equal(tc.expect, actual, "case #%d fail", idx)
	}
}

func TestIntToBytes(t *testing.T) {
	type args struct {
		n int
		b byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "1",
			args:    args{n: 10, b: 4},
			want:    []byte{0, 0, 0, 0x0a},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := helper.IntToBytes(tt.args.n, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("IntToBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntToBytes() got = %v, want %v", got, tt.want)
			}
		})
	}
}
