package databus

import (
	"github.com/vilsongwei/vilgo/vlog"
	"ray-seep/ray-seep/common/errs"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/model"
)

type BaseDao interface {
	UserLogin(connId int64, userId int64, user string, appKey string, token string) (*model.UserLoginDao, error)
	GetToken(connId int64, user string) string
	DelToken(connId int64, user string, isDelKeys bool)
	Close()
}

func NewDao(cfg *conf.Server) BaseDao {
	if cfg.DataBase.OpenDb {
		return NewRaySeepServerDao(cfg)
	}
	vlog.INFO("dao use NewNotSqlDao")
	return NewNotSqlDao(cfg)
}

type RaySeepServerDao struct {
	rdsDb *RedisClient
	sqlDb *MysqlClient
}

func NewRaySeepServerDao(cfg *conf.Server) *RaySeepServerDao {
	r := &RaySeepServerDao{
		sqlDb: NewMysqlClient(cfg.DataBase.Mysql),
		rdsDb: NewRedisClient(cfg.DataBase.Redis),
	}
	return r
}

func (sel *RaySeepServerDao) Close() {
	if sel.sqlDb != nil {
		_ = sel.sqlDb.Close()
	}
	if sel.rdsDb != nil {
		_ = sel.rdsDb.Close()
	}
}

func (sel *RaySeepServerDao) UserLogin(connId int64, userId int64, user string, appKey string, token string) (*model.UserLoginDao, error) {
	ul := &model.UserLoginDao{}
	if err := sel.sqlDb.UserAuth(userId, user, appKey, ul); err != nil {
		return nil, err
	}
	if ul.Secret == "" {
		return nil, errs.ErrUserInfoValidFail
	}
	return ul, sel.rdsDb.SetUserToken(connId, user, token)
}

func (sel *RaySeepServerDao) GetToken(connId int64, user string) string {
	return sel.rdsDb.GetUserToken(connId, user)
}

func (sel *RaySeepServerDao) DelToken(connId int64, user string, isDelKeys bool) {
	_ = sel.rdsDb.DelUserToken(connId, user, isDelKeys)
	return
}
