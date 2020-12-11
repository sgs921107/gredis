/*************************************************************************
	> File Name: scripts.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月11日 星期五 10时47分32秒
*************************************************************************/

package gredis

var hsetScript = `
    local name = KEYS[1]
	local expire_time = ARGV[1]
	local key = ARGV[2]
	local value = ARGV[3]
	redis.call("hset", name, key, value)
	local ret = redis.call("expire", name, expire_time)
	return ret
`

var hmsetScript = `
	local name = KEYS[1]
    local expire_time = ARGV[1]
    local index = 2
    while index < table.getn(ARGV) do
        redis.call("hset", name, ARGV[index], ARGV[index + 1])
        index = index + 2
    end
	local ret = redis.call("expire", name, expire_time)
    return ret
`

var zaddRemByRankScript = `
	local name = KEYS[1]
    local startNum = ARGV[1]
    local stopNum = ARGV[2]
    local index = 3
    while index < table.getn(ARGV) do
        redis.call("zadd", name, ARGV[index], ARGV[index + 1])
        index = index + 2
    end
    local ret = redis.call("zremrangebyrank", name, startNum, stopNum)
    return ret
`

var lpushTrimScript = `
	local name = KEYS[1]
    local startNum = ARGV[1]
    local stopNum = ARGV[2]
    local index = 3
    while index <= table.getn(ARGV) do
        redis.call("lpush", name, ARGV[index])
        index = index + 1
    end
    local ret = redis.call("ltrim", name, startNum, stopNum)
    return ret
`

var rpushTrimScript = `
	local name = KEYS[1]
    local startNum = ARGV[1]
    local stopNum = ARGV[2]
    local index = 3
    while index <= table.getn(ARGV) do
        redis.call("rpush", name, ARGV[index])
        index = index + 1
    end
    local ret = redis.call("ltrim", name, startNum, stopNum)
    return ret
`
