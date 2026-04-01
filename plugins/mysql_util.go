package plugins

import (
	"database/sql"
	"fmt"
	"os/exec"
	"strings"

	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
	_ "github.com/go-sql-driver/mysql"
)

// 安全处理 MySQL 密码中的特殊字符
func escapeMysqlPassword(password string) string {
	// 按照 MySQL 的转义规则处理特殊字符
	password = strings.ReplaceAll(password, "\\", "\\\\")  // 反斜杠必须最先处理
	password = strings.ReplaceAll(password, "'", "\\'")    // 单引号
	password = strings.ReplaceAll(password, "\"", "\\\"")  // 双引号
	password = strings.ReplaceAll(password, "\n", "\\n")   // 换行符
	password = strings.ReplaceAll(password, "\r", "\\r")   // 回车符
	password = strings.ReplaceAll(password, "\t", "\\t")   // 制表符
	password = strings.ReplaceAll(password, "\b", "\\b")   // 退格符
	password = strings.ReplaceAll(password, "\f", "\\f")   // 换页符
	password = strings.ReplaceAll(password, "\v", "\\v")   // 垂直制表符
	password = strings.ReplaceAll(password, "\000", "\\0") // 空字符
	return password
}

var DB *sql.DB

func initDB(host, port, loginUser, loginPwd, dbName string) error {
	var err error
	connParam := "?charset=utf8"
	if MysqlSSLEnable {
		connParam = "?tls=true&charset=utf8"
	}
	path := strings.Join([]string{loginUser, ":", loginPwd, "@tcp(", host, ":", port, ")/", dbName, connParam}, "")

	DB, err = sql.Open("mysql", path)
	if err != nil {
		return fmt.Errorf("connect to mysql fail,%s", err.Error())
	}
	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("keep alive to mysql fail,%s", err.Error())
	}

	return nil
}

func getAllUserByDB(host, port, loginUser, loginPwd, dbName, language string) ([]string, error) {
	users := []string{}

	// initDB param dbName = "mysql", not getUserByDB.dbName
	err := initDB(host, port, loginUser, loginPwd, "mysql")
	if err != nil {
		return users, getMysqlConnectError(language, err)
	}

	querySql := fmt.Sprintf("select User from db where db.Db='%s'", dbName)
	rows, err := DB.Query(querySql)
	if err != nil {
		return users, fmt.Errorf("Query mysql user fail,%s ", err.Error())
	}

	for rows.Next() {
		var user string
		err := rows.Scan(&user)
		if err != nil {
			return users, fmt.Errorf("Mysql rows scan fail,%s ", err.Error())
		}
		users = append(users, user)
	}
	return users, nil
}

func getAllDBByUser(host, port, loginUser, loginPwd, userName, language string) ([]string, error) {
	dbs := []string{}
	// initDB param dbName = "mysql".
	err := initDB(host, port, loginUser, loginPwd, "mysql")
	if err != nil {
		err = getMysqlConnectError(language, err)
		return dbs, err
	}

	querySql := fmt.Sprintf("select Db from db where db.User='%s'", userName)
	rows, err := DB.Query(querySql)
	if err != nil {
		return dbs, fmt.Errorf("Query mysql database user fail,%s ", err.Error())
	}
	for rows.Next() {
		var db string
		err := rows.Scan(&db)
		if err != nil {
			return dbs, fmt.Errorf("Mysql rows scan fail,%s ", err.Error())
		}
		dbs = append(dbs, db)
	}
	return dbs, nil
}

func checkDBExistOrNot(host, port, loginUser, loginPwd, dbName, language string) (bool, error) {
	// initDB param dbName = "mysql", not getUserByDB.dbName
	err := initDB(host, port, loginUser, loginPwd, "mysql")
	if err != nil {
		return false, getMysqlConnectError(language, err)
	}

	rows, err := DB.Query("SHOW DATABASES LIKE ?", dbName)
	if err != nil {
		return false, fmt.Errorf("Query mysql database fail,%s ", err.Error())
	}
	defer rows.Close()

	return rows.Next(), nil
}

func checkUserExistOrNot(host, port, loginUser, loginPwd, userName, language string) (bool, error) {
	// initDB param dbName = "mysql".
	err := initDB(host, port, loginUser, loginPwd, "mysql")
	if err != nil {
		return false, getMysqlConnectError(language, err)
	}

	querySql := fmt.Sprintf("SELECT 1 FROM mysql.user WHERE user = '%s'", userName)
	rows, err := DB.Query(querySql)
	if err != nil {
		return false, fmt.Errorf("Query mysql user fail,%s ", err.Error())
	}

	return rows.Next(), nil
}

func runDatabaseCommand(host string, port string, loginUser string, loginPwd string, cmd string) error {
	argv := []string{
		"-h" + host,
		"-u" + loginUser,
		"-p" + loginPwd,
		"-P" + port,
		"-e",
		cmd,
	}
	if !MysqlSSLEnable {
		argv = append([]string{"--ssl-mode=DISABLED"}, argv...)
	}
	command := exec.Command("/usr/bin/mysql", argv...)
	out, err := command.CombinedOutput()
	if err != nil {
		log.Logger.Debug("Run mysql command", log.StringList("command", argv), log.String("output", string(out)), log.Error(err))
		return fmt.Errorf("output:%s,error:%s", string(out), err.Error())
	}
	return nil
}
