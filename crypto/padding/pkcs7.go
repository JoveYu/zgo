package padding

import (
	"bytes"
	"errors"
)

type Pkcs7Padding struct {
	BlockSize int
}

func (p Pkcs7Padding) Pad(data []byte) ([]byte, error) {
	if p.BlockSize < 1 || p.BlockSize > 255 {
		return nil, errors.New("block size error")
	}

	length := p.BlockSize - len(data)%p.BlockSize

	return append(data, bytes.Repeat([]byte{byte(length)}, length)...), nil
}

func (p Pkcs7Padding) UnPad(data []byte) ([]byte, error) {
	datalen := len(data)
	length := int(data[datalen-1])

	if datalen%p.BlockSize != 0 {
		return nil, errors.New("not padded correctly")
	}

	if length > p.BlockSize || length <= 0 {
		return nil, errors.New("not padded correctly")
	}

	padding := data[datalen-length : datalen-1]
	for _, i := range padding {
		if int(i) != length {
			return nil, errors.New("not padded correctly")
		}
	}

	return data[:datalen-length], nil
}
