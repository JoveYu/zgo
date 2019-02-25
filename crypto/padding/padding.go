package padding

type Padding interface {
	Pad(data []byte) ([]byte, error)
	UnPad(data []byte) ([]byte, error)
}
