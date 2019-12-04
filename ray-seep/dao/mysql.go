package dao

import (
	"database/sql"
	"fmt"
	"ray-seep/ray-seep/model"
	"vilgo/vsql"
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

func (sel *MysqlClient) UserAuth(userId int64, appKey string, ul *model.UserLoginDao) error {
	sqlStr := fmt.Sprintf("SELECT secret FROM user_account WHERE user_id='%d' and app_key='%s' LIMIT 1; ", userId, appKey)
	rows, err := sel.db.Query(sqlStr)
	if err != nil {
		return err
	}
	for rows.Next() {
		if err := rows.Scan(&ul.Secret); err != nil {
			return nil
		}
	}
	return err
}
