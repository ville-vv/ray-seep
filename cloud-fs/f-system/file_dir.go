package f_system

import (
	"bytes"
	"fmt"
	"net/http"
)

type DirFile struct{}

func (df *DirFile) Display(root string, path string, rsp http.ResponseWriter, req *http.Request) error {
	buf := bytes.NewBufferString("")

	fileList, err := PathEasyWolk(path)
	if err != nil {
		return err
	}
	_, _ = buf.Write([]byte(fmt.Sprintf("<div><a href=\"./\">./<a/></div>")))
	if path != root {
		_, _ = buf.Write([]byte(fmt.Sprintf("<div><a href=\"../\">../<a/></div>")))
	}
	for _, file := range fileList {
		if path == file.Path() {
			continue
		}

		isDir, err := IsDir(path + "/" + file.Name())
		if err != nil {
			continue
		}
		if isDir {
			_, _ = buf.Write([]byte(fmt.Sprintf("<div><a href=\"%s/\">%s<a/></div>", file.Name(), file.Name())))
			continue
		}
		_, _ = buf.Write([]byte(fmt.Sprintf("<div><a href=\"%s\">%s<a/></div>", file.Name(), file.Name())))
	}
	return UseTemplate(rsp, htmlPageTemp, map[string]string{
		"Context":    buf.String(),
		"RemoterIp":  req.RemoteAddr,
		"User_Agent": req.Header.Get("User-Agent"),
	})
}
