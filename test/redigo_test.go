package test

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

func init_conn() {
	pool = &redis.Pool{
		MaxIdle:     30,
		MaxActive:   30,
		IdleTimeout: time.Duration(time.Second),
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", "192.168.1.44:6377",
				redis.DialConnectTimeout(time.Second),
				redis.DialReadTimeout(time.Second),
				redis.DialWriteTimeout(time.Second),
				redis.DialPassword(""),
			)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
}

func TestRedigo(t *testing.T) {
	init_conn()
	t.Run("rmslots", func(t *testing.T) {
		rmslots()
	})
	t.Run("do", func(t *testing.T) {
		do()
	})
	t.Run("rmblackinfo", func(t *testing.T) {
		rmblackinfo()
	})
	t.Run("rmslots", func(t *testing.T) {
		rmslots()
	})
	t.Run("pipeline", func(t *testing.T) {
		pipeline()
	})
}

func rmslots() {
	conn := pool.Get()
	defer conn.Close()
	cursor := 0
	for {
		fmt.Println(cursor)
		_reply, err := conn.Do("SCAN", cursor, "MATCH", "15:slot:*:data*", "COUNT", 100)
		if err != nil {
			fmt.Println(err)
			break
		}
		reply, ok := _reply.([]interface{})
		if !ok {
			fmt.Println("to []interface{} failed")
			break
		}

		cursor, _ = strconv.Atoi(string(reply[0].([]uint8)))

		result := reply[1].([]interface{})
		fmt.Println(cursor, len(result))
		for i := 0; i < len(result); i++ {
			key := string(result[i].([]uint8))
			_, err := conn.Do("Del", key)
			fmt.Println(key)
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println("del", key, "success")
		}
		if cursor == 0 {
			fmt.Println("迭代结束")
			break
		}
	}
}

func do() {
	conn := pool.Get()
	defer conn.Close()
	reply, _ := redis.Strings(conn.Do("ZRANGE", "13:autoblack", 1, -1))
	for i := range reply {
		_, err := conn.Do("ZREM", "13:autoblack", reply[i])
		if err != nil {
			fmt.Println(err)
			break
		}
		_, err = conn.Do("Del", fmt.Sprintf("13:autoblackinfo:%s", reply))
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(reply[i], "数据删除成功")
	}
}

func rmblackinfo() {
	conn := pool.Get()
	defer conn.Close()
	cursor := 0
	for {
		fmt.Println(cursor)
		_reply, err := conn.Do("SCAN", cursor, "MATCH", "15:autoblackinfo:*", "COUNT", 100)
		if err != nil {
			fmt.Println(err)
			break
		}
		reply, ok := _reply.([]interface{})
		if !ok {
			fmt.Println("to []interface{} failed")
			break
		}

		cursor, _ = strconv.Atoi(string(reply[0].([]uint8)))

		result := reply[1].([]interface{})
		fmt.Println(cursor, len(result))
		for i := 0; i < len(result); i++ {
			key := string(result[i].([]uint8))
			_, err := conn.Do("Del", key)
			fmt.Println(key)
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println("del", key, "success")
		}
		if cursor == 0 {
			fmt.Println("迭代结束")
			break
		}
	}
}

func key_mutex() {
	wg := sync.WaitGroup{}
	wg.Add(5)
	for i := 1; i < 5; i++ {
		go func(idx int) {

			conn := pool.Get()
			defer conn.Close()
			for {
				ok, err := redis.Bool(conn.Do("SETNX", "key_mutex", strconv.Itoa(idx)))
				if err != nil {
					panic(err)
				}
				if !ok {
					fmt.Println(idx, "设置排他锁失败，等待排它锁释放")
					time.Sleep(time.Second)
				} else {
					fmt.Println(idx, "数据库load数据")
					time.Sleep(time.Second)
					fmt.Println(idx, "数据库load数据成功")
					ok, err = redis.Bool(conn.Do("DEL", "key_mutex"))
					if err != nil {
						panic(err)
					}
					if ok {
						fmt.Println(idx, "删除缓存成功")
					}
					wg.Done()
					break
				}
			}
		}(i)
	}
	wg.Wait()
}

// 虽然是pipeline，redis中的实现是是使用send、flush、receive来实现pipeline的功能
func pipeline() {
	conn := pool.Get()
	defer conn.Close()
	for i := 0; i < 10; i++ {
		err := conn.Send("SET", fmt.Sprintf("myname:%d", i), i)
		fmt.Println(err)
	}
	err := conn.Flush()
	fmt.Println(err)
	for i := 0; i < 10; i++ {
		reply, _ := conn.Receive()
		fmt.Println(reply)
	}

	for i := 0; i < 10; i++ {
		err := conn.Send("GET", fmt.Sprintf("myname:%d", i))
		fmt.Println(err)
	}
	err = conn.Flush()
	fmt.Println(err)
	for i := 0; i < 10; i++ {
		reply, _ := redis.Int(conn.Receive())
		fmt.Println(reply)
	}

	for i := 0; i < 10; i++ {
		err := conn.Send("DLE", fmt.Sprintf("myname:%d", i))
		fmt.Println(err)
	}
	err = conn.Flush()
	fmt.Println(err)
	for i := 0; i < 10; i++ {
		reply, _ := redis.Int(conn.Receive())
		fmt.Println(reply)
	}
}
