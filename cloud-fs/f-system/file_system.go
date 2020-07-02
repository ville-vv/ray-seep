package f_system

import (
	"io"
	"net/http"
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

type FileDisplayer interface {
	Display(root string, filePath string, w http.ResponseWriter, req *http.Request) error
}

type FileSystem struct {
	root string        // 文件系统的根目录
	dir  FileDisplayer // 目录文件
	text FileDisplayer // 文本文件
}

func NewFileSystem(root string) *FileSystem {
	if strings.Trim(root, " ") == "" {
		root = "./"
	}
	return &FileSystem{
		root: filepath.Join("", root),
		dir:  &DirFile{},
		text: &TextFile{},
	}
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
	switch isDir {
	case false:
		err = f.text.Display(f.root, uriPath, rsp, req)
	default:
		err = f.dir.Display(f.root, uriPath, rsp, req)
	}
	if err != nil {
		_, _ = rsp.Write([]byte(err.Error()))
	}
}
