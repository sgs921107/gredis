/*************************************************************************
	> File Name: gredis.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月11日 星期五 10时17分29秒
*************************************************************************/
/*
对redis client进行扩展
添加和重写一些方法
*/

package gredis

import (
	"github.com/go-redis/redis"
	"github.com/sgs921107/gcommon"
	"time"
)

// 起别名
type (
	// RedisClusterClient redis集群客户端
	RedisClusterClient = redis.ClusterClient
	// ClusterOptions  创建集群客户端的参数结构类型
	ClusterOptions     = redis.ClusterOptions
)

/*
ClusterClient  集群客户端结构类型
*/
type ClusterClient struct {
	RedisClusterClient
}

/*
ScanIter 获取匹配指定pattern的所有redis key  替代keys方法
*/
func (c *ClusterClient) ScanIter(pattern string, count int64) *ScansCmd {
	scansCmd := NewScansCmd()
	cursor := uint64(0)
	for {
		scanCmd := c.Scan(cursor, pattern, count)
		scansCmd.addScanCmd(scanCmd)
		_, cursor, err := scanCmd.Result()
		if err != nil || cursor == 0 {
			return scansCmd
		}
	}
}

/*
SScanIter 获取集合中匹配指定pattern的所有元素  替代sismembers
*/
func (c *ClusterClient) SScanIter(key, match string, count int64) *ScansCmd {
	scansCmd := NewScansCmd()
	cursor := uint64(0)
	for {
		scanCmd := c.SScan(key, cursor, match, count)
		scansCmd.addScanCmd(scanCmd)
		_, cursor, err := scanCmd.Result()
		if err != nil || cursor == 0 {
			return scansCmd
		}
	}
}

/*
ZScanIter 获取有序集合中匹配指定pattern的所有元素
*/
func (c *ClusterClient) ZScanIter(key, match string, count int64) *ScansCmd {
	scansCmd := NewScansCmd()
	cursor := uint64(0)
	for {
		scanCmd := c.ZScan(key, cursor, match, count)
		scansCmd.addScanCmd(scanCmd)
		_, cursor, err := scanCmd.Result()
		if err != nil || cursor == 0 {
			return scansCmd
		}
	}
}

/*
HScanIter 获取字典中匹配指定pattern的所有字段
*/
func (c *ClusterClient) HScanIter(key, match string, count int64) *ScansCmd {
	scansCmd := NewScansCmd()
	cursor := uint64(0)
	for {
		scanCmd := c.HScan(key, cursor, match, count)
		scansCmd.addScanCmd(scanCmd)
		_, cursor, err := scanCmd.Result()
		if err != nil || cursor == 0 {
			return scansCmd
		}
	}
}

/*
HSetEX 执行hset命令并设置过期时间 单位: 秒
*/
func (c *ClusterClient) HSetEX(key, field string, value interface{}, expiration time.Duration) *Cmd {
	keys := []string{key}
	ex := gcommon.DurationToIntSecond(expiration)
	return c.Eval(hsetexScript, keys, ex, field, value)
}

/*
HMSetEX 执行hmset命令并设置过期时间 单位: 秒
*/
func (c *ClusterClient) HMSetEX(key string, fields map[string]interface{}, expiration time.Duration) *Cmd {
	keys := []string{key}
	ex := gcommon.DurationToIntSecond(expiration)
	args := []interface{}{ex}
	for k, v := range fields {
		if v != nil {
			args = append(args, k, v)
		}
	}
	return c.Eval(hmsetexScript, keys, args...)
}

/*
ZAddRemByRank 向zset中插入成员并剪切，并截取只保留分数最高的length个成员
*/
func (c *ClusterClient) ZAddRemByRank(key string, length int64, members ...Z) *Cmd {
	keys := []string{key}
	args := []interface{}{0, -(length + 1)}
	for _, member := range members {
		args = append(args, member.Score, member.Member)
	}
	return c.Eval(zaddRemByRankScript, keys, args...)
}

/*
LPushTrim 从左边向list插入元素，并截取只保留左起length个元素
*/
func (c *ClusterClient) LPushTrim(key string, length int64, values ...interface{}) *Cmd {
	keys := []string{key}
	args := []interface{}{0, length - 1}
	args = append(args, values...)
	return c.Eval(lpushTrimScript, keys, args...)
}

/*
RPushTrim 从右边向list插入元素，并截取只保留右起length个元素
*/
func (c *ClusterClient) RPushTrim(key string, length int64, values ...interface{}) *Cmd {
	keys := []string{key}
	args := []interface{}{-length, -1}
	args = append(args, values...)
	return c.Eval(rpushTrimScript, keys, args...)
}

/*
NewClusterClient 实例化一个redis集群客户端
*/
func NewClusterClient(opt *ClusterOptions) *ClusterClient {
	return &ClusterClient{
		RedisClusterClient: *redis.NewClusterClient(opt),
	}
}

/*
NewClusterClientFromRedisClient 通过reids集群客户端实例生成客户端
*/
func NewClusterClientFromRedisClient(client *RedisClusterClient) *ClusterClient {
	return &ClusterClient{
		RedisClusterClient: *client,
	}
}
