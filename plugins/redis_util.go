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

/*
Redis 支持多种数据类型，每种数据类型都有一组专门的读写命令。以下是常见数据类型及其对应的读写命令：

1.1 字符串（Strings）
读命令：
GET：获取指定键的值。
MGET：获取多个键的值。

写命令：
SET：设置指定键的值。
MSET：同时设置多个键的值。
INCR / DECR：对键的整数值进行递增或递减。
APPEND：在键的现有值后追加数据。

1.2 哈希（Hashes）
读命令：
HGET：获取哈希表中指定字段的值。
HGETALL：获取哈希表中所有字段和值。
HMGET：获取哈希表中多个字段的值。
HKEYS：获取哈希表中所有字段。
HVALS：获取哈希表中所有值。

写命令：
HSET：为哈希表中的字段赋值。
HMSET：同时为哈希表中的多个字段赋值（已弃用，推荐使用 HSET）。
HDEL：删除哈希表中的一个或多个字段。
HINCRBY / HINCRBYFLOAT：对哈希表中指定字段的数值进行递增。

1.3 列表（Lists）
读命令：
LINDEX：通过索引获取列表中的元素。
LRANGE：获取列表中指定范围的元素。
LLEN：获取列表的长度。
LPOP / RPOP：移除并返回列表头部或尾部的元素。

写命令：
LPUSH / RPUSH：在列表头部或尾部添加元素。
LSET：通过索引设置列表元素的值。
LREM：根据值移除列表中的元素。
LTRIM：裁剪列表，使其只保留指定范围内的元素。

1.4 集合（Sets）
读命令：
SMEMBERS：获取集合中的所有成员。
SISMEMBER：判断成员是否在集合中。
SCARD：获取集合的基数（成员数量）。
SRANDMEMBER：随机获取集合中的一个或多个成员。

写命令：
SADD：向集合添加一个或多个成员。
SREM：移除集合中的一个或多个成员。
SMOVE：将成员从一个集合移动到另一个集合。
SPOP：随机移除并返回集合中的一个或多个成员。
SUNION、SINTER、SDIFF：集合的并集、交集、差集操作。

1.5 有序集合（Sorted Sets）
读命令：
ZRANGE：获取有序集合中指定范围的成员。
ZREVRANGE：逆序获取指定范围的成员。
ZSCORE：获取有序集合中指定成员的分数。
ZRANK / ZREVRANK：获取成员的排名。
ZCARD：获取有序集合的基数。
ZRANGEBYSCORE：按分数范围获取成员。

写命令：
ZADD：向有序集合添加一个或多个成员，或者更新已存在成员的分数。
ZREM：移除有序集合中的一个或多个成员。
ZINCRBY：对有序集合中指定成员的分数进行递增。
ZINTERSTORE、ZUNIONSTORE：有序集合的交集、并集存储操作。

1.6 键操作（Keys）
读命令：
EXISTS：检查键是否存在。
TTL：获取键的剩余生存时间（秒）。
PTTL：获取键的剩余生存时间（毫秒）。

写命令：
DEL：删除一个或多个键。
EXPIRE / PEXPIRE：为键设置过期时间。
RENAME / RENAMENX：重命名键。
TYPE：获取键的数据类型。

1.7 其他命令
只读命令：
PING：测试连接是否可用。
AUTH：进行身份验证。
INFO：获取服务器信息和统计数据。
管理类命令（需要谨慎授予权限）：

ACL：管理访问控制列表。
CONFIG：获取或修改服务器配置。
SAVE / BGSAVE：同步或异步保存数据到磁盘。
FLUSHALL / FLUSHDB：删除所有数据库中的所有键或指定数据库的所有键
*/

var (
	redisReadCmds = []string{
		"GET", "MGET",
		"HGET", "HGETALL", "HMGET",
		"LRANGE", "LINDEX",
		"SMEMBERS", "SISMEMBER",
		"ZRANGE", "ZREVRANGE", "ZSCORE",
		"EXISTS", "TTL", "PTTL",
	}
	redisWriteCmds = []string{
		"SET", "MSET",
		"HSET",
		"LPUSH", "RPUSH", "LSET", "LREM",
		"SADD", "SREM", "SUNION", "SINTER", "SDIFF",
		"ZADD", "ZREM", "ZINTERSTORE", "ZUNIONSTORE",
		"DEL", "EXPIRE", "PEXPIRE", "RENAME",
	}

	redisGrantOp  = "grant"
	redisRevokeOp = "revoke"
)

func redisGetReadWriteArgs(readKeyPrefixList []string, writeKeyPrefixList []string, operation string) (permissionArgs []string, err error) {
	permissionArgs = []string{}

	var operator string
	if operation == redisGrantOp {
		operator = "+"
	} else if operation == redisRevokeOp {
		operator = "-"
	} else {
		err = fmt.Errorf("operation should be [%s, %s]", redisGrantOp, redisRevokeOp)
		return
	}

	readKeyPrefixMap := make(map[string]struct{})
	for _, keyPrefix := range readKeyPrefixList {
		if keyPrefix != "" {
			readKeyPrefixMap[keyPrefix] = struct{}{}
		}
	}
	if len(readKeyPrefixMap) > 0 {
		for keyPrefix := range readKeyPrefixMap {
			for _, cmd := range redisReadCmds {
				permissionArgs = append(permissionArgs, operator+cmd+"|"+keyPrefix+"*")
			}
			permissionArgs = append(permissionArgs, "~"+keyPrefix+"*")
		}
	}

	writeKeyPrefixMap := make(map[string]struct{})
	for _, keyPrefix := range writeKeyPrefixList {
		if keyPrefix != "" {
			writeKeyPrefixMap[keyPrefix] = struct{}{}
		}
	}
	if len(writeKeyPrefixMap) > 0 {
		for keyPrefix := range writeKeyPrefixMap {
			for _, cmd := range redisWriteCmds {
				permissionArgs = append(permissionArgs, operator+cmd+"|"+keyPrefix+"*")
			}
			permissionArgs = append(permissionArgs, "~"+keyPrefix+"*")
		}
	}
	return
}

// 创建redis用户
func redisCreateUser(host, port, adminUser, adminPassword, userName, password string, userReadKeyPrefix, userWriteKeyPrefix []string) (err error) {
	args := []string{
		"-h", host,
		"-p", port,
		"-u", adminUser,
		"-a", adminPassword,
		"ACL", "SETUSER", userName,
		"on", ">" + password,
	}

	permissionArgs, tmpErr := redisGetReadWriteArgs(userReadKeyPrefix, userWriteKeyPrefix, "grant")
	if tmpErr != nil {
		err = fmt.Errorf("redis get read write args failed:%s", tmpErr.Error())
		return
	}
	if len(permissionArgs) > 0 {
		args = append(args, permissionArgs...)
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

// 授予redis用户读/写权限
func redisGrantReadWritePermission(host, port, adminUser, adminPassword, userName string, userReadKeyPrefix, userWriteKeyPrefix []string) (err error) {
	args := []string{
		"-h", host,
		"-p", port,
		"-u", adminUser,
		"-a", adminPassword,
		"ACL", "SETUSER", userName,
	}

	permissionArgs, tmpErr := redisGetReadWriteArgs(userReadKeyPrefix, userWriteKeyPrefix, "grant")
	if tmpErr != nil {
		err = fmt.Errorf("redis get read write args failed:%s", tmpErr.Error())
		return
	}
	if len(permissionArgs) > 0 {
		args = append(args, permissionArgs...)
	} else {
		// 无读/写权限需要授予
		return
	}

	output, tmpErr := runRedisCli(args...)
	if tmpErr != nil {
		err = fmt.Errorf("output:%s, error:%s", output, tmpErr.Error())
		return
	}
	return
}

// 撤销redis用户读/写权限
func redisRevokeReadWritePermission(host, port, adminUser, adminPassword, userName string, userReadKeyPrefix, userWriteKeyPrefix []string) (err error) {
	args := []string{
		"-h", host,
		"-p", port,
		"-u", adminUser,
		"-a", adminPassword,
		"ACL", "SETUSER", userName,
	}

	permissionArgs, tmpErr := redisGetReadWriteArgs(userReadKeyPrefix, userWriteKeyPrefix, "revoke")
	if tmpErr != nil {
		err = fmt.Errorf("redis get read write args failed:%s", tmpErr.Error())
		return
	}
	if len(permissionArgs) > 0 {
		args = append(args, permissionArgs...)
	} else {
		// 无读/写权限需要撤销
		return
	}

	output, tmpErr := runRedisCli(args...)
	if tmpErr != nil {
		err = fmt.Errorf("output:%s, error:%s", output, tmpErr.Error())
		return
	}
	return
}
