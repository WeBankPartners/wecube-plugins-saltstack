package plugins

import (
	"database/sql"
	"fmt"
	"os/exec"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
)

var DB *sql.DB

func initDB(host, port, loginUser, loginPwd, dbName string) error {
	var err error
	path := strings.Join([]string{loginUser, ":", loginPwd, "@tcp(", host, ":", port, ")/", dbName, "?charset=utf8"}, "")

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

func getAllDBByUser(host, port, loginUser, loginPwd, userName string) ([]string, error) {
	dbs := []string{}
	// initDB param dbName = "mysql".
	err := initDB(host, port, loginUser, loginPwd, "mysql")
	if err != nil {
		logrus.Errorf("init myhsql db failed, err=%v ", err)
		return dbs, err
	}

	querySql := fmt.Sprintf("select Db from db where db.User='%s'", userName)
	rows, err := DB.Query(querySql)
	if err != nil {
		logrus.Infof("db.query meet err=%v", err)
		return dbs, err
	}
	for rows.Next() {
		var db string
		err := rows.Scan(&db)
		if err != nil {
			logrus.Infof("rows.Scan meet err=%v", err)
			return dbs, err
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

	querySql := fmt.Sprintf("SELECT 1 FROM mysql.db WHERE Db = '%s'", dbName)
	rows, err := DB.Query(querySql)
	if err != nil {
		return false, fmt.Errorf("Query mysql database fail,%s ", err.Error())
	}

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
		return false, fmt.Errorf("query mysql user fail,%s", err.Error())
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
	command := exec.Command("/usr/bin/mysql", argv...)
	out, err := command.CombinedOutput()
	if err != nil {
		log.Logger.Error("Run mysql command", log.String("command", command.String()), log.String("output", string(out)), log.Error(err))
		return fmt.Errorf("output:%s,error:%s", string(out), err.Error())
	}
	return nil
}
