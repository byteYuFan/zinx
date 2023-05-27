package utils

import (
	"encoding/json"
	"github.com/byteYuFan/zinx/zinterfance"
	"os"
)

// 第三方的功能在此
// 这个包是关于解析Zinx的全局控制模块

// GlobalObj GlobalObj 存储一切有关zinx框架的全局参数，供其他模块使用
//
//	一些参数是可以通过用户通过zinx.json文件进行配置
type GlobalObj struct {
	/*
		Server端
	*/
	// TcpServer 当前zinx全局的server对象
	TcpServer zinterfance.IServer
	// Host 当前服务器主机监听的IP
	Host string
	// TcpPort 当前服务器监听的端口号
	TcpPort int
	// Name 当前服务器的名称
	Name string
	/*
	 Zinx
	*/
	// Version 当前zinx的版本号
	Version string
	// MaxConn 当前服务器主机允许的最大连接数
	MaxConn int
	// MaxPackageSize 当前zinx框架数据包的最大值
	MaxPackageSize uint32
}

// GlobalObject GlobalObject 定义一个全局的对象
var GlobalObject *GlobalObj

// Reload 加载用户自定义的参数
func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("../conf/zinx.json")
	if err != nil {
		panic(err)
	}
	// 将json文件数据解析到struct中
	err = json.Unmarshal(data, g)
	if err != nil {
		panic(err)
	}
}

// init 初始化全局的GlobalObject对象
func init() {
	// 如果配置文件没有加载，就是一个默认的值
	GlobalObject = &GlobalObj{
		TcpServer:      nil,
		Host:           "0.0.0.0",
		TcpPort:        8999,
		Name:           "ZinxServerApp",
		Version:        "V0.4",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}
	// 尝试从zinx.json读取
	GlobalObject.Reload()
}
