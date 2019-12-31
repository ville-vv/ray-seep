package input

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"ray-seep/ray-seep/common/util"
	"ray-seep/ray-seep/run/admin/store"
	"time"
)

var (
	String string
	Number int
	Input  string
)

func getInput() {
	f := bufio.NewReader(os.Stdin) //读取输入的内容
	for {
		fmt.Print("请输入用户名-> ")
		Input, _ = f.ReadString('\n') //定义一行输入的内容分隔符。
		fmt.Printf("您输入的是:%s\n", Input)
		if String == "stop" {
			break
		}
		fmt.Printf("输入密码-> ")
		pas, _ := GetPass(os.Stdin, os.Stdout)
		fmt.Printf("您输入的密码是:%s \n", string(pas))
	}
}

var getCh = func(r io.Reader) (byte, error) {
	buf := make([]byte, 1)
	if n, err := r.Read(buf); n == 0 || err != nil {
		if err != nil {
			return 0, err
		}
		return 0, io.EOF
	}
	return buf[0], nil
}

// 获取密码
func GetPass(r *os.File, w io.Writer) ([]byte, error) {
	var err error
	var pass, bs []byte
	//主要就是这一行终端文件替换
	if terminal.IsTerminal(int(r.Fd())) {
		oldState, err := terminal.MakeRaw(int(r.Fd()))
		if err != nil {
			return pass, err
		}
		defer func() {
			terminal.Restore(int(r.Fd()), oldState)
			fmt.Println("")
		}()
	}

	var counter int
	for counter = 0; counter <= 200; counter++ {
		if v, e := getCh(r); e != nil {
			err = e
			break
		} else if v == 127 || v == 8 {
			if l := len(pass); l > 0 {
				pass = pass[:l-1]
				fmt.Fprint(w, string(bs))
			}
		} else if v == 13 || v == 10 {
			break
		} else if v == 3 {
			err = errors.New("interrupted")
			break
		} else if v != 0 {
			pass = append(pass, v)
			fmt.Fprint(w, "*")
		}
	}
	return pass, err

}

func input(r *bufio.Reader, notice string, isPas bool) string {
	fmt.Print(notice)
	if isPas {
		pas, _ := GetPass(os.Stdin, os.Stdout)
		return string(pas)
	}
	dt, _ := r.ReadString('\n')
	if String == "exit" {
		os.Exit(0)
	}
	return dt[:len(dt)-1]
}

func InputStoreInfo() (addr, user, passwd string) {
	f := bufio.NewReader(os.Stdin) //读取输入的内容
	addr = input(f, "请输入数据库地址 ->: ", false)
	user = input(f, "请输入数据库用户名 ->: ", false)
	passwd = input(f, "请输入数据库密码 ->: ", true)
	return
}

func UserAccountInfo() (*store.RayAccount, *store.RayProtocol) {
	ac := new(store.RayAccount)
	pt := new(store.RayProtocol)
	f := bufio.NewReader(os.Stdin) //读取输入的内容
	ac.UserName = input(f, "请输入用户名： ->: ", false)
	fmt.Println("您输入的用户名为：", ac.UserName)
	yes := input(f, "请输入确认 yes/no ? ->: ", false)
	if yes != "yes" {
		os.Exit(1)
	}
	ac.Secret = input(f, "请输入生成 Secret 随机字符 ->: ", false)
	ac.AppKey = util.RandString(32)
	content := fmt.Sprintf("%s&%s&%d", ac.Secret, ac.AppKey, time.Now().UnixNano())
	ac.Secret = util.GetMd5String(util.HmacSha256String(ac.UserName, content))
	ac.UserId = int64(util.GenRandID())
	pt.ProtocolPort = input(f, "请输入Http端口号 ->: ", false)
	pt.ProtocolName = "http"
	pt.UserId = ac.UserId

	fmt.Println("UserId: ", ac.UserId)
	fmt.Println("Secret: ", ac.Secret)
	fmt.Println("AppKey: ", ac.AppKey)

	return ac, pt
}
