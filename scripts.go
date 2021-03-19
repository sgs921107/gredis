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
    local unpack = unpack or table.unpack
    local ret = redis.call("hmset", name, unpack(ARGV, 2))
    redis.call("expire", name, expire_time)
    return ret
`

var zaddRemByRankScript = `
    local name = KEYS[1]
    local startNum = ARGV[1]
    local stopNum = ARGV[2]
    local unpack = unpack or table.unpack
    local ret = redis.call("zadd", name, unpack(ARGV, 3))
    redis.call("zremrangebyrank", name, startNum, stopNum)
    return ret
`

var lpushTrimScript = `
    local name = KEYS[1]
    local startNum = ARGV[1]
    local stopNum = ARGV[2]
    local unpack = unpack or table.unpack
    local ret = redis.call("lpush", name, unpack(ARGV, 3))
    if( ret > stopNum + 1 )
    then
        redis.call("ltrim", name, startNum, stopNum)
    end
    return ret
`

var rpushTrimScript = `
    local name = KEYS[1]
    local startNum = ARGV[1]
    local stopNum = ARGV[2]
    local unpack = unpack or table.unpack
    local ret = redis.call("rpush", name, unpack(ARGV, 3))
    if( ret > stopNum + 1 )
    then
        redis.call("ltrim", name, startNum, stopNum)
    end
    return ret
`

var saddex_script = `
    local name = KEYS[1]
    local expire_time = ARGV[1]
    local unpack = unpack or table.unpack
    local ret = redis.call("sadd", name, unpack(ARGV, 2))
    redis.call("expire", name, expire_time)
    return ret
`
