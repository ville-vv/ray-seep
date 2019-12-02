-- 用户账号数据库
CREATE TABLE IF NOT EXISTS user_account (
seq int(10) NOT NULL AUTO_INCREMENT,
user_id int(10) NOT NULL PRIMARY KEY,
secret varchar(128) NOT NULL UNIQUE,
app_key varchar(128) NOT NULL UNIQUE,
INDEX idx_user_account_seq (seq),
INDEX idx_user_account_app_id (app_key)
)ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;