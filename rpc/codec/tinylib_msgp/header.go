package tinylib_msgp

type RequestHeader struct {
	Method string `msg:"method"`
	Seq    uint64 `msg:"seq"`
}

type ResponseHeader struct {
	Method string `msg:"method"`
	Seq    uint64 `msg:"seq"`
	Error  string `msg:"error"`
}
