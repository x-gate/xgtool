package internal

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

// ErrInvalidFlag is returned when the flag is invalid.
var ErrInvalidFlag = errors.New("invalid flag")

// Decode from Run-Length Encoding.
func Decode(encoded []byte) (decoded []byte, err error) {
	buf := bytes.NewBuffer(encoded)

	for {
		var fb byte // first byte
		if fb, err = buf.ReadByte(); err != nil && errors.Is(err, io.EOF) {
			err = nil
			break
		} else if err != nil {
			return
		}

		var data *byte
		if data, err = readData(fb, buf); err != nil {
			return
		}

		var bc int32 // bytes count
		if bc, err = readBytesCnt(fb, buf); err != nil {
			return
		}

		var app []byte
		if app, err = makeAppend(data, bc, buf); err != nil {
			return
		}

		decoded = append(decoded, app...)
	}

	return
}

func readData(fb byte, buf *bytes.Buffer) (*byte, error) {
	switch fb & 0xf0 {
	case 0x00, 0x10, 0x20:
		return nil, nil
	case 0x80, 0x90, 0xa0:
		b, err := buf.ReadByte()
		return &b, err
	case 0xc0, 0xd0, 0xe0:
		var b byte
		return &b, nil
	default:
		return nil, fmt.Errorf("%w: %x", ErrInvalidFlag, fb)
	}
}

func readBytesCnt(fb byte, buf *bytes.Buffer) (int32, error) {
	switch fb & 0xf0 {
	case 0x00, 0x80, 0xc0:
		return int32(fb & 0x0f), nil
	case 0x10, 0x90, 0xd0:
		b, err := buf.ReadByte()

		return int32(fb&0x0f)<<8 + int32(b), err
	case 0x20, 0xa0, 0xe0:
		var b [2]byte
		var err error
		b[0], err = buf.ReadByte()
		b[1], err = buf.ReadByte()

		return int32(fb&0x0f)<<16 + int32(b[0])<<8 + int32(b[1]), err
	default:
		return 0, fmt.Errorf("%w: %x", ErrInvalidFlag, fb)
	}
}

func makeAppend(data *byte, cnt int32, buf *bytes.Buffer) (app []byte, err error) {
	app = make([]byte, cnt)

	if data == nil {
		if _, err = io.ReadFull(buf, app); err != nil && !errors.Is(err, io.EOF) {
			return
		}
	} else {
		app = bytes.Repeat([]byte{*data}, int(cnt))
	}

	return app, nil
}
