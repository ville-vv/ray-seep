package databus

import (
	"fmt"
	"ray-seep/ray-seep/conf"
	"ray-seep/ray-seep/model"
	"sync"
)

type NotSqlDao struct {
	lock   sync.RWMutex
	user   *UserConfig
	tokens map[string]string
}

func NewNotSqlDao(cfg *conf.Server) *NotSqlDao {
	return &NotSqlDao{
		user:   NewUserConfig(cfg.User),
		tokens: make(map[string]string),
	}
}

func (sel *NotSqlDao) Close() {
}

func (sel *NotSqlDao) UserLogin(connId int64, userId int64, user string, appKey string, token string) (*model.UserLoginDao, error) {
	ul := &model.UserLoginDao{}
	if err := sel.user.UserAuth(userId, appKey, ul); err != nil {
		return nil, err
	}
	sel.lock.Lock()
	sel.tokens[fmt.Sprintf("login_token_%s_%d", user, connId)] = token
	sel.lock.Unlock()
	return ul, nil
}

func (sel *NotSqlDao) GetToken(connId int64, user string) string {
	sel.lock.RLock()
	defer sel.lock.RUnlock()
	return sel.tokens[fmt.Sprintf("login_token_%s_%d", user, connId)]
}

func (sel *NotSqlDao) DelToken(connId int64, user string, isDelKeys bool) {
	sel.lock.Lock()
	defer sel.lock.Unlock()
	delete(sel.tokens, fmt.Sprintf("login_token_%s_%d", user, connId))
}
