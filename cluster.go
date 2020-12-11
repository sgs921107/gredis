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
	RedisClusterClient = redis.ClusterClient
	ClusterOptions     = redis.ClusterOptions
)

type ClusterClient struct {
	RedisClusterClient
}

// keys
func (c *ClusterClient) ScanIter(pattern string, count int64) (keys []string, err error) {
	cursor := uint64(0)
	for {
		page, cursor, err := c.Scan(cursor, pattern, count).Result()
		if err != nil {
			return keys, err
		}
		keys = append(keys, page...)
		if cursor == 0 {
			return keys, nil
		}
	}
}

// 查找集合元素
func (c *ClusterClient) SScanIter(key, match string, count int64) (members []string, err error) {
	cursor := uint64(0)
	for {
		page, cursor, err := c.SScan(key, cursor, match, count).Result()
		if err != nil {
			return members, err
		}
		members = append(members, page...)
		if cursor == 0 {
			return members, nil
		}
	}
}

// 查找zset中的元素
func (c *ClusterClient) ZScanIter(key, match string, count int64) (members []string, err error) {
	corsur := uint64(0)
	for {
		page, cursor, err := c.ZScan(key, corsur, match, count).Result()
		if err != nil {
			return members, err
		}
		members = append(members, page...)
		if cursor == 0 {
			return members, nil
		}
	}
}

// 查找hash中的元素
func (c *ClusterClient) HScanIter(key, match string, count int64) (members []string, err error) {
	cursor := uint64(0)
	for {
		page, cursor, err := c.HScan(key, cursor, match, count).Result()
		if err != nil {
			return members, err
		}
		members = append(members, page...)
		if cursor == 0 {
			return members, nil
		}
	}
}

// HSet
func (c *ClusterClient) HSet(key, field string, value interface{}, expiration time.Duration) *Cmd {
	keys := []string{key}
	ex := gcommon.DurationToIntSecond(expiration)
	return c.Eval(hsetScript, keys, ex, field, value)
}

// HMSet
func (c *ClusterClient) HMSet(key string, fields map[string]interface{}, expiration time.Duration) *Cmd {
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
func (c *ClusterClient) ZAddRemByRank(key string, length int, members ...Z) *Cmd {
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
func (c *ClusterClient) LPushTrim(key string, length int, values ...interface{}) *Cmd {
	keys := []string{key}
	args := []interface{}{0, length - 1}
	args = append(args, values...)
	return c.Eval(lpushTrimScript, keys, args...)
}

/*
从右边向list插入元素，并截取只保留右起length个元素
*/
func (c *ClusterClient) RPushTrim(key string, length int, values ...interface{}) *Cmd {
	keys := []string{key}
	args := []interface{}{-length, -1}
	args = append(args, values...)
	return c.Eval(rpushTrimScript, keys, args...)
}

func NewClusterClient(opt *ClusterOptions) *ClusterClient {
	return &ClusterClient{
		RedisClusterClient: *redis.NewClusterClient(opt),
	}
}
