#总则:  
1. 所有对外暴露的接口只有一个就是/callapi.do,参与需使用post json（application/json）的格式调用
   调用参数为：
   {
   	 "req":"字符串类型的请求数据包，由Request序列化成字符串而成",
   	 "sig":"调用者使用私钥对req签名后的结果"
   }
   其中req为真正的请求结构体Request如下，序列化而成的字符串
   ```
   {
    "timestamp": 1550825521,  	 //时间戳，用于对抗 重放攻击
    "fromid": "user public key", //调用者的ID或者公钥,如果是公钥，使用base64 编码的x509 序列化公钥
    "func": "adduser", 			 //要调用的功能函数
    "args": ["arg1","arg2"] 	 //传递给功能函数的参数
   }
   ```
   由上述说明可知，对于调用功能而言，每次需要变更的是func 和args参数，下面针对的是func和args参数的说明   
2. 返回结果均是标准化的JSON，其格式如下   
```
{
    "errmsg": "sucess", 		//错误信息，成功为sucess，否则为具体的错误信息
    "txid": "e490e9149a12239e3a1f5c29aee2d912cbeb68bcae4d5f5ec45920422b226fd3", //transcation ID
    "valid_code": "VALID",  //transcation 有效性
    "data": {} 				//这个值为调用func之后返回的信息，json格式，根据func不同而不同
}
```

#用户接口  

##添加用户  
func : "adduser"
args :["新增用户的RSA公钥"]
权限： 任何人都可以调用，第一个添加的用户为root账户
返回值：用户的信息，例子
```
请求：
{
    "timestamp": 1550825521, 
    "fromid": "MIIBCgKCAQEA12ytFRV9NMC6LWh8HGPSQSoxKjwgBZ+3DrcVgKEj9t2a6AOVNMQNl9sH87Bkp208UFoneA0j6ry9s8l2da1a5fRj2opRPFF2S3cK7AhVzHHe7WNQESgYLpHxORxFC5Y5C5LGa0cj6megxu95GvEu61lbkxmScdWWfLiLLhI5Cr/5jOaFzNahsx1W3tdD66EkYMqKaCED67TUjVRU2b0EOTVz/Yw5xJSerSz35WDk4uo19rHt+Vn91DmbT+CVpooGQeTdiPbX2d39LMafqLrWOUgRhOAs3jlS8x9v0ERiLQWV9HMwABjQs+DmMwB7YW+Zxq62CrU5wzsPP6uNbDWXGwIDAQAB", 
    "func": "adduser", 
    "args": [
        "MIIBCgKCAQEA12ytFRV9NMC6LWh8HGPSQSoxKjwgBZ+3DrcVgKEj9t2a6AOVNMQNl9sH87Bkp208UFoneA0j6ry9s8l2da1a5fRj2opRPFF2S3cK7AhVzHHe7WNQESgYLpHxORxFC5Y5C5LGa0cj6megxu95GvEu61lbkxmScdWWfLiLLhI5Cr/5jOaFzNahsx1W3tdD66EkYMqKaCED67TUjVRU2b0EOTVz/Yw5xJSerSz35WDk4uo19rHt+Vn91DmbT+CVpooGQeTdiPbX2d39LMafqLrWOUgRhOAs3jlS8x9v0ERiLQWV9HMwABjQs+DmMwB7YW+Zxq62CrU5wzsPP6uNbDWXGwIDAQAB"
    ]
}
```

```
回应：
{
    "errmsg": "sucess", 
    "txid": "e490e9149a12239e3a1f5c29aee2d912cbeb68bcae4d5f5ec45920422b226fd3", 
    "valid_code": "VALID", 
    "data": {
        "balance": { }, // 用户的资产信息，例如 {"USD":100,"TokenA":800}
        "id": "fbb3f222c61d9092927dd066460290af", //用户ID
        "lockedbalance": { }, // 用户的锁定资产信息
        "pub_key": "MIIBCgKCAQEA12ytFRV9NMC6LWh8HGPSQSoxKjwgBZ+3DrcVgKEj9t2a6AOVNMQNl9sH87Bkp208UFoneA0j6ry9s8l2da1a5fRj2opRPFF2S3cK7AhVzHHe7WNQESgYLpHxORxFC5Y5C5LGa0cj6megxu95GvEu61lbkxmScdWWfLiLLhI5Cr/5jOaFzNahsx1W3tdD66EkYMqKaCED67TUjVRU2b0EOTVz/Yw5xJSerSz35WDk4uo19rHt+Vn91DmbT+CVpooGQeTdiPbX2d39LMafqLrWOUgRhOAs3jlS8x9v0ERiLQWV9HMwABjQs+DmMwB7YW+Zxq62CrU5wzsPP6uNbDWXGwIDAQAB", 			 // 用户公钥
        "type": 0		// 0 root账户，1 bank manger  2 normal user
    }
}
```

##查询用户信息  
func : "getuser"
args :["用户的public key或者用户ID"]
备注：用户ID在adduser之后返回
权限： 任何人都可以调用，既可以查询自己，也可以查询别人
返回值：用户的信息，例子
```
请求：
{
    "timestamp": 1550825521, 
    "fromid": "MIIBCgKCAQEA12ytFRV9NMC6LWh8HGPSQSoxKjwgBZ+3DrcVgKEj9t2a6AOVNMQNl9sH87Bkp208UFoneA0j6ry9s8l2da1a5fRj2opRPFF2S3cK7AhVzHHe7WNQESgYLpHxORxFC5Y5C5LGa0cj6megxu95GvEu61lbkxmScdWWfLiLLhI5Cr/5jOaFzNahsx1W3tdD66EkYMqKaCED67TUjVRU2b0EOTVz/Yw5xJSerSz35WDk4uo19rHt+Vn91DmbT+CVpooGQeTdiPbX2d39LMafqLrWOUgRhOAs3jlS8x9v0ERiLQWV9HMwABjQs+DmMwB7YW+Zxq62CrU5wzsPP6uNbDWXGwIDAQAB", 
    "func": "getuser", 
    "args": [
        "MIIBCgKCAQEA12ytFRV9NMC6LWh8HGPSQSoxKjwgBZ+3DrcVgKEj9t2a6AOVNMQNl9sH87Bkp208UFoneA0j6ry9s8l2da1a5fRj2opRPFF2S3cK7AhVzHHe7WNQESgYLpHxORxFC5Y5C5LGa0cj6megxu95GvEu61lbkxmScdWWfLiLLhI5Cr/5jOaFzNahsx1W3tdD66EkYMqKaCED67TUjVRU2b0EOTVz/Yw5xJSerSz35WDk4uo19rHt+Vn91DmbT+CVpooGQeTdiPbX2d39LMafqLrWOUgRhOAs3jlS8x9v0ERiLQWV9HMwABjQs+DmMwB7YW+Zxq62CrU5wzsPP6uNbDWXGwIDAQAB"
    ]
}
```

```
回应：
{
    "errmsg": "sucess", 
    "txid": "e490e9149a12239e3a1f5c29aee2d912cbeb68bcae4d5f5ec45920422b226fd3", 
    "valid_code": "VALID", 
    "data": {
        "balance": { }, 
        "id": "fbb3f222c61d9092927dd066460290af", 
        "lockedbalance": { }, 
        "pub_key": "MIIBCgKCAQEA12ytFRV9NMC6LWh8HGPSQSoxKjwgBZ+3DrcVgKEj9t2a6AOVNMQNl9sH87Bkp208UFoneA0j6ry9s8l2da1a5fRj2opRPFF2S3cK7AhVzHHe7WNQESgYLpHxORxFC5Y5C5LGa0cj6megxu95GvEu61lbkxmScdWWfLiLLhI5Cr/5jOaFzNahsx1W3tdD66EkYMqKaCED67TUjVRU2b0EOTVz/Yw5xJSerSz35WDk4uo19rHt+Vn91DmbT+CVpooGQeTdiPbX2d39LMafqLrWOUgRhOAs3jlS8x9v0ERiLQWV9HMwABjQs+DmMwB7YW+Zxq62CrU5wzsPP6uNbDWXGwIDAQAB", 
        "type": 0
    }
}
```

#banker操作接口  
##创建bank  

http://127.0.0.1:8789/createbank?bankname=bankA&currency=USD&chip=chipA&exchanger=ex_A  
参数：bankname 	银行名  
 	 currency 	法币  
 	 chip 	  	Token名字  
 	 exchanger 	和x管保持一致性，我们这里目前没有意义  
备注：所有参数都要求是全局唯一的，也就是banname全局唯一，currency全局唯一，chip全局唯一，exchanger全局唯一  
返回值：创建后，银行的具体信息，json格式  
```
{
	"status": 200,
	"errmsg": "sucess",
	"txid": "4ab10d98cc9f2dfb85be992f05b67f2aff66a562acdb9a99e355c3cfbfcd7a84",
	"valid_code": "VALID",
	"data": {
		"bankname": "bankA",
		"chip": "chipA",
		"currency": "USD",
		"exchanger": "ex_A",
		"totalamount": 0,
		"usedamount": 0
	}
}
```

##获取bank信息  
http://127.0.0.1:8789/getbank?bankname=bankA  
参数：bankname 	银行名  
返回值：银行的目前信息  
```
```

##增加（预支）资金池  
http://127.0.0.1:8789/changebanklimit?bankname=bankA&add_threshold=50000  
参数：bankname 银行名  
	 add_threshold 增加值  
备注：预支资金池影响issue的功能，issue将会使用资金池里面的资金，如果资金不够，issue将会失败，刚创建的银行，资金池为0，需要进行首次调整才能正常使用issue功能  
返回值：操作完成后银行的信息  
```
```

#其它功能  
##充值法币  
http://127.0.0.1:8789/cashin?username=user001&currency=USD&amount=100
参数：username 用户名
	 currency 法币名
	 amount   数量
备注：每个用户可拥有多种法币，所以充值的时候，可以接受任意法币币种
返回值：充值后用户信息
```


```
##法币转Token（issue）  
http://127.0.0.1:8789/issue?username=user001&bankname=bankA&currency=USD&amount=100
充值法币  

##提现法币  
http://127.0.0.1:8789/cashout?username=user001&bankname=banA&fromcurrency=USD&dstcurrency=USD&amount=100
参数：username 		用户名
	 bankname 		从哪个银行支出
	 fromcurrency 	支出的货币类型，因为每个银行只有一种法币，如果这个法币不是支出银行所有的，将会返回错误
	 dstcurrency	提取成哪种法币，目前仅支持fromcurrency == dstcurrency的情况
	 amount   		数量
返回值：提现后返回用户信息
```

```

##花费Token  

##法币/Token转账  


