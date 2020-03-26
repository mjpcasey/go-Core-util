package msgpack

type RequestHeader struct {
	Method string
	Seq    uint64
}

type ResponseHeader struct {
	Method string
	Seq    uint64
	Error  string
}
