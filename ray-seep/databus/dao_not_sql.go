package databus

import (
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/model"
)

type NotSqlDao struct {
	user *UserConfig
}

func NewNotSqlDao(cfg *conf.Server) *NotSqlDao {
	return &NotSqlDao{
		user: NewUserConfig(cfg.User),
	}
}

func (sel *NotSqlDao) Close() {
}

func (sel *NotSqlDao) UserLogin(connId int64, userId int64, user string, appKey string, token string) (*model.UserLoginDao, error) {
	ul := &model.UserLoginDao{}
	if err := sel.user.UserAuth(userId, appKey, ul); err != nil {
		return nil, err
	}
	return ul, nil
}

func (sel *NotSqlDao) GetToken(connId int64, user string) string {
	return ""
}
