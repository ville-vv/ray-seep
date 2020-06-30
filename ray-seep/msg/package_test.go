package msg

import (
	"fmt"
	"testing"
)

func TestPackerManager01_Pack(t *testing.T) {
	pkm01 := &packerManager01{}
	data, err := pkm01.Pack(&Package{
		Cmd:  23345,
		Body: []byte("PASS: TestInt32ToBytes (0.00s)"),
	})
	fmt.Println(err)
	var pkgResp Package
	err = pkm01.UnPack(data, &pkgResp)
	fmt.Println(pkgResp.Cmd, string(pkgResp.Body), err)
}

func TestPackerManagerJson_Pack(t *testing.T) {
	pkm01 := &packerManagerJson{}
	data, err := pkm01.Pack(&Package{
		Cmd:  23345,
		Body: []byte("PASS: TestInt32ToBytes (0.00s)"),
	})
	fmt.Println(err)
	var pkgResp Package
	err = pkm01.UnPack(data, &pkgResp)
	fmt.Println(pkgResp.Cmd, string(pkgResp.Body), err)
}

func BenchmarkPackerManager01_Pack(b *testing.B) {
	// BenchmarkPackerManager01_Pack-4          3827990               280 ns/op             188 B/op          6 allocs/op
	pkm01 := &packerManager01{}
	for i := 0; i < b.N; i++ {
		_, _ = pkm01.Pack(&Package{
			Cmd:  23345,
			Body: []byte("PASS: TestInt32ToBytes (0.00s)"),
		})
	}
}

func BenchmarkPackerManagerJson_Pack(b *testing.B) {
	// BenchmarkPackerManagerJson_Pack-4        2484488               477 ns/op             296 B/op          6 allocs/op
	pkm01 := &packerManagerJson{}
	for i := 0; i < b.N; i++ {
		_, _ = pkm01.Pack(&Package{
			Cmd:  23345,
			Body: []byte("PASS: TestInt32ToBytes (0.00s)"),
		})
	}
}
