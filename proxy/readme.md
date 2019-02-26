#总则:  
1. 所有对外暴露的接口只有一个 POST/GET <127.0.0.1:8789/callapi.do>调用参数为：  
   ```
   Post 的JSON格式如下：
   {  
   	 "req":"请求数据包由Request序列化成字符串",  
   	 "sig":"签名信息Signature序列化字符串"  
   }

   Request 如下：
   {
    "timestamp": 1550825521,     //时间戳，用于对抗 重放攻击
    "fromid": "user public key", //调用者的ID或者公钥,如果是公钥，使用base58
    "func": "adduser",           //要调用的功能函数
    "args": ["arg1","arg2"]      //传递给功能函数的参数，字符串数组
   }
   调用不同的功能，需要不同的func和args参数

   Signature 如下：
   {
     "caller":"函数调用者使用私钥对req进行RSA签名的结果的base58",
     "optuser":"目标用户使用私钥对req进行RSA签名的结果的base58"
   }
   其中，所有的功能调用，caller的签名是必须的，optuser的签名只有在需要在场证明的情况下才需要，目前只有cashin，cashout需要在场证明  
   即需要bank manger和用户同时在场，功能由bank manger调用，但是同时需要用户的签名  
   ```  
2. 返回结果均是标准化的JSON，其格式如下   
```
{
    "errmsg": "sucess", 		//错误信息，成功为sucess，否则为具体的错误信息
    "txid": "e490e9149a12239e3a1f5c29aee2d912cbeb68bcae4d5f5ec45920422b226fd3", //transcation ID
    "valid_code": "VALID",  //transcation 有效性
    "data": {} 				//这个值为调用func之后返回的信息，json格式，根据func不同而不同
}
```
3. 所有接口的测试代码可以参考:https://github.com/Clivebi/xcoin/blob/master/proxy/proxy_test.go  
4. 用户ID算法，base58(sha256(publicKey))


#用户接口  

##添加用户  
func : "adduser"  
args :["新增用户的RSA公钥"]  
权限： 任何人都可以调用，第一个添加的用户为root账户  
返回值：用户的信息，例子  
```
请求：
{
    "timestamp": 1551153414,
    "fromid": "e14Lihsy8n7MBnQvUMooQWy9exJ7j3evooxmhXjaJwtizR",
    "func": "adduser",
    "args": [
        "4D1btsFgbEQuvgZS6pu7rmfweC3UH1QprFALPMDeXVo8dsw5brJWFegfQff5owgQvqsQmUJUW8JHvaKGxBVppugbzkihPDz49YrgEQfzASHvNUfY56GEMDHvVF3fhYZG8fKTJS32JZRpFFo865HFgTUF8fVkxjhwanGiQnf8NmDC5WoDeBeempAseknTbu6ScjbzaxKMU1vyX78KHR73ocoAyY8B7eUFurBqkYkYWjRNh9h9DmksHFus1ZBUUZKRXeyRnH7iTURob4RCQCeYCDwEZ9mud4cJZZbSiZVLT9LM3EXpVQpKDd81ricDfKtPshPzN7feWa7d1J5An2npRZsJsMxzHTCjWuPk8kSzM3mkWCDwz"
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
        "id": "e14Lihsy8n7MBnQvUMooQWy9exJ7j3evooxmhXjaJwtizR", //用户ID
        "lockedbalance": { }, // 用户的锁定资产信息
        "pub_key": "4D1btsFgbEQuvgZS6pu7rmfweC3UH1QprFALPMDeXVo8dsw5brJWFegfQff5owgQvqsQmUJUW8JHvaKGxBVppugbzkihPDz49YrgEQfzASHvNUfY56GEMDHvVF3fhYZG8fKTJS32JZRpFFo865HFgTUF8fVkxjhwanGiQnf8NmDC5WoDeBeempAseknTbu6ScjbzaxKMU1vyX78KHR73ocoAyY8B7eUFurBqkYkYWjRNh9h9DmksHFus1ZBUUZKRXeyRnH7iTURob4RCQCeYCDwEZ9mud4cJZZbSiZVLT9LM3EXpVQpKDd81ricDfKtPshPzN7feWa7d1J5An2npRZsJsMxzHTCjWuPk8kSzM3mkWCDwz", 			 // 用户公钥
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
    "timestamp": 1551153460,
    "fromid": "e122aNbAASxBPMRcpEL2Eut98sDFfjhH9cGxKYxEdwAR6s",
    "func": "getuser",
    "args": [
        "4D1btsFgbEQuvRZTbkqUCdC1btMGSN8a2K7Z8Go2t2YgLUFhbZuKwva59PyKJ7Ndh7NHQPpDCjQXGb49cAzjvhDhVVWR5J1pyvpfWc2rhkJNyGZ2xfXaFDNWkgc2Dnhyb3KteHZ2wXUZRyfAgYKNFg8XDvqKF8ns9cB3UnGbYn7pTeZDphtwLrFezExadbDC1PtSEizYNYye52AoskAYrCoiUnzk3TKocwYUdmaufTmH9tYepheyV8MbPFKSDpJqXFV6ZwiQ3mgRrvmphtXz4Xbd4YiG67g878dvsGVMTTcQ8DWYX1zYXZ8xFyJZy2G2tHB7vbELJGvrFfDfxuQtXSmU5dW7TKPnbcg2mb4dk5TiTN9LH"
    ]
}
```

```
回应：  
{
    "errmsg": "sucess",
    "txid": "22de3f90b950fce614026f05898772594498922834236273b8deed4437b9c81d",
    "valid_code": "VALID",
    "data": {
        "balance": {},
        "id": "e1dW9MvoiNX4SajEU1VVVMkYsL2ehijSmhRbVWUKXPXXYk",
        "lockedbalance": {},
        "pub_key": "4D1btsFgbEQuvRZTbkqUCdC1btMGSN8a2K7Z8Go2t2YgLUFhbZuKwva59PyKJ7Ndh7NHQPpDCjQXGb49cAzjvhDhVVWR5J1pyvpfWc2rhkJNyGZ2xfXaFDNWkgc2Dnhyb3KteHZ2wXUZRyfAgYKNFg8XDvqKF8ns9cB3UnGbYn7pTeZDphtwLrFezExadbDC1PtSEizYNYye52AoskAYrCoiUnzk3TKocwYUdmaufTmH9tYepheyV8MbPFKSDpJqXFV6ZwiQ3mgRrvmphtXz4Xbd4YiG67g878dvsGVMTTcQ8DWYX1zYXZ8xFyJZy2G2tHB7vbELJGvrFfDfxuQtXSmU5dW7TKPnbcg2mb4dk5TiTN9LH",
        "type": 2
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
    "timestamp": 1551153467,
    "fromid": "e14Lihsy8n7MBnQvUMooQWy9exJ7j3evooxmhXjaJwtizR",
    "func": "addbank",
    "args": [
        "bankA",
        "USD",
        "TokenA",
        "e1dzoVvkTN1frBnmWnKy7wmFRim8UWaf34Ej1yQ4AP4D1c"
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
        "mangername": "e1dzoVvkTN1frBnmWnKy7wmFRim8UWaf34Ej1yQ4AP4D1c"	//	管理员ID
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
    "timestamp": 1551153651,
    "fromid": "e13dZV5pMXpY66Pr2HH6q7BfVZ5BzZHvRowjAxvCxJXM9Y",
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
    "txid": "e5c29d9cb46858bc553b8281ec66763f4b2c6402319d42a406653cdf2cf5d430",
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
        "fiexedexchangemap": {},
        "mangername": "e13dZV5pMXpY66Pr2HH6q7BfVZ5BzZHvRowjAxvCxJXM9Y"
    }
}
```
##法币充值  
func : "cashin"  
args :["目标用户","目标货币类型","数量"]   
权限： 只有银行管理员能够调用，而且只能操作自己管理的银行，即第二个参数必须和自己管理银行的法币相同   
返回值：充值后目标用户的账户信息，同addbank  
```
请求：
{
    "timestamp": 1551153525,
    "fromid": "e1dzoVvkTN1frBnmWnKy7wmFRim8UWaf34Ej1yQ4AP4D1c",
    "func": "cashin",
    "args": [
        "e122aNbAASxBPMRcpEL2Eut98sDFfjhH9cGxKYxEdwAR6s",
        "USD",
        "2000"
    ]
}
```

```
回应：  
{
    "errmsg": "sucess",
    "txid": "19ea0f7e1e19db5b6ecb60d1c6dd76c72d90280f8f391db5884d01cdaa5ef178",
    "valid_code": "VALID",
    "data": {
        "balance": {
            "USD": 2000
        },
        "id": "e122aNbAASxBPMRcpEL2Eut98sDFfjhH9cGxKYxEdwAR6s",
        "lockedbalance": {},
        "pub_key": "4D1btsFgbEQuw3aCeAgrEiqRBoNBabW8jd9L14pcTbGrE5Ub41A8fTUU6k3taUt5LeAqX1D4dEmsYou6PvVxJnhLcDBJMMvvtGdhj61DR9hRZZVzMBMBMewvPANkdAdGkRBHwWcUh6z6iU9AWAmVUfmZbzDvaqXpBkh4cMddg2sfhxAHhr419j7iU5zcmErEGjfB8Y1peFLUibMxacES9V11qP7aYah4iFmEk6r8TDxfnG4JX5dTFmdvaBmoADDNXsoE6A7ANxsUBCLCtw56ZRDABwVAJNy6uFC2TGrftujNEXnVMnt6Dg4AGVQ6AyrEeZpd2C5LHg7vW91roExXQrZsRU3VQvjFMVrAR2wrbebVEC5CH",
        "type": 2
    }
}
```
##法币提取  
func : "cashout"  
args :["目标用户","目标货币类型","数量"]      
权限： 只有银行管理员能够调用，而且只能操作自己管理的银行，即第二个参数必须和自己管理银行的法币相同   
返回值：提取后目标的账户信息，同addbank  
```
请求：
{
    "timestamp": 1551153546,
    "fromid": "e1dzoVvkTN1frBnmWnKy7wmFRim8UWaf34Ej1yQ4AP4D1c",
    "func": "cashout",
    "args": [
        "e122aNbAASxBPMRcpEL2Eut98sDFfjhH9cGxKYxEdwAR6s",
        "USD",
        "100"
    ]
}
```

```
回应：  
{
    "errmsg": "sucess",
    "txid": "e277751a4445321eeaeae3121a31c6cafee363d793f26e02a77a3b409b48547b",
    "valid_code": "VALID",
    "data": {
        "balance": {
            "USD": 1900
        },
        "id": "e122aNbAASxBPMRcpEL2Eut98sDFfjhH9cGxKYxEdwAR6s",
        "lockedbalance": {},
        "pub_key": "4D1btsFgbEQuw3aCeAgrEiqRBoNBabW8jd9L14pcTbGrE5Ub41A8fTUU6k3taUt5LeAqX1D4dEmsYou6PvVxJnhLcDBJMMvvtGdhj61DR9hRZZVzMBMBMewvPANkdAdGkRBHwWcUh6z6iU9AWAmVUfmZbzDvaqXpBkh4cMddg2sfhxAHhr419j7iU5zcmErEGjfB8Y1peFLUibMxacES9V11qP7aYah4iFmEk6r8TDxfnG4JX5dTFmdvaBmoADDNXsoE6A7ANxsUBCLCtw56ZRDABwVAJNy6uFC2TGrftujNEXnVMnt6Dg4AGVQ6AyrEeZpd2C5LHg7vW91roExXQrZsRU3VQvjFMVrAR2wrbebVEC5CH",
        "type": 2
    }
}
```
##转账  
func : "transfer"  
args :["目标用户的公钥或者ID","货币类型","数量","是否是从locked的资金转出"]   
权限： 任何人都可以调用此接口   
备注：普通用户可以转给普通用户，银行管理员可以转给普通用户，普通用户可以转个银行管理员，除此之外的转账将引发一个错误  
返回值：转账后自己的账户信息，同addbank  
```
请求：
{
    "timestamp": 1551153567,
    "fromid": "e122aNbAASxBPMRcpEL2Eut98sDFfjhH9cGxKYxEdwAR6s",
    "func": "transfer",
    "args": [
        "e1dW9MvoiNX4SajEU1VVVMkYsL2ehijSmhRbVWUKXPXXYk",
        "USD",
        "100",
        "false"
    ]
}
```

```
回应：  
{
    "errmsg": "sucess",
    "txid": "194aab5e7da5329fc95297d466a3e531b2ad7f9aa22c48ccddbf4b1a2899303c",
    "valid_code": "VALID",
    "data": {
        "balance": {
            "USD": 1800
        },
        "id": "e122aNbAASxBPMRcpEL2Eut98sDFfjhH9cGxKYxEdwAR6s",
        "lockedbalance": {},
        "pub_key": "4D1btsFgbEQuw3aCeAgrEiqRBoNBabW8jd9L14pcTbGrE5Ub41A8fTUU6k3taUt5LeAqX1D4dEmsYou6PvVxJnhLcDBJMMvvtGdhj61DR9hRZZVzMBMBMewvPANkdAdGkRBHwWcUh6z6iU9AWAmVUfmZbzDvaqXpBkh4cMddg2sfhxAHhr419j7iU5zcmErEGjfB8Y1peFLUibMxacES9V11qP7aYah4iFmEk6r8TDxfnG4JX5dTFmdvaBmoADDNXsoE6A7ANxsUBCLCtw56ZRDABwVAJNy6uFC2TGrftujNEXnVMnt6Dg4AGVQ6AyrEeZpd2C5LHg7vW91roExXQrZsRU3VQvjFMVrAR2wrbebVEC5CH",
        "type": 2
    }
}

```
##货币兑换  
func : "exchange"  
args :["付出的货币类型","目标货币类型","数量","是否使用固定汇率"]   
权限： 任何人都可以调用此接口   
返回值：兑换后自己的账户信息，同addbank  
```
请求：
{
    "timestamp": 1551153693,
    "fromid": "e122aNbAASxBPMRcpEL2Eut98sDFfjhH9cGxKYxEdwAR6s",
    "func": "exchange",
    "args": [
        "USD",
        "TokenB",
        "100",
        "false"
    ]
}
```

```
回应：  

{
    "errmsg": "sucess",
    "txid": "9f58042349cf54d7478dbf43d1ab11282083f7510be4694d5f2bedb4f5b39a25",
    "valid_code": "VALID",
    "data": {
        "balance": {
            "HKD": 700,
            "TokenA": 400,
            "TokenB": 700,
            "USD": 1100
        },
        "id": "e122aNbAASxBPMRcpEL2Eut98sDFfjhH9cGxKYxEdwAR6s",
        "lockedbalance": {},
        "pub_key": "4D1btsFgbEQuw3aCeAgrEiqRBoNBabW8jd9L14pcTbGrE5Ub41A8fTUU6k3taUt5LeAqX1D4dEmsYou6PvVxJnhLcDBJMMvvtGdhj61DR9hRZZVzMBMBMewvPANkdAdGkRBHwWcUh6z6iU9AWAmVUfmZbzDvaqXpBkh4cMddg2sfhxAHhr419j7iU5zcmErEGjfB8Y1peFLUibMxacES9V11qP7aYah4iFmEk6r8TDxfnG4JX5dTFmdvaBmoADDNXsoE6A7ANxsUBCLCtw56ZRDABwVAJNy6uFC2TGrftujNEXnVMnt6Dg4AGVQ6AyrEeZpd2C5LHg7vW91roExXQrZsRU3VQvjFMVrAR2wrbebVEC5CH",
        "type": 2
    }
}
```

###golang编程接口  
所有对外的API位于client.go  
```
把RSA公钥编码成字符串
func PublicKeyToString(publicKey *rsa.PublicKey) string

从公钥获取用户ID（钱包ID）
func PublicKeyToID(publicKey *rsa.PublicKey) string

编码发起函数调用的req字符串
func NewRequest(callID string, function string, args []string) (string, error)
callID      -- 调用者的ID或者使用PublicKeyToString编码的公钥
function    -- 要调用的功能
args        -- 调用function时使用的参数，如果其中一个参数是JSON格式，则需要序列化成字符串

获取请求的签名
func SignRequest(request string, privatekey *rsa.PrivateKey) (string, error)
request     -- NewRequest返回的字符串

发起调用
func CallAPI(apiURI string, request string, callersig string, optusersig string) (string, error) 
apiURI  -- proxy所在的HTTP URL
request -- NewRequest返回的字符串
callersig -- 调用者的签名
optusersig -- 参与者的签名，除了cashin，cashout 功能这个不能为空外，其它功能，这个参数为空字符串

通常调用流程：
// example
//   req,_ := NewRequest(PublicKeyToID(publicKey),"adduser",[]string{PublicKeyToString(publicKey)},)
//   callersig,_ :=SignRequest(req,privateKey)
//   rsp,err := CallAPI("http://127.0.0.1:8789/callapi.do",req,callersig,"")
//   parse rsp ...
更多例子，参考：https://github.com/Clivebi/xcoin/blob/master/proxy/proxy_test.go
```