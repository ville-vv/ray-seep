package main

import (
	"flag"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/dao"
	"vilgo/vlog"
)

var (
	configPath = ""
	help       bool
)

func main() {
	vlog.DefaultLogger()
	flag.StringVar(&configPath, "c", "", "the config file")
	flag.BoolVar(&help, "h", false, "the tool use help")
	flag.Parse()
	if help {
		flag.PrintDefaults()
		return
	}
	cfg := conf.InitServer(configPath)
	mig := dao.NewMysqlMigrate(cfg.DataBase.Mysql)
	if err := mig.CreateDatabase(); err != nil {
		vlog.ERROR("create databases error")
	}

	if err := mig.CreateTable(dao.TableUserAccount); err != nil {
		vlog.ERROR("create table error %s", err.Error())
		return
	}
	vlog.INFO("数据库初始化成功")
}
