package handler

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var IPLimitRequestNumber = 10

type RateLimitHandler interface {
	CheckIPLimit(*gin.Context)
}

type ratelimitHandler struct {
	Rdb *redis.Client
}

func NewRateLimitHandler() RateLimitHandler {
	return &ratelimitHandler{
		Rdb: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
}

func (rlh *ratelimitHandler) CheckIPLimit(c *gin.Context) {
	clientIP := c.ClientIP()
	data, _ := rlh.Rdb.HGetAll(c, clientIP).Result()
	if len(data) == 0 {
		newdata := make(map[string]interface{})
		newdata["remain"] = IPLimitRequestNumber - 1
		newdata["deadline"] = time.Now().Add(time.Hour * 1).Unix()

		if err := rlh.Rdb.HMSet(c, clientIP, newdata).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Header("X-RateLimit-Remaining", strconv.Itoa(newdata["remain"].(int)))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(newdata["deadline"].(int64), 10))
		return
	}

	remain, _ := strconv.ParseInt(data["remain"], 10, 64)
	deadline, _ := strconv.ParseInt(data["deadline"], 10, 64)
	if deadline < time.Now().Unix() {
		newdata := make(map[string]interface{})
		newdata["remain"] = IPLimitRequestNumber - 1
		newdata["deadline"] = time.Now().Add(time.Hour * 1).Unix()

		if err := rlh.Rdb.HMSet(c, clientIP, newdata).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Header("X-RateLimit-Remaining", newdata["remain"].(string))
		c.Header("X-RateLimit-Reset", newdata["deadline"].(string))
		return
	}

	if remain <= 0 {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests"})
		return
	}

	remain -= 1
	if err := rlh.Rdb.HSet(c, clientIP, "remain", remain).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("X-RateLimit-Remaining", strconv.FormatInt(remain, 10))
	c.Header("X-RateLimit-Reset", data["deadline"])
}
