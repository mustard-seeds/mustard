package proto_util

// proto3 syntax
// https://developers.google.com/protocol-buffers/docs/proto3

import (
	"mustard/base/log"
	"mustard/internal/github.com/golang/protobuf/proto"
)

func Serialize(pb proto.Message) ([]byte, error) {
	r,e := proto.Marshal(pb)
	if e != nil {
		log.Log.Error("MarShal failed")
	}
	return r,e
}
func Deserialize(s []byte, pb proto.Message) error {
	e := proto.Unmarshal(s, pb)
	if e != nil {
		log.Log.Error("UnMarShal failed")
	}
	return e
}
func FromProtoToString(pb proto.Message) string {
	return proto.MarshalTextString(pb)
}
func FromStringToProto(s string, pb proto.Message) error {
	e := proto.UnmarshalText(s, pb)
	if e != nil {
		log.Log.Error("UnMarShal failed")
	}
	return e
}
