package main

type Request struct {
	Time     int64    `json:"timestamp"` //时间戳
	FromID   string   `json:"fromid"`    //调用者，用户的ID或者public key
	Function string   `json:"func"`      //调用函数
	Args     []string `json:"args"`      //调用参数
}

/*
调用合约功能
o.client.Execute(channel.Request{ChaincodeID: o.conf.ChainCode, Fcn: "callapi", Args: [requrst,signature]}
*/
