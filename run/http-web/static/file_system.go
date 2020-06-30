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
func (f *FileSystem) writeFile(w io.Writer, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, file)
	return err
}
func (f *FileSystem) displayFile(fileName string, w http.ResponseWriter, req *http.Request) error {
	//strs := strings.Split(fileName, "/")
	//w.Header().Add("Content-Type", "application/octet-stream")
	//w.Header().Add("Content-Disposition", "attachment; filename=\""+strs[len(strs)-1]+"\"")
	return f.writeFile(NewFileResponse(w, fileName), fileName)
}
func (f *FileSystem) displayDir(path string, rsp http.ResponseWriter, req *http.Request) error {
	buf := bytes.NewBufferString("")
	fileList, err := PathEasyWolk(path)
	if err != nil {
		return err
	}
	_, _ = buf.Write([]byte(fmt.Sprintf("<div><a href=\"./\">./<a/></div>")))
	if path != f.root {
		_, _ = buf.Write([]byte(fmt.Sprintf("<div><a href=\"../\">../<a/></div>")))
	}
	for _, file := range fileList {
		if path == file.Path() {
			continue
		}
		_, _ = buf.Write([]byte(fmt.Sprintf("<div><a href=\"%s/\">%s<a/></div>", file.Name(), file.Name())))
	}
	return UseTemplate(rsp, htmlPageTemp, map[string]string{
		"Context":    buf.String(),
		"RemoterIp":  req.RemoteAddr,
		"User_Agent": req.Header.Get("User-Agent"),
	})
}
func (f *FileSystem) Display(rsp http.ResponseWriter, req *http.Request) {
	var err error
	uriPath := filepath.Join(f.root, req.URL.Path)
	isDir, err := IsDir(uriPath)
	if err != nil {
		_, _ = rsp.Write([]byte(err.Error()))
		return
	}
	if !isDir {
		if err = f.displayFile(uriPath, rsp, req); err != nil {
			_, _ = rsp.Write([]byte(err.Error()))
		}
		return
	}
	if err = f.displayDir(uriPath, rsp, req); err != nil {
		_, _ = rsp.Write([]byte(err.Error()))
	}
}

type FileResponse struct {
	fileName string
	wt       http.ResponseWriter
	isFirst  bool
}

func NewFileResponse(wt http.ResponseWriter, fileName string) *FileResponse {
	strList := strings.Split(fileName, "/")
	fileName = strList[len(strList)-1]
	return &FileResponse{wt: wt, fileName: fileName}
}

func (sel *FileResponse) Write(p []byte) (n int, err error) {
	if !sel.isFirst && len(p) > 10 {
		sel.isFirst = true
		if !ShowWeb(p[:10]) {
			sel.wt.Header().Add("Content-Type", "application/octet-stream")
			sel.wt.Header().Add("Content-Disposition", "attachment; filename=\""+sel.fileName+"\"")
		}
	}
	return sel.wt.Write(p)
}
