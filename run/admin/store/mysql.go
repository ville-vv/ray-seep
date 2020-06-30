package store

import (
	"database/sql"
	"fmt"
	"github.com/vilsongwei/vilgo/vsql"
)

type MySqlStore struct {
	db *sql.DB
}

func NewMysqlStore(addr, user, passwd string) *MySqlStore {
	fmt.Println("数据库地址：", addr)
	fmt.Println("数据库用户：", user)
	db := vsql.NewNormalSqlDrive(&vsql.MySqlCnf{
		Version:   "8",
		UserName:  user,
		Address:   addr,
		Password:  passwd,
		Default:   "information_schema",
		MaxIdles:  10,
		MaxOpens:  100,
		Databases: []string{"ray_seep"},
	})
	if err := db.Open(); err != nil {
		panic(err)
	}
	return &MySqlStore{
		db: db.GetDb("ray_seep"),
	}
}

func (m *MySqlStore) AddRaySeepUser(use *RayAccount, pto *RayProtocol) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()
	if _, err := tx.Exec(use.InsertStr()); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(pto.InsertStr()); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
