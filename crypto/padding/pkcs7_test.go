package padding

import "testing"
import "github.com/JoveYu/zgo/log"

func TestPkcs7(t *testing.T) {
	log.Install("stdout")

	blocksize := 8
	alldata := [][]byte{
		{0xAB, 0xCD, 0xEF},
		{0xAB, 0xCD, 0xEF, 0xEF, 0xEF, 0xEF, 0xEF},
		{0xAB, 0xCD, 0xEF, 0xEF, 0xEF, 0xEF, 0xEF, 0xEF},
	}

	for _, data := range alldata {
		pad := Pkcs7Padding{blocksize}
		padded, err := pad.Pad(data)
		if err != nil {
			log.Error(err)
			return
		}
		unpadded, err := pad.UnPad(padded)
		if err != nil {
			log.Error(err)
			return
		}

		log.Debug("before: %X", data)
		log.Debug("pad:    %X", padded)
		log.Debug("unpad:  %X", unpadded)

	}

}
