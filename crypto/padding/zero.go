package padding

import (
	"bytes"
	"errors"
)

type ZeroPadding struct {
	BlockSize int
}

func (p ZeroPadding) Pad(data []byte) ([]byte, error) {
	if p.BlockSize < 1 || p.BlockSize > 255 {
		return nil, errors.New("block size error")
	}

	length := p.BlockSize - len(data)%p.BlockSize

	return append(data, bytes.Repeat([]byte{byte(0)}, length)...), nil
}

func (p ZeroPadding) UnPad(data []byte) ([]byte, error) {
	datalen := len(data)

	var length int
	for length = 0; length <= datalen; length++ {
		if int(data[datalen-1-length]) != 0 {
			break
		}
	}

	return data[:datalen-length], nil
}
