// @File     : cli_cfg_test
// @Author   : Ville
// @Time     : 19-10-15 上午10:27
// conf
package conf

import (
	"fmt"
	"github.com/vilsongwei/vilgo/vfile"

	"testing"
)

func TestInitClient(t *testing.T) {
	cliCfg := InitClient()
	fmt.Println(cliCfg.Pxy)
	//fmt.Println(cliCfg.Node)
}

func TestInitClient2(t *testing.T) {
	fName, err := vfile.WriteToExecuteDir("client_config_test.toml", clientDefaultConfig)
	if err != nil {
		t.Error(err)
		return
	}
	defer vfile.Remove(fName)
	cliCfg := InitClient(fName)
	fmt.Println(cliCfg.Pxy)
	//fmt.Println(cliCfg.Node)

}
