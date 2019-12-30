package env_init

import (
	"ray-seep/ray-seep/conf"
	"vilgo/vlog"
)

func InitDb(cfg *conf.Server) {
	mig := NewMysqlMigrate(cfg.DataBase.Mysql)
	if err := mig.CreateDatabase(); err != nil {
		vlog.ERROR("create databases error")
		return
	}

	if err := mig.CreateTable(TableUserAccount, TableUserProtocol); err != nil {
		vlog.ERROR("create table error %s", err.Error())
		return
	}

	if err := mig.TableInitDataInsert(TablesInitDataSqls); err != nil {
		vlog.ERROR("init table data error %s", err.Error())
		return
	}

	mig.Close()
	vlog.INFO("数据库初始化成功")
}
