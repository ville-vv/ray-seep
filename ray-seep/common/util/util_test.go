package util

import (
	"fmt"
	"testing"
)

func TestInt32ToBytes(t *testing.T) {
	tmpByte, _ := Int32ToBytes(6)
	fmt.Println(tmpByte)
	fmt.Println(BytesToInt32(tmpByte))
}
