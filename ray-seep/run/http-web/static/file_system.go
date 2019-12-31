package static

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var htmlPageTemp = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Ray Seep Test</title>
</head>
<body>

<div style="margin: 0px 100px 100px 100px">
	<div style="">
		<h1>简单文件服务器</h1>
		<div style="font-weight: bolder;color:#00cb1a">
			<span>你当前的设备:&nbsp&nbsp&nbsp{{.User_Agent}}</span>
		</div>
		<p>如果没有文件，请添加文件到你的服务目录下</p>
		<h4>目录文件：</h4>
	</div>
	<div style="background-color: #c6d5e9">
		{{.Context}}
	<div>
</div>

</body>
</html>
`

func UseTemplate(w io.Writer, tmpl string, data interface{}) error {
	tp, err := template.New("tmpl").Parse(tmpl)
	if err != nil {
		return err
	}
	return tp.Execute(w, data)
}

type FileSystem struct {
	root string
}

func NewFileSystem(root string) *FileSystem {
	if strings.Trim(root, " ") == "" {
		root = "./"
	}
	return &FileSystem{root: filepath.Join("", root)}
}
func (f *FileSystem) displayContent(w io.Writer, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, file)
	return err
}
func (f *FileSystem) displayDir(w io.Writer, path string) error {
	fileList, err := PathEasyWolk(path)
	if err != nil {
		return err
	}
	_, _ = w.Write([]byte(fmt.Sprintf("<div><a href=\"./\">./<a/></div>")))
	if path != f.root {
		_, _ = w.Write([]byte(fmt.Sprintf("<div><a href=\"../\">../<a/></div>")))
	}
	for _, file := range fileList {
		if path == file.Path() {
			continue
		}
		_, _ = w.Write([]byte(fmt.Sprintf("<div><a href=\"%s/\">%s<a/></div>", file.Name(), file.Name())))
	}
	return nil
}
func (f *FileSystem) DisplayFile(rsp http.ResponseWriter, req *http.Request) {
	var err error
	uriPath := filepath.Join(f.root, req.URL.Path)
	isDir, err := IsDir(uriPath)
	if err != nil {
		_, _ = rsp.Write([]byte(err.Error()))
		return
	}
	if !isDir {
		if err = f.displayContent(rsp, uriPath); err != nil {
			_, _ = rsp.Write([]byte(err.Error()))
		}
		return
	}
	buf := bytes.NewBufferString("")
	if err = f.displayDir(buf, uriPath); err != nil {
		_, _ = rsp.Write([]byte(err.Error()))
	}

	if err = UseTemplate(rsp, htmlPageTemp, map[string]string{
		"Context":    buf.String(),
		"RemoterIp":  req.RemoteAddr,
		"User_Agent": req.Header.Get("User-Agent"),
	}); err != nil {
		_, _ = rsp.Write([]byte(err.Error()))
	}
}
