//
// check_redis.go
// Quick and dirty redis SET/GET checker
//
// Created by Jeraimee Hughes
// Copyright (c) 2014 Jeraimee Hughes
// License: GNU GPL Version 2
// http://www.gnu.org/licenses/gpl-2.0.html
//

package main

import (
  "flag"
  "fmt"
  "os"
  "os/signal"
	"crypto/rand"
  "github.com/garyburd/redigo/redis"
)

var EndChannel chan string

var server string

func registerSignalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			EndChannel <- fmt.Sprintf("Received signal %d.", sig)
		}
	}()
}

func RandomString(length int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, length)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func main() {

  flag.StringVar(&server, "server", "127.0.0.1:6379", "Redis server to check including port")
  flag.Parse()

  EndChannel = make(chan string)

  registerSignalHandler()

  c, err := redis.Dial("tcp", server)
  if err != nil {
    fmt.Printf("CRITICAL: Failed to connect: %s", err.Error())
    return
  }
  defer c.Close()

  randString := RandomString(10)

  c.Do("SET", "check_redis", randString)

  s, err := redis.String(c.Do("GET", "check_redis"))
  if err != nil {
    fmt.Printf("CRITICAL: error: %s", err.Error())
    return
  }

  if s != randString {
    fmt.Printf("CRITICAL: check string was not correct! Got %s instead of %s", s, randString)
    return
  }

  fmt.Printf("OK")
}
