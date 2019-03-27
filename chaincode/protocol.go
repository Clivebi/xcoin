package main

//Request 数据打包格式
type Request struct {
	Time     int64    `json:"timestamp"` //时间戳
	Function string   `json:"func"`      //调用函数
	Args     []string `json:"args"`      //调用参数
}
