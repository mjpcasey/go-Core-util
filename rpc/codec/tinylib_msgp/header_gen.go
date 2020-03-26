package tinylib_msgp

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *RequestHeader) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zxvk uint32
	zxvk, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zxvk > 0 {
		zxvk--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "method":
			z.Method, err = dc.ReadString()
			if err != nil {
				return
			}
		case "seq":
			z.Seq, err = dc.ReadUint64()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z RequestHeader) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "method"
	err = en.Append(0x82, 0xa6, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Method)
	if err != nil {
		return
	}
	// write "seq"
	err = en.Append(0xa3, 0x73, 0x65, 0x71)
	if err != nil {
		return err
	}
	err = en.WriteUint64(z.Seq)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z RequestHeader) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "method"
	o = append(o, 0x82, 0xa6, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64)
	o = msgp.AppendString(o, z.Method)
	// string "seq"
	o = append(o, 0xa3, 0x73, 0x65, 0x71)
	o = msgp.AppendUint64(o, z.Seq)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *RequestHeader) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zbzg uint32
	zbzg, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zbzg > 0 {
		zbzg--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "method":
			z.Method, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "seq":
			z.Seq, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z RequestHeader) Msgsize() (s int) {
	s = 1 + 7 + msgp.StringPrefixSize + len(z.Method) + 4 + msgp.Uint64Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *ResponseHeader) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zbai uint32
	zbai, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zbai > 0 {
		zbai--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "method":
			z.Method, err = dc.ReadString()
			if err != nil {
				return
			}
		case "seq":
			z.Seq, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "error":
			z.Error, err = dc.ReadString()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z ResponseHeader) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "method"
	err = en.Append(0x83, 0xa6, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Method)
	if err != nil {
		return
	}
	// write "seq"
	err = en.Append(0xa3, 0x73, 0x65, 0x71)
	if err != nil {
		return err
	}
	err = en.WriteUint64(z.Seq)
	if err != nil {
		return
	}
	// write "error"
	err = en.Append(0xa5, 0x65, 0x72, 0x72, 0x6f, 0x72)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Error)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z ResponseHeader) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "method"
	o = append(o, 0x83, 0xa6, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64)
	o = msgp.AppendString(o, z.Method)
	// string "seq"
	o = append(o, 0xa3, 0x73, 0x65, 0x71)
	o = msgp.AppendUint64(o, z.Seq)
	// string "error"
	o = append(o, 0xa5, 0x65, 0x72, 0x72, 0x6f, 0x72)
	o = msgp.AppendString(o, z.Error)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ResponseHeader) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zcmr uint32
	zcmr, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zcmr > 0 {
		zcmr--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "method":
			z.Method, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "seq":
			z.Seq, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "error":
			z.Error, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z ResponseHeader) Msgsize() (s int) {
	s = 1 + 7 + msgp.StringPrefixSize + len(z.Method) + 4 + msgp.Uint64Size + 6 + msgp.StringPrefixSize + len(z.Error)
	return
}
