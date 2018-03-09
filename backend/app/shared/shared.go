package shared

import(
	"github.com/go-redis/redis"
	"time"
)

const (
	SUCCESS = "0"
	CONTACT_ALREADY_EXISTS = "1"
	DATA_STORE_WRITE_ERR = "2"
	CONTACT_DOES_NOT_EXIST = "3"
	DATA_STORE_DELETE_ERR = "4"
	IMPL_ERR = "10"
)

//Api Command codes
const (
	POST = 1
	PUT = 2
)

type RedisClient interface {
	Get(key string) *redis.StringCmd
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(keys ...string) *redis.IntCmd
	Keys(pattern string) *redis.StringSliceCmd
}