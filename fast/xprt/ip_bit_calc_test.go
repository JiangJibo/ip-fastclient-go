package clientUtils

import "testing"

func TestReadVInt3(t *testing.T) {
	data := []byte{0, 1, 1, 1}
	value := ReadVInt3(data, 1)
	t.Log(value)
}

func TestWriteVint3(t *testing.T) {
	data := make([]byte, 3)
	WriteVInt3(data, 0, 65537)
	t.Log(data)
}
