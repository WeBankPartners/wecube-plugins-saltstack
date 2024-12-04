package plugins

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
)

// 执行 redis-cli 命令并返回输出
func runRedisCli(args ...string) (string, error) {
	cmd := exec.Command("redis-cli", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Logger.Debug("Run redis command", log.StringList("command", args), log.String("output", string(output)), log.Error(err))
	}

	return string(output), err
}

// 检查redis用户是否存在
func redisCheckUserExistedOrNot(host, port, adminUser, adminPassword, userName string) (isExisted bool, err error) {
	isExisted = false

	args := []string{
		"-h", host,
		"-p", port,
		"-u", adminUser,
		"-a", adminPassword,
		"ACL", "GETUSER", userName}

	output, tmpErr := runRedisCli(args...)
	if tmpErr != nil {
		if strings.Contains(output, "no such user") {
			return
		}
		err = fmt.Errorf("output:%s, error:%s", output, err.Error())
		return
	}
	isExisted = true
	return
}

// 创建redis用户
func redisCreateUser(host, port, adminUser, adminPassword, userName, password, userReadKeyPrefix, userWriteKeyPrefix string) (err error) {
	args := []string{
		"-h", host,
		"-p", port,
		"-u", adminUser,
		"-a", adminPassword,
		"ACL", "SETUSER", userName,
		"on", ">" + password,
	}

	if userReadKeyPrefix != "" {
		args = append(args, "+get", "+mget", "~"+userReadKeyPrefix+"*")
	}
	if userWriteKeyPrefix != "" {
		args = append(args, "+set", "+hset", "+lpush", "+rpush", "~"+userWriteKeyPrefix+"*")
	}

	output, tmpErr := runRedisCli(args...)
	if tmpErr != nil {
		err = fmt.Errorf("output:%s, error:%s", output, tmpErr.Error())
		return
	}
	return
}

// 删除redis用户
func redisDeleteUser(host, port, adminUser, adminPassword, userName string) (err error) {
	args := []string{
		"-h", host,
		"-p", port,
		"-u", adminUser,
		"-a", adminPassword,
		"ACL", "DELUSER", userName,
	}

	output, tmpErr := runRedisCli(args...)
	if tmpErr != nil {
		err = fmt.Errorf("output:%s, error:%s", output, tmpErr.Error())
		return
	}
	return
}

// 授予redis用户读权限
func redisGrantReadPermission(host, port, adminUser, adminPassword, userName, keyPrefix string) (err error) {
	args := []string{
		"-h", host,
		"-p", port,
		"-u", adminUser,
		"-a", adminPassword,
		"ACL", "SETUSER", userName,
		"+get", "+mget",
		"~" + keyPrefix + "*",
	}

	output, tmpErr := runRedisCli(args...)
	if tmpErr != nil {
		err = fmt.Errorf("output:%s, error:%s", output, tmpErr.Error())
		return
	}
	return
}

// 撤销redis用户读权限
func redisRevokeReadPermission(host, port, adminUser, adminPassword, userName, keyPrefix string) (err error) {
	args := []string{
		"-h", host,
		"-p", port,
		"-u", adminUser,
		"-a", adminPassword,
		"ACL", "SETUSER", userName,
		"-get", "-mget",
		"~" + keyPrefix + "*",
	}

	output, tmpErr := runRedisCli(args...)
	if tmpErr != nil {
		err = fmt.Errorf("output:%s, error:%s", output, tmpErr.Error())
		return
	}
	return
}

// 授予redis用户写权限
func redisGrantWritePermission(host, port, adminUser, adminPassword, userName, keyPrefix string) (err error) {
	args := []string{
		"-h", host,
		"-p", port,
		"-u", adminUser,
		"-a", adminPassword,
		"ACL", "SETUSER", userName,
		"+set", "+hset", "+lpush", "+rpush",
		"~" + keyPrefix + "*",
	}

	output, tmpErr := runRedisCli(args...)
	if tmpErr != nil {
		err = fmt.Errorf("output:%s, error:%s", output, tmpErr.Error())
		return
	}
	return
}

// 撤销redis用户写权限
func redisRevokeWritePermission(host, port, adminUser, adminPassword, userName, keyPrefix string) (err error) {
	args := []string{
		"-h", host,
		"-p", port,
		"-u", adminUser,
		"-a", adminPassword,
		"ACL", "SETUSER", userName,
		"-set", "-hset", "-lpush", "-rpush",
		"~" + keyPrefix + "*",
	}

	output, tmpErr := runRedisCli(args...)
	if tmpErr != nil {
		err = fmt.Errorf("output:%s, error:%s", output, tmpErr.Error())
		return
	}
	return
}
