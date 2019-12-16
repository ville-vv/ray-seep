package control

import (
	"ray-seep/ray-seep/common/errs"
	"ray-seep/ray-seep/databus"
	"ray-seep/ray-seep/model"
)

type PodHandler struct {
	db databus.BaseDao
}

func NewPodHandler(db databus.BaseDao) *PodHandler {
	return &PodHandler{db: db}
}

func (sel *PodHandler) OnLogin(connId, userId int64, appKey string, token string) (loginDao *model.UserLoginDao, err error) {
	loginDao, err = sel.db.UserLogin(userId, appKey, token)
	if err != nil {
		return
	}
	if loginDao.Secret == "" {
		return nil, errs.ErrSecretIsInValid
	}
	if loginDao.HttpPort == "" {
		return nil, errs.ErrHttpPortIsInValid
	}
	// TODO redis
	return
}

func (sel *PodHandler) OnLogout(name string, id int64) error {
	return nil
}
