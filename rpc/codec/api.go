package codec

import (
	gencode_codec "gcore/rpc/codec/gencode"
	gob_codec "gcore/rpc/codec/gob"
	tinylib_msgpack_codec "gcore/rpc/codec/tinylib_msgp"
	protobuf_codec "github.com/mars9/codec"
	"net"
	"net/rpc"
)

const (
	CODEC_PROTOBUF = "protobuf"
	CODEC_GOB      = "gob"
	CODEC_MSGPACK  = "messagepack"
	CODEC_GENCODE  = "gencode"
)

func NewRpcClientCodec(codecType string, conn net.Conn) (codec rpc.ClientCodec) {
	switch codecType {
	case CODEC_PROTOBUF:
		codec = protobuf_codec.NewClientCodec(conn)
	case CODEC_MSGPACK:
		codec = tinylib_msgpack_codec.NewClientCodec(conn)
	case CODEC_GENCODE:
		codec = gencode_codec.NewClientCodec(conn)
	default:
		codec = gob_codec.NewClientCodec(conn)
	}

	return
}

func NewRpcServerCodec(codecType string, conn net.Conn) (codec rpc.ServerCodec) {
	switch codecType {
	case CODEC_GOB:
		codec = gob_codec.NewServerCodec(conn)
	case CODEC_PROTOBUF:
		codec = protobuf_codec.NewServerCodec(conn)
	case CODEC_MSGPACK:
		codec = tinylib_msgpack_codec.NewServerCodec(conn)
	case CODEC_GENCODE:
		codec = gencode_codec.NewServerCodec(conn)
	default:
		codec = gob_codec.NewServerCodec(conn)
	}

	return
}
