package dao

import (
	"ray-seep/ray-seep/conf"
)

type RaySeepServer struct {
	rdsDb *RedisClient
	sqlDb *MysqlClient
	user  *UserConfig
}

func NewRaySeepServer(cfg *conf.Server) *RaySeepServer {
	r := &RaySeepServer{
		user: NewUserConfig(cfg.User),
	}
	if cfg.DataBase.OpenMysql {
		r.sqlDb = NewMysqlClient(cfg.DataBase.Mysql)
	}
	if cfg.DataBase.OpenRedis {
		r.rdsDb = NewRedisClient(cfg.DataBase.Redis)
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

func (sel *RaySeepServer) UserLogin(userId int64, appKey string, token string) (string, error) {
	secret, err := sel.user.UserAuth(userId, appKey)
	if err != nil {
		return "", err
	}
	return secret, err
}
