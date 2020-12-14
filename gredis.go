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
	RedisClient = redis.Client
	Options     = redis.Options
	BoolCmd     = redis.BoolCmd
	Cmd         = redis.Cmd
	Z           = redis.Z
	ScanCmd     = redis.ScanCmd
)

type Client struct {
	RedisClient
}

// keys
func (c *Client) ScanIter(pattern string, count int64) *ScansCmd {
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

// 查找集合元素
func (c *Client) SScanIter(key, match string, count int64) *ScansCmd {
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

// 查找zset中的元素
func (c *Client) ZScanIter(key, match string, count int64) *ScansCmd {
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

// 查找hash中的元素
func (c *Client) HScanIter(key, match string, count int64) *ScansCmd {
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

// HSet
func (c *Client) HSet(key, field string, value interface{}, expiration time.Duration) *Cmd {
	keys := []string{key}
	ex := gcommon.DurationToIntSecond(expiration)
	return c.Eval(hsetScript, keys, ex, field, value)
}

// HMSet
func (c *Client) HMSet(key string, fields map[string]interface{}, expiration time.Duration) *Cmd {
	keys := []string{key}
	ex := gcommon.DurationToIntSecond(expiration)
	args := []interface{}{ex}
	for k, v := range fields {
		if v != nil {
			args = append(args, k, v)
		}
	}
	return c.Eval(hmsetScript, keys, args...)
}

/*
向zset中插入成员并剪切，并截取只保留分数最高的length个成员
*/
func (c *Client) ZAddRemByRank(key string, length int64, members ...Z) *Cmd {
	keys := []string{key}
	args := []interface{}{0, -(length + 1)}
	for _, member := range members {
		args = append(args, member.Score, member.Member)
	}
	return c.Eval(zaddRemByRankScript, keys, args...)
}

/*
从左边向list插入元素，并截取只保留左起length个元素
*/
func (c *Client) LPushTrim(key string, length int64, values ...interface{}) *Cmd {
	keys := []string{key}
	args := []interface{}{0, length - 1}
	args = append(args, values...)
	return c.Eval(lpushTrimScript, keys, args...)
}

/*
从右边向list插入元素，并截取只保留右起length个元素
*/
func (c *Client) RPushTrim(key string, length int64, values ...interface{}) *Cmd {
	keys := []string{key}
	args := []interface{}{-length, -1}
	args = append(args, values...)
	return c.Eval(rpushTrimScript, keys, args...)
}

/*
通过配置生成客户端
*/
func NewClient(opt *Options) *Client {
	return &Client{
		RedisClient: *redis.NewClient(opt),
	}
}

/*
通过reids客户端实例生成客户端
*/
func NewClientFromRedisClient(client *RedisClient) *Client {
	return &Client{
		RedisClient: *client,
	}
}
