package models

import (
	fmt "fmt"
	utils "vulnlabs-rest-api/utils"

	"github.com/gomodule/redigo/redis"
	redisgo "github.com/gomodule/redigo/redis"
)

const (
	RedisSessionStoragePrefix       = "session"
	RedisSessionStorageUserIDSuffix = "userID"
	RedisUserStoragePrefix          = "user"
	RedisUserStorageSessionSuffix   = "session"
)

// RedisInterface : Redis Communication interface
type RedisInterface interface {
	CloseConnection() error
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	SetWithExpiration(key string, value []byte, expirationInSeconds int) error
	Exists(key string) (bool, error)
	Delete(key string) error
	Incr(counterKey string) (int, error)
	Multi(commands []RedisCommand) ([]interface{}, error)
}

// Redis : Redis communication interface
type Redis struct {
	Connection redis.Conn
}

// RedisCommand : Redis command struct
type RedisCommand struct {
	Command string
	Args    []interface{}
}

// NewRedis : Return a new Redis abstraction struct
func NewRedis(connectionURL string, password string) *Redis {

	// Initialize the redis connection to a redis instance running on your local machine
	conn, err := redisgo.DialURL(connectionURL)

	if err != nil {
		utils.PanicOnError(err, "Failed to connect to Redis")
	}

	// Authenticate to Redis
	conn.Do("AUTH", password)

	// Return new MongoDB abstraction struct
	return &Redis{
		Connection: conn,
	}
}

// CloseConnection : Close Redis Connection
func (redis *Redis) CloseConnection() error {

	return redis.Connection.Close()
}

func (redis *Redis) Get(key string) ([]byte, error) {

	var data []byte

	data, err := redisgo.Bytes(redis.Connection.Do("GET", key))

	if err != nil {
		return nil, fmt.Errorf("error getting key %s : %v", key, err)
	}
	return data, nil
}

func (redis *Redis) HGet(key string, field string) ([]byte, error) {

	var data []byte
	data, err := redisgo.Bytes(redis.Connection.Do("HGET", key, field))

	if err != nil {
		return nil, fmt.Errorf("error getting key %s : %v", key, err)
	}
	return data, nil
}

func (redis *Redis) HSet(key string, field1 string, value1 []byte, field2 string, value2 []byte) error {

	_, err := redis.Connection.Do("HSET", key, field1, value1, field2, value2)
	if err != nil {
		return fmt.Errorf("error setting key %s to %s : %v", key, value1, err)
	}
	return nil
}

func (redis *Redis) SetWithExpiration(key string, value []byte, expirationInSeconds int) error {

	_, err := redis.Connection.Do("SETEX", key, fmt.Sprintf("%d", expirationInSeconds), value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error setting key %s to %s : %v", key, v, err)
	}
	return nil
}

func (redis *Redis) Set(key string, value []byte) error {

	_, err := redis.Connection.Do("SET", key, value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error setting key %s to %s : %v", key, v, err)
	}
	return nil
}

func (redis *Redis) Rename(oldKey string, newKey string) error {

	_, err := redis.Connection.Do("RENAME", oldKey, newKey)
	if err != nil {
		return fmt.Errorf("error renaming key %s to %s : %v", oldKey, newKey, err)
	}
	return nil
}

func (redis *Redis) Exists(key string) (bool, error) {

	ok, err := redisgo.Bool(redis.Connection.Do("EXISTS", key))
	if err != nil {
		return ok, fmt.Errorf("error checking if key %s exists : %v", key, err)
	}
	return ok, nil
}

func (redis *Redis) Delete(key string) error {

	_, err := redis.Connection.Do("DEL", key)

	if err != nil {
		return err
	}

	return nil
}

func (redis *Redis) GetKeys(pattern string) ([]string, error) {

	iter := 0
	keys := []string{}
	for {
		arr, err := redisgo.Values(redis.Connection.Do("SCAN", iter, "MATCH", pattern))
		if err != nil {
			return keys, fmt.Errorf("error retrieving '%s' keys", pattern)
		}

		iter, _ = redisgo.Int(arr[0], nil)
		k, _ := redisgo.Strings(arr[1], nil)
		keys = append(keys, k...)

		if iter == 0 {
			break
		}
	}

	return keys, nil
}

func (redis *Redis) Incr(counterKey string) (int, error) {

	return redisgo.Int(redis.Connection.Do("INCR", counterKey))
}

func (redis *Redis) Multi(commands []RedisCommand) ([]interface{}, error) {

	redis.Connection.Send("MULTI")

	for _, cmd := range commands {
		redis.Connection.Send(cmd.Command, cmd.Args...)
	}

	r, err := redisgo.Values(redis.Connection.Do("EXEC"))

	if err != nil {
		return nil, err
	}

	return r, nil
}
