package main

import (
	"fmt"
	"os"

	"github.com/dpordomingo/sre-counter/counter"
)

var envServerPort = os.Getenv("SERVER_PORT")
var envRedisHostPort = os.Getenv("REDIS_HOST_PORT")
var envInstanceName = os.Getenv("INSTANCE_NAME")

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		help()
		os.Exit(0)
	}

	if envInstanceName == "" {
		envInstanceName = "first_node"
	}

	if envServerPort == "" {
		envServerPort = "8090"
	}

	fmt.Printf("Redis server...\n%s\n\n", envRedisHostPort)
	var cacheClient = counter.NewRedisClient(envRedisHostPort)

	fmt.Println("Starting server...")
	var server = counter.NewServer(cacheClient, envServerPort, envInstanceName)

	if err := server.Run(); err != nil {
		fmt.Println(fmt.Errorf("error. server rised a panic. %s", err))
	}
}

func help() {
	fmt.Println(`
Run the server exposing the counter backed in Redis.

Environment variables:
- REDIS_HOST_PORT: hostname and port where Redis is listening (e.g. hostname:6379)
- SERVER_PORT:     port where the server will be listening (default: 8090)
- INSTANCE_NAME:   name for the server instance (default: first_node)

example:
$ REDIS_HOST_PORT="hostname:6379" \
  SERVER_PORT=9999 \
  INSTANCE_NAME="instance name" \
  redis-counter
`)
}
