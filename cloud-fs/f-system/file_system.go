package f_system

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
	<title>Cloud-FS</title>
</head>



<body>
	<div style="margin: 0px 200px 100px 200px">
		<div style="">
			<div style="text-align: center">
				<h1>Cloud-FS</h1>
			</div>
			<div style="font-weight:unset;color:#000000">
				<label id="time_now"></label>
			</div>
		</div>
		<div style="background-color: #dfe5ec; margin-top: 10px">
			{{.Context}}
		</div>
	</div>
</body>
<script type="text/javascript">
    document.getElementById("time_now").innerHTML = getNowDate();
    function getNowDate() {
        var date = new Date();
        var year = date.getFullYear() ;// 年
        var month = date.getMonth() + 1; // 月
        var day  = date.getDate(); // 日
        var hour = date.getHours(); // 时
        var minutes = date.getMinutes(); // 分
        var seconds = date.getSeconds() //秒
        var weekArr = ['星期日','星期一', '星期二', '星期三', '星期四', '星期五', '星期六' ];
        var week = weekArr[date.getDay()];
		// 给一位数数据前面加 “0”
        if (month >= 1 && month <= 9) {
            month = "0" + month;
        }
        if (day >= 0 && day <= 9) {
            day = "0" + day;
        }
        if (hour >= 0 && hour <= 9) {
            hour = "0" + hour;
        }
        if (minutes >= 0 && minutes <= 9) {
            minutes = "0" + minutes;
        }
        if (seconds >= 0 && seconds <= 9) {
            seconds = "0" + seconds;
        }
         return year + "年" + month + "月" + day + "日" + hour +":" + minutes+ ":"+ seconds +"&nbsp&nbsp"+ week;
    }
    setInterval(function () {
        document.getElementById("time_now").innerHTML = getNowDate()
    }, 1000);

</script>
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
	root string // 文件系统的根目录
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

// 文件展示
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
