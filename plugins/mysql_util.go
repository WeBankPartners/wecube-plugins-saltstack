package plugins

import (
	"database/sql"
	"fmt"
	"os/exec"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

var DB *sql.DB

func initDB(host, port, loginUser, loginPwd, dbName string) (*sql.DB, error) {
	var err error
	path := strings.Join([]string{loginUser, ":", loginPwd, "@tcp(", host, ":", port, ")/", dbName, "?charset=utf8"}, "")
	logrus.Infof("Init mysql db path=[%v]", path)

	DB, err = sql.Open("mysql", path)
	if err != nil {
		logrus.Errorf("opening mysql db[%v] meet err=%v", dbName, err)
		return nil, err
	}
	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)

	if err := DB.Ping(); err != nil {
		logrus.Errorf("opening mysql db[%v] failed, err=%v", dbName, err)
		return nil, err
	}

	logrus.Infof("connected mysql db[%v] successfully", dbName)
	return DB, nil
}

func getAllUserByDB(host, port, loginUser, loginPwd, dbName string) ([]string, error) {
	users := []string{}
	var err error

	// initDB param dbName = "mysql", not getUserByDB.dbName
	DB, err = initDB(host, port, loginUser, loginPwd, "mysql")
	if err != nil {
		logrus.Errorf("getting user by db[%v] failed, err=%v ", dbName, err)
		return users, err
	}
	defer DB.Close()

	querySql := fmt.Sprintf("select User from db where db.Db='%s'", dbName)
	rows, err := DB.Query(querySql)
	if err != nil {
		logrus.Errorf("db.query meet err=%v", err)
		return users, err
	}

	for rows.Next() {
		var user string
		err := rows.Scan(&user)
		if err != nil {
			logrus.Errorf("rows.Scan meet err=%v", err)
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

func getAllDBByUser(host, port, loginUser, loginPwd, userName string) ([]string, error) {
	dbs := []string{}
	var err error
	// initDB param dbName = "mysql".
	DB, err = initDB(host, port, loginUser, loginPwd, "mysql")
	if err != nil {
		logrus.Errorf("init myhsql db failed, err=%v ", err)
		return dbs, err
	}
	defer DB.Close()

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

func checkDBExistOrNot(host, port, loginUser, loginPwd, dbName string) (bool, error) {
	var err error
	// initDB param dbName = "mysql", not getUserByDB.dbName
	DB, err = initDB(host, port, loginUser, loginPwd, "mysql")
	if err != nil {
		logrus.Errorf("init myhsql db failed, err=%v ", err)
		return false, err
	}
	defer DB.Close()

	querySql := fmt.Sprintf("SELECT 1 FROM mysql.db WHERE Db = '%s'", dbName)
	rows, err := DB.Query(querySql)
	if err != nil {
		logrus.Errorf("db.query meet err=%v", err)
		return false, err
	}

	return rows.Next(), nil
}

func checkUserExistOrNot(host, port, loginUser, loginPwd, userName string) (bool, error) {
	var err error
	// initDB param dbName = "mysql".
	DB, err = initDB(host, port, loginUser, loginPwd, "mysql")
	if err != nil {
		logrus.Errorf("init myhsql db failed, err=%v ", err)
		return false, err
	}
	defer DB.Close()

	querySql := fmt.Sprintf("SELECT 1 FROM mysql.user WHERE user = '%s'", userName)
	rows, err := DB.Query(querySql)
	if err != nil {
		logrus.Errorf("db.query meet err=%v", err)
		return false, err
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
	fmt.Printf("runDatabaseCommand(%v) output=%v,err=%v\n", command, string(out), err)
	return err
}
