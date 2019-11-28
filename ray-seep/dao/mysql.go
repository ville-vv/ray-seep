package dao

import (
	"database/sql"
	"fmt"
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

func (sel *MysqlClient) UserAuth(userId int64, appId string) (string, error) {
	sqlStr := fmt.Sprintf("SELECT secret FROM user_account WHERE user_id='%d' and app_id='%s' LIMIT 1; ", userId, appId)
	rows, err := sel.db.Query(sqlStr)
	if err != nil {
		return "", err
	}
	var secret string
	for rows.Next() {
		if err := rows.Scan(&secret); err != nil {
			return "", nil
		}
	}
	return secret, err
}
