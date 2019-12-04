package env_init

import (
	"database/sql"
	"fmt"
	"strings"
	"vilgo/vsql"
)

var (
	TableUserAccount = `
CREATE TABLE IF NOT EXISTS user_accounts (
seq int(10) NOT NULL AUTO_INCREMENT PRIMARY KEY,
user_id int(10) NOT NULL UNIQUE COMMENT '用户ID',
secret varchar(128) NOT NULL UNIQUE,
app_key varchar(128) NOT NULL UNIQUE,
INDEX idx_user_account_seq (seq),
INDEX idx_user_account_app_id (app_key)
)ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;
`
	TableUserProtocol = `
CREATE TABLE IF NOT EXISTS user_protocols(
seq int(10) NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '',
created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '',
updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '',
user_id int(10) NOT NULL COMMENT '',
protocol_name varchar (20) NOT NULL COMMENT '',
protocol_port int(10) NOT NULL UNIQUE COMMENT '' ,
INDEX idx_user_protocols_user_id (user_id),
INDEX idx_user_protocols_protocol_name (protocol_name)
)ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;
`
)

func SnakeToCameString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

func CamelToSnakeString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

type MysqlMigrate struct {
	cnf    *vsql.MySqlCnf
	driver *vsql.NormalSqlDrive
	db     *sql.DB
}

func NewMysqlMigrate(cnf *vsql.MySqlCnf) *MysqlMigrate {
	m := &MysqlMigrate{cnf: cnf}
	db := vsql.NewNormalSqlDrive(cnf)
	if err := db.Open(); err != nil {
		if db.GetDefDb() == nil {
			panic(err)
		}
	}
	m.driver = db
	return m
}

func (sel *MysqlMigrate) Close() {
	if sel.db != nil {
		_ = sel.db.Close()
	}
}
func (sel *MysqlMigrate) CleanDatabase() {
}

func (sel *MysqlMigrate) CreateDatabase() error {
	for _, v := range sel.cnf.Databases {
		sqlStr := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET utf8;", v)
		if _, err := sel.driver.GetDefDb().Exec(sqlStr); err != nil {
			fmt.Println("创建数据库错误：", sqlStr)
			return err
		}
		err := sel.driver.Add(v)
		if err != nil {
			return err
		}
		sel.db = sel.driver.GetDb(v)
	}
	return nil
}

func (sel *MysqlMigrate) CreateTable(nms ...interface{}) error {
	for _, v := range nms {
		if _, err := sel.db.Exec(v.(string)); err != nil {
			return err
		}
	}
	return nil
}

func (sel *MysqlMigrate) batchExec(nms []string) error {
	for _, v := range nms {
		_, err := sel.db.Exec(v)
		if err != nil {
			return err
		}
	}
	return nil
}
