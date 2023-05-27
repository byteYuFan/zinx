package main

import "github.com/byteYuFan/zinx/znet"

func main() {
	s := znet.NewServer("test")
	s.Server()
}
