/*************************************************************************
	> File Name: scripts.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月11日 星期五 10时47分32秒
*************************************************************************/

package gredis

var hsetexScript = `
    local name = KEYS[1]
    local expire_time = ARGV[1]
    local key = ARGV[2]
    local value = ARGV[3]
    local ret = redis.call("hset", name, key, value)
    redis.call("expire", name, expire_time)
    return ret
`

var hmsetexScript = `
    local name = KEYS[1]
    local expire_time = ARGV[1]
    local ret = redis.call("hmset", name, unpack(ARGV, 2))
    redis.call("expire", name, expire_time)
    return ret
`

// 向zset中加入元素并只保留指定数量的元素
var zaddRemByRankScript = `
    local name = KEYS[1]
    local startNum = ARGV[1]
    local stopNum = ARGV[2]
    local ret = redis.call("zadd", name, unpack(ARGV, 3))
    redis.call("zremrangebyrank", name, startNum, stopNum)
    return ret
`

// 从list左边截取若干元素
var lpushTrimScript = `
    local name = KEYS[1]
    local startNum = ARGV[1]
    local stopNum = ARGV[2]
    local ret = redis.call("lpush", name, unpack(ARGV, 3))
    if( ret > stopNum + 1 )
    then
        redis.call("ltrim", name, startNum, stopNum)
    end
    return ret
`

// 从list右边截取若干元素
var rpushTrimScript = `
    local name = KEYS[1]
    local startNum = ARGV[1]
    local stopNum = ARGV[2]
    local ret = redis.call("rpush", name, unpack(ARGV, 3))
    if( ret > stopNum + 1 )
    then
        redis.call("ltrim", name, startNum, stopNum)
    end
    return ret
`

// 向集合中插入数据并设置过期时间
var saddex_script = `
    local name = KEYS[1]
    local expire_time = ARGV[1]
    local ret = redis.call("sadd", name, unpack(ARGV, 2))
    redis.call("expire", name, expire_time)
    return ret
`

// 从list左端插入数据并设置过期时间
var lpushex_script = `
    local name = KEYS[1]
    local expire_time = ARGV[1]
    local ret = redis.call("lpush", name, unpack(ARGV, 2))
    redis.call("expire", name, expire_time)
    return ret
    `

// 从list右端插入数据并设置过期时间
var rpushex_script = `
    local name = KEYS[1]
    local expire_time = ARGV[1]
    local ret = redis.call("rpush", name, unpack(ARGV, 2))
    redis.call("expire", name, expire_time)
    return ret
    `
