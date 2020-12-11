/*************************************************************************
	> File Name: test/gredis_test.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月11日 星期五 16时29分22秒
*************************************************************************/

package main

import (
	"github.com/sgs921107/gredis"
	"strconv"
	"testing"
	"time"
)

var (
	addr     = "172.17.0.1:6379"
	db       = 0
	password = "qaz123"
)

var ex = time.Second * 30
var c = gredis.NewClient(&gredis.Options{
	Addr:     addr,
	DB:       db,
	Password: password,
})

func TestKeys(t *testing.T) {
	if _, err := c.ScanIter("*", 1000).Result(); err != nil {
		t.Errorf("redis keys cmd err: %s", err.Error())
	}
}

func TestSet(t *testing.T) {
	key := "test_string"
	val := "hello"
	if _, err := c.Set(key, val, ex).Result(); err != nil {
		t.Errorf("redis set cmd err: %s", err.Error())
	}
	if ret, err := c.Get(key).Result(); ret != val {
		t.Errorf(`c.Get("%s").Result() == ("%v", %v), want ("%v", nil)`, key, ret, err, val)
	}
}

func TestHSet(t *testing.T) {
	key := "test_hash1"
	field := "hello"
	val := "go"
	if _, err := c.HSet(key, field, val, ex).Result(); err != nil {
		t.Errorf("redis hset cmd err: %s", err.Error())
	}
	if ret, err := c.HGet(key, field).Result(); ret != val {
		t.Errorf(`c.HGet("%s", "%s").Result() == ("%v", %v), want ("%v", nil)`, key, field, ret, err, val)
	}
	if ttl, err := c.TTL(key).Result(); ttl == -1 {
		t.Errorf(`c.TTL("%s").Result() == (-1, %v) want (uint, nil)`, key, err)
	}
}

func TestHMset(t *testing.T) {
	key := "test_hash2"
	fields := map[string]interface{}{
		"name": "tom",
		"age":  18,
		"sex":  nil,
	}
	if _, err := c.HMSet(key, fields, ex).Result(); err != nil {
		t.Errorf("redis hmset cmd err: %s", err.Error())
	}
	if ret, err := c.HGet(key, "name").Result(); ret != "tom" {
		t.Errorf(`c.HGet("%s", "name").Result() == ("%v", %v), want ("tom", nil)`, key, ret, err)
	}
	if ret, err := c.HGet(key, "sex").Result(); ret != "" {
		t.Errorf(`c.HGet("%s", "sex").Result() == ("%v", %v), want ("", nil)`, key, ret, err)
	}
	if ret, err := c.HGet(key, "age").Result(); ret != "18" {
		t.Errorf(`c.HGet("%s", "age").Result() == ("%v", %v), want ("18", nil)`, key, ret, err)
	}
	if ttl, err := c.TTL(key).Result(); ttl == -1 {
		t.Errorf(`c.TTL("%s").Result() == (-1, %v) want (uint, nil)`, key, err)
	}
}

func TestZAddRemByRank(t *testing.T) {
	key := "test_zadd"
	length := int64(10)
	var members []gredis.Z
	for i := 0; i < 20; i++ {
		member := gredis.Z{
			Score:  float64(i),
			Member: i * 2,
		}
		members = append(members, member)
	}
	if _, err := c.ZAddRemByRank(key, length, members...).Result(); err != nil {
		t.Errorf("redis ZAddRemByRank cmd err: %s", err.Error())
	}
	if ret, err := c.ZCard(key).Result(); ret != length {
		t.Errorf(`c.ZCard("%s").Result() == (%d, %v), want (%d, nil)`, key, ret, err, length)
	}
	if ret, err := c.ZScore(key, "10").Result(); err == nil {
		t.Errorf(`c.ZScore("%s", "10").Result() == (%v, nil), want (0, err)`, key, ret)
	}
	if ret, err := c.ZScore(key, "20").Result(); ret != float64(10) {
		t.Errorf(`c.ZScore("%s", "20").Result() == (%v, %v), want (10.0, nil)`, key, ret, err)
	}
	c.Expire(key, ex)
}

func TestPushTrim(t *testing.T) {
	lpushKey := "test_list1"
	var length int64 = 10
	var num int64 = 20
	var list = make([]interface{}, num)
	for i, _ := range list {
		list[i] = i
	}
	// 测试lpushtrim
	if _, err := c.LPushTrim(lpushKey, length, list...).Result(); err != nil {
		t.Errorf("redis LPushTrim cmd err: %s", err.Error())
	}
	if ret, err := c.LLen(lpushKey).Result(); ret != length {
		t.Errorf(`c.LPushTrim("%s").Result() == (%d, %v), want (%d, nil)`, lpushKey, ret, err.Error(), length)
	}
	expectRet := int(num - 1)
	if ret, err := c.LIndex(lpushKey, 0).Result(); ret != strconv.Itoa(expectRet) {
		t.Errorf(`c.LIndex("%s").Result() == ("%s", %v), want ("%d", nil)`, lpushKey, ret, err, expectRet)
	}
	expectRet = int(num - length)
	if ret, err := c.LIndex(lpushKey, length-1).Result(); ret != strconv.Itoa(expectRet) {
		t.Errorf(`c.LIndex("%s").Result() == ("%s", %v), want ("%d", nil)`, lpushKey, ret, err, expectRet)
	}
	rpushKey := "test_list2"
	// 测试rpushtrim
	if _, err := c.RPushTrim(rpushKey, length, list...).Result(); err != nil {
		t.Errorf("redis RPushTrim cmd err: %s", err.Error())
	}
	if ret, err := c.LLen(rpushKey).Result(); ret != length {
		t.Errorf(`c.RPushTrim("%s").Result() == (%d, %v), want (%d, nil)`, lpushKey, ret, err.Error(), length)
	}
	expectRet = int(num - length)
	if ret, err := c.LIndex(rpushKey, 0).Result(); ret != strconv.Itoa(expectRet) {
		t.Errorf(`c.LIndex("%s").Result() == ("%s", %v), want ("%d", nil)`, lpushKey, ret, err, expectRet)
	}
	expectRet = int(num - 1)
	if ret, err := c.LIndex(rpushKey, int64(length-1)).Result(); ret != strconv.Itoa(expectRet) {
		t.Errorf(`c.LIndex("%s").Result() == ("%s", %v), want ("%d", nil)`, lpushKey, ret, err, expectRet)
	}
	c.Expire(lpushKey, ex)
	c.Expire(rpushKey, ex)
}
