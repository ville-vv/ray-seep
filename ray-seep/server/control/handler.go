package control

import (
	"ray-seep/ray-seep/common/errs"
	"ray-seep/ray-seep/dao"
)

type PodHandler struct {
	db *dao.RaySeepServer
}

func NewPodHandler(db *dao.RaySeepServer) *PodHandler {
	return &PodHandler{db: db}
}

func (sel *PodHandler) OnLogin(connId, userId int64, appKey string, token string) (secret string, err error) {
	secret, err = sel.db.UserLogin(userId, appKey, token)
	if err != nil {
		return
	}
	if secret == "" {
		return "", errs.ErrSecretIsInValid
	}
	// TODO redis
	return
}

func (sel *PodHandler) OnLogout(name string, id int64) error {
	return nil
}
