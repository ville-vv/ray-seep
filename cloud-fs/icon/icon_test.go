package icon

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestPrintIcon(t *testing.T) {
	f, err := os.Open("../icon/cloud-fs.ico")
	if err != nil {
		t.Logf("open error: %v", err)
	}
	defer f.Close()
	fSrc, err := ioutil.ReadAll(f)
	fmt.Print("[]byte{")
	n := 0
	for _, v := range fSrc {
		fmt.Print(v, ",")
		n++
		if n >= 40 {
			fmt.Println("")
			n = 0
		}
	}
	fmt.Println("}")
}
