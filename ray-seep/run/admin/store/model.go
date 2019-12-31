package store

import "fmt"

type RayAccount struct {
	Seq      int64
	UserId   int64
	UserName string
	Secret   string
	AppKey   string
}

func (r *RayAccount) InsertStr() string {
	sql := `insert into user_accounts(user_id, user_name, secret, app_key)value(%d, '%s','%s', '%s');`
	return fmt.Sprintf(sql, r.UserId, r.UserName, r.Secret, r.AppKey)
}

type RayProtocol struct {
	UserId       int64
	ProtocolName string
	ProtocolPort string
}

func (r *RayProtocol) InsertStr() string {
	sql := `insert into user_protocols(user_id,protocol_name,protocol_port)value(%d, '%s','%s');`
	return fmt.Sprintf(sql, r.UserId, r.ProtocolName, r.ProtocolPort)
}
