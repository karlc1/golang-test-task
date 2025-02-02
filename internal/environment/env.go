package environment

import (
	"fmt"
	"os"
	"strconv"
)

type Environment struct {
	RMQqueueName string
	RMQuser      string
	RMQpassword  string
	RMQhost      string
	RMQport      int
	RMQexchange  string
	RedisHost    string
	RedisPort    int
	ApiPort      int
}

func MustGetApiEnv() Environment {
	return Environment{
		RMQqueueName: mustGetString("RMQ_QUEUE_NAME"),
		RMQuser:      mustGetString("RMQ_USER"),
		RMQpassword:  mustGetString("RMQ_PASS"),
		RMQhost:      mustGetString("RMQ_HOST"),
		RMQport:      mustGetInt("RMQ_PORT"),
		ApiPort: mustGetInt("API_PORT"),
	}
}

func MustGetWorkerEnv() Environment {
	return Environment{
		RMQqueueName: mustGetString("RMQ_QUEUE_NAME"),
		RMQuser:      mustGetString("RMQ_USER"),
		RMQpassword:  mustGetString("RMQ_PASS"),
		RMQhost:      mustGetString("RMQ_HOST"),
		RMQport:      mustGetInt("RMQ_PORT"),
		RMQexchange:  mustGetString("RMQ_EXCHANGE"),
		RedisHost:    mustGetString("REDIS_HOST"),
		RedisPort:    mustGetInt("REDIS_PORT"),
	}
}

func MustGetReporterEnv() Environment {
	return Environment{
		RedisHost: mustGetString("REDIS_HOST"),
		RedisPort: mustGetInt("REDIS_PORT"),
		ApiPort: mustGetInt("API_PORT"),
	}
}

func mustGetString(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("Missing required env var '%s'", key))
	}
	return val
}

func mustGetInt(key string) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("Missing required env var '%s'", key))
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Sprintf("malformed int env var: '%s': %v", key, err))
	}
	return i
}
