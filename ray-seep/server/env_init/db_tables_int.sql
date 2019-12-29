-- 用户账号数据库
CREATE TABLE IF NOT EXISTS user_accounts (
seq int(10) NOT NULL AUTO_INCREMENT PRIMARY KEY,
user_id int(10) NOT NULL UNIQUE COMMENT '用户ID',
user_name varchar(128) NOT NULL COMMENT '用户名',
secret varchar(128) NOT NULL UNIQUE,
app_key varchar(128) NOT NULL UNIQUE,
INDEX idx_user_account_seq (seq),
INDEX idx_user_account_app_id (app_key)
)ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;

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