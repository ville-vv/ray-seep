package dao

import "ray-seep/ray-seep/conf"

type RaySeepServer struct {
	rdsDb *RedisClient
	sqlDb *MysqlClient
}

func NewRaySeepServer(cfg *conf.DataBaseSrv) *RaySeepServer {
	return &RaySeepServer{
		//rdsDb: NewRedisClient(cfg.Redis),
		sqlDb: NewMysqlClient(cfg.Mysql),
	}
}

func (sel *RaySeepServer) Close() {
	_ = sel.sqlDb.Close()
	//_ = sel.rdsDb.Close()
}

func (sel *RaySeepServer) UserLogin(userId int64, appKey string, token string) (string, error) {
	secret, err := sel.sqlDb.UserAuth(userId, appKey)
	if err != nil {
		return "", err
	}

	return secret, err
}
