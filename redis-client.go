package main

import (
    "context"

    "github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func Client(clientAddr string) redis.Client {
    return *redis.NewClient(&redis.Options{
        Addr:    clientAddr,
        Password: "", // no password set
        DB:       0,  // use default DB
    })
}


func PutInRedis(clientAddr, key, value string) {
    err := Client(clientAddr).Set(ctx, key, value, 0).Err()
    if err != nil {
        panic(err)
    }
}

func GetFromRedis(clientAddr, key string) (string, error) {
    val, err := Client(clientAddr).Get(ctx, key).Result()
    if err == redis.Nil {
        return "", nil
    } else if err != nil {
        return "", err
    } else {
        return val, nil
    }
}
