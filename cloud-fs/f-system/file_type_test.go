package f_system

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGetFileType(t *testing.T) {
	f, err := os.Open("bootstrap.css")
	if err != nil {
		t.Logf("open error: %v", err)
	}
	defer f.Close()

	fSrc, err := ioutil.ReadAll(f)
	t.Log(GetFileType(fSrc[:10]))
}

func TestGetFileType3(t *testing.T) {
	f, err := os.Open("css/bootstrap.css")
	if err != nil {
		t.Logf("open error: %v", err)
	}
	defer f.Close()

	fSrc, err := ioutil.ReadAll(f)
	t.Log(bytesToHexString(fSrc[:10]))

	f2, err := os.Open("css/style.css")
	if err != nil {
		t.Logf("open error: %v", err)
	}
	defer f2.Close()

	fSrc2, err := ioutil.ReadAll(f2)
	t.Log(bytesToHexString(fSrc2[:10]))

	f3, err := os.Open("css/chocolat.css")
	if err != nil {
		t.Logf("open error: %v", err)
	}
	defer f2.Close()

	fSrc3, err := ioutil.ReadAll(f3)
	t.Log(bytesToHexString(fSrc3[:10]))
}
