package databus

import (
	"database/sql"
	"fmt"
	"github.com/vilsongwei/vilgo/vsql"
	"ray-seep/ray-seep/model"
)

type MysqlClient struct {
	db *sql.DB
}

func NewMysqlClient(cnf *vsql.MySqlCnf) *MysqlClient {
	db := vsql.NewNormalSqlDrive(cnf)
	if err := db.Open(); err != nil {
		panic(err)
	}
	return &MysqlClient{
		db: db.GetDb(cnf.Databases[0]),
	}
}
func (sel *MysqlClient) Close() error {
	return sel.db.Close()
}

func (sel *MysqlClient) UserAuth(userId int64, userName, appKey string, ul *model.UserLoginDao) error {
	sqlStr := fmt.Sprintf("SELECT ua.secret, up.protocol_port FROM user_accounts as ua "+
		"JOIN user_protocols as up ON up.user_id=ua.user_id"+
		" WHERE ua.user_id='%d' and ua.user_name='%s' and ua.app_key='%s' and up.protocol_name='http' LIMIT 1; ", userId, userName, appKey)
	rows, err := sel.db.Query(sqlStr)
	if err != nil {
		return err
	}
	for rows.Next() {
		if err := rows.Scan(&ul.Secret, &ul.HttpPort); err != nil {
			return err
		}
	}
	return nil
}
