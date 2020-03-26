package tinylib_msgp

import (
	"bufio"
	"errors"
	"io"
	"net/rpc"
)

type clientCodec struct {
	req RequestHeader
	enc *Encoder
	w   *bufio.Writer

	resp ResponseHeader
	dec  *Decoder
	c    io.Closer
}

// NewClientCodec returns a new rpc.Client.
//
// A ClientCodec implements writing of RPC requests and reading of RPC
// responses for the client side of an RPC session. The client calls
// WriteRequest to write a request to the connection and calls
// ReadResponseHeader and ReadResponseBody in pairs to read responses. The
// client calls Close when finished with the connection. ReadResponseBody
// may be called with a nil argument to force the body of the response to
// be read and then discarded.
func NewClientCodec(rwc io.ReadWriteCloser) rpc.ClientCodec {
	w := bufio.NewWriterSize(rwc, defaultBufferSize)
	r := bufio.NewReaderSize(rwc, defaultBufferSize)
	return &clientCodec{
		enc: NewEncoder(w),
		w:   w,
		dec: NewDecoder(r),
		c:   rwc,
	}
}

func (c *clientCodec) WriteRequest(req *rpc.Request, body interface{}) (err error) {
	c.req.Method = req.ServiceMethod
	c.req.Seq = req.Seq

	if err = c.enc.Encode(&c.req); err != nil {
		return err
	}

	b, ok := body.(MsgpCodeType)
	if !ok {
		return errors.New("request body invalid")
	}

	if err = c.enc.Encode(b); err != nil {
		return err
	}

	return c.w.Flush()
}

func (c *clientCodec) ReadResponseHeader(resp *rpc.Response) error {
	c.resp = ResponseHeader{}
	if err := c.dec.Decode(&c.resp); err != nil {
		return err
	}

	resp.ServiceMethod = c.resp.Method
	resp.Seq = c.resp.Seq
	resp.Error = c.resp.Error
	return nil
}

func (c *clientCodec) ReadResponseBody(body interface{}) (err error) {
	if body == nil {
		return nil
	}

	b, ok := body.(MsgpCodeType)
	if !ok {
		return errors.New("response body invalid")
	}

	return c.dec.Decode(b)
}

func (c *clientCodec) Close() error { return c.c.Close() }
