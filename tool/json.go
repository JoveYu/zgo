package tool

import (
	"fmt"
	"time"
)

type JSONTimeISO time.Time

func (t JSONTimeISO) MarshalJSON() ([]byte, error) {
	s := time.Time(t).Format(time.RFC3339)
	return []byte(fmt.Sprintf("\"%s\"", s)), nil
}

type JSONTimeTimestamp time.Time

func (t JSONTimeTimestamp) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", time.Time(t).Unix())), nil
}
