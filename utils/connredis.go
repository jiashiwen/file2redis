package utils

import (
	// "fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

type PoolConf struct {
	MaxIdle int
	Server  string
}

func GetPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		MaxActive:   10000,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialTimeout("tcp", server, time.Second, time.Second, time.Second)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func GetConnection(pool *redis.Pool) redis.Conn {
	return pool.Get()
}

/*func main() {

	// conn, _ := redis.DialTimeout("tcp", "10.100.21.90:6379", 0, 1*time.Second, 1*time.Second)
	pool := GetPool("192.168.3.6:6379")

	conn := GetConnection(pool)
	size, err := conn.Do("DBSIZE")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("size is %d \n", size)
	fmt.Println(pool.ActiveCount())
	conn.Close()
}*/
