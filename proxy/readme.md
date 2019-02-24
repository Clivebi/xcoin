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
func : "addbank"  
args :["bankname","使用的法币名称","使用的Token名字","管理员的公钥或者ID"]   
权限： 只有root用户才能调用此接口  
备注：每个银行只能使用一种法币和一种Token币，而且名字是全局排重的，就是不能重复，管理员需要先adduser  
	 否则这个函数会调用失败，调用后，这个这个ID指定的用户被升级成银行管理员  
返回值：银行的信息，例子  
```
请求：
{
    "timestamp": 1550825585, 
    "fromid": "MIIBCgKCAQEA12ytFRV9NMC6LWh8HGPSQSoxKjwgBZ+3DrcVgKEj9t2a6AOVNMQNl9sH87Bkp208UFoneA0j6ry9s8l2da1a5fRj2opRPFF2S3cK7AhVzHHe7WNQESgYLpHxORxFC5Y5C5LGa0cj6megxu95GvEu61lbkxmScdWWfLiLLhI5Cr/5jOaFzNahsx1W3tdD66EkYMqKaCED67TUjVRU2b0EOTVz/Yw5xJSerSz35WDk4uo19rHt+Vn91DmbT+CVpooGQeTdiPbX2d39LMafqLrWOUgRhOAs3jlS8x9v0ERiLQWV9HMwABjQs+DmMwB7YW+Zxq62CrU5wzsPP6uNbDWXGwIDAQAB", 
    "func": "addbank", 
    "args": [
        "bankA", 
        "USD", 
        "TokenA", 
        "MIIBCgKCAQEAvMk0UUJK0r1wKbuqBlrEphKK7UR/EFOCz7p5oAyKpqX2hBZZQgaiQc9ebZRouGj+Giui2/S1eZDHAVzqeYNQ665l/TjvJerHTphp2NRVyBDESawuQxnyLDx81dY/iEPN9yg0YzXpssN8SvTOslo15O2SnkxJ6Wkno90pg34jjIM0oZWQ/K3u4W1alN9urOzYKzcC6ycJKbeDfBTHhEF/vm+HpzyHDsXNQ6Ax85CODu74+SE6JLqIBek0dSSn09VzWcS3C6tD/4+0IoqTb01MdoT72toUYUTV5p1zEJBbmr4/VJXkaJ4ecgraNe6URDl/TlMV5UtMExhdBqErhTXUlQIDAQAB"
    ]
}
```

```
回应：  
{
    "errmsg": "sucess", 
    "txid": "581614bd197abdc5d65a80ee76081a161810014cb4726c41d67cc24d5dcb48dc", 
    "valid_code": "VALID", 
    "data": {
        "bankname": "bankA", 	//	银行名字
        "chip": "TokenA", 		//	Token名字
        "chiplimit": 0, 		//	Token限额，用户可使用这种Token的数量
        "chipused": 0, 			//	已经使用的Token
        "currency": "USD", 		//	法币名字
        "currencyCount": 0, 	//	法币池值
        "exchangemap": { }, 	//	汇率表，只需要指定法币即可，例如 {"USD2HKD":7.0,"HKD2USD":0.14}
        "fiexedexchangemap": { }, 	//	固定汇率表
        "mangername": "5e9c47ef7d0565c43a4d17d53f0fc0da"	//	管理员ID
    }
}
```
##获取bank信息  
func : "getbank"  
args :["bankname"]   
权限： 任何人都可调用   
返回值：银行的信息，同addbank  

##设置bank的token限额  
func : "adjustbanklimit"  
args :["bankname","限额值"]   
权限： 只有root账户可调用此功能   
返回值：银行的信息，同addbank  

##设置汇率
func : "setexchanemap"  
args :["isfixedmap","newvalue"]   
权限： 只有银行管理员可以调用此接口，而且只能操作自己管理的银行   
返回值：银行的信息，同addbank  
```
请求：
{
    "timestamp": 1550825759, 
    "fromid": "MIIBCgKCAQEAzRjSSSk5Y4Rve2Yk4fViAnSa01iB6d35qqsWoNs/F4XoRTxp03j8jVGWnJFbG+oAoSnWjbf7Ba7K6BN5ClKDwRjh1T6DEiJAJmfzLArpZrMZbP7JnLCV4TYmOUPDzHSKz9//23NZO6wuhDTgwEqxPhSBe2zaNBI7PQkM6WNdc1ldN+Km6cwxg0P2mn+ltjKVjh4NY3LEBI45vs6vgO+8aZLVjfXqTM6DCMqZTGo6/IzU3N+AwtLR/m7KkDMFZkQGnj8J9+fn7WPqV+Dr9UxA7B7l1Pkm2BaosBKREjRuiU8oBTnoJLe/PTdSrgLblEwKqJGkwYRen7/JUsq9nz7s6wIDAQAB", 
    "func": "setexchanemap", 
    "args": [
        "false", 
        "{\"USD2HKD\":7.0,\"HKD2USD\":0.14}"
    ]
}
```

```
回应：  
{
    "errmsg": "sucess", 
    "txid": "aea750280173a0dcb9b6e5f17362585b809968f1c8818fd376cc89bce0ad4191", 
    "valid_code": "VALID", 
    "data": {
        "bankname": "bankB", 
        "chip": "TokenB", 
        "chiplimit": 1000000, 
        "chipused": 0, 
        "currency": "HKD", 
        "currencyCount": 0, 
        "exchangemap": {
            "HKD2USD": 0.14, 
            "USD2HKD": 7
        }, 
        "fiexedexchangemap": { }, 
        "mangername": "854c7458f81ecd02cb2d7d05503ea272"
    }
}
```
