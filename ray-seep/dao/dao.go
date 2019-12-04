package dao

import (
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/model"
	"vilgo/vlog"
)

type BaseDao interface {
	UserLogin(userId int64, appKey string, token string) (*model.UserLoginDao, error)
	Close()
}

func NewDao(cfg *conf.Server) BaseDao {
	if cfg.DataBase.OpenMysql {
		return NewRaySeepServer(cfg)
	}
	vlog.INFO("dao use NewNotSqlDao")
	return NewNotSqlDao(cfg)
}

type RaySeepServer struct {
	rdsDb *RedisClient
	sqlDb *MysqlClient
}

func NewRaySeepServer(cfg *conf.Server) *RaySeepServer {
	r := &RaySeepServer{
		sqlDb: NewMysqlClient(cfg.DataBase.Mysql),
		rdsDb: NewRedisClient(cfg.DataBase.Redis),
	}
	return r
}

func (sel *RaySeepServer) Close() {
	if sel.sqlDb != nil {
		_ = sel.sqlDb.Close()
	}
	if sel.rdsDb != nil {
		_ = sel.rdsDb.Close()
	}
}

func (sel *RaySeepServer) UserLogin(userId int64, appKey string, token string) (*model.UserLoginDao, error) {
	ul := &model.UserLoginDao{}
	if err := sel.sqlDb.UserAuth(userId, appKey, ul); err != nil {
		return nil, err
	}
	return ul, nil
}
