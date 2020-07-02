package f_system

import (
	"io"
	"net/http"
	"os"
	"strings"
)

type fileResponse struct {
	fileName string
	wt       http.ResponseWriter
	isFirst  bool
}

func newFileResponse(wt http.ResponseWriter, fileName string) *fileResponse {
	strList := strings.Split(fileName, "/")
	fileName = strList[len(strList)-1]
	return &fileResponse{wt: wt, fileName: fileName}
}

func (sel *fileResponse) Write(p []byte) (n int, err error) {
	if !sel.isFirst && len(p) > 10 {
		sel.setHeader(HeaderType(p[:10]))
		sel.isFirst = true
	}
	return sel.wt.Write(p)
}
func (sel *fileResponse) setHeader(tp string) {
	switch tp {
	case "download":
		sel.wt.Header().Add("Content-Type", "application/octet-stream")
		sel.wt.Header().Add("Content-Disposition", "attachment; filename=\""+sel.fileName+"\"")
	case "css":
		sel.wt.Header().Add("Content-Type", "text/css")
	case "js":
		sel.wt.Header().Add("Context-Type", "application/javascript")
	}
}

type TextFile struct {
}

func (df *TextFile) writeFile(fileName string, w io.Writer) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(w, file)
	return err
}

func (df *TextFile) Display(root string, fileName string, rsp http.ResponseWriter, req *http.Request) error {
	return df.writeFile(fileName, newFileResponse(rsp, fileName))
}
