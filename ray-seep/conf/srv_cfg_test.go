// @File     : srv_cfg_test
// @Author   : Ville
// @Time     : 19-10-15 上午10:27
// conf
package conf

import (
	"github.com/vilsongwei/vilgo/vfile"
	"testing"
)

func TestInitServer(t *testing.T) {
	InitServer()
}

func TestInitServer2(t *testing.T) {
	fName, err := vfile.WriteToExecuteDir("server_config_test.toml", serverDefaultConfig)
	if err != nil {
		t.Error(err)
		return
	}
	defer vfile.Remove(fName)
	InitServer(fName)

}
