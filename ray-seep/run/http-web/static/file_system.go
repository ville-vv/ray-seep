package static

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
		<h1> Ray Seep(射线渗透) Test</h1>
		<span> 你可以再文件中添加几个文件然后刷新网页看看有什么变化。</span>
		<p>1、这个是一个部署在本地的内网web服务网站，此服务网站试用以体验 Ray Seep 工具的内网穿透能力。</p>
		<p>2、转发次网站的数据技术是由 Ray-Seep 工具实现。</p>
		<p>3、部署此工具可以学生共享外网服务器，在宿舍自己电脑学习开发web应用，并且可以在公网访问。</p>
		<p>4、部署此工具可以提供企业，在本地开发软件，并可以与第三方公司联调本地HTTP API接口。</p>
		<h4>这个网站的示例功能是显示当前运行环境目录的所有文件名称，以下是你当前软件目录下的文件：<h4>
	</div>
	<div style="background-color: #c6d5e9">
		{{.Context}}
	<div>
</div>

</body>
</html>
`

func UseTemplate(w io.Writer, tmpl string, data interface{}) error {
	tp, err := template.New("tmpl").Parse(htmlPageTemp)
	if err != nil {
		return err
	}
	return tp.Execute(w, data)
}

type FileSystem struct {
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
	for _, file := range fileList {
		one := fmt.Sprintf("<div><a href=\"%s/%s\">%s<a/></div>", path, file.Name(), file.Name())
		_, _ = w.Write([]byte(one))
	}
	return nil
}
func (f *FileSystem) ServeHTTP(rsp http.ResponseWriter, req *http.Request) {
	var err error
	uriPath := filepath.Join("./", req.URL.Path)
	fmt.Println("文件名称：", uriPath)
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
	if err = UseTemplate(rsp, htmlPageTemp, map[string]string{"Context": buf.String()}); err != nil {
		_, _ = rsp.Write([]byte(err.Error()))
	}
}
