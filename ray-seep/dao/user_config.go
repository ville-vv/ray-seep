package dao

import (
	"ray-seep/ray-seep/common/errs"
	"ray-seep/ray-seep/conf"
	"sync"
)

type UserConfig struct {
	lock sync.RWMutex
	user map[int64]*conf.User
}

func NewUserConfig(user map[string]*conf.User) *UserConfig {
	u := make(map[int64]*conf.User)

	for _, v := range user {
		u[v.UserId] = &conf.User{
			UserId:   v.UserId,
			UserName: v.UserName,
			Secret:   v.Secret,
			AppKey:   v.AppKey,
		}
	}

	return &UserConfig{user: u}
}

func (sel *UserConfig) UserAuth(userId int64, appKey string) (string, error) {
	sel.lock.RLock()
	defer sel.lock.RUnlock()
	user, ok := sel.user[userId]
	if !ok {
		return "", errs.ErrSecretIsInValid
	}
	return user.Secret, nil
}
