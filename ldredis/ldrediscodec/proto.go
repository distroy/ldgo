/*
 * Copyright (C) distroy
 */

package ldrediscodec

import (
	"github.com/golang/protobuf/proto"
	protov2 "google.golang.org/protobuf/proto"
)

type ProtoCodec[T proto.Message] struct {
	ProtoV1Codec[T]
}

type ProtoV1Codec[T proto.Message] struct {
	codecBase[T]
}

func (c ProtoV1Codec[T]) Marshal(v T) ([]byte, error) {
	return proto.Marshal(v)
}

func (c ProtoV1Codec[T]) Unmarshal(b []byte) (T, error) {
	var v T
	c.fillIfPointer(&v)
	return v, proto.Unmarshal(b, v)
}

type ProtoV2Codec[T proto.GeneratedMessage] struct {
	codecBase[T]
}

func (c ProtoV2Codec[T]) Marshal(v T) ([]byte, error) {
	vv := proto.MessageV2(v)
	return protov2.Marshal(vv)
}

func (c ProtoV2Codec[T]) Unmarshal(b []byte) (T, error) {
	var v T
	c.fillIfPointer(&v)
	vv := proto.MessageV2(v)
	return v, protov2.Unmarshal(b, vv)
}
