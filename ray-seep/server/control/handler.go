package control

import (
	"ray-seep/ray-seep/dao"
)

type PodHandler struct {
	db *dao.RaySeepServer
}

func NewPodHandler(db *dao.RaySeepServer) *PodHandler {
	return &PodHandler{db: db}
}

func (sel *PodHandler) OnLogin(connId, userId int64, appId string, token string) (secret string, err error) {
	secret, err = sel.db.UserLogin(userId, appId, token)
	if err != nil {
		return
	}
	// TODO redis
	return
}

func (sel *PodHandler) OnLogout(name string, id int64) error {
	return nil
}
