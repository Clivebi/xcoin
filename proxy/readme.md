#总则:  
1. 所有接口不区分post或者get，但是post只能使用标准FORM格式（ application/x-www-form-urlencoded）或者application/json  
2. 返回结果均是标准化的JSON  
3. 所有参数，若无特别说明，字符集默认为:所有字符数字加下划线，使用其它字符可能造成无法预知的逻辑错误  
4. 因为没有签名保证，此系统只能运行在内网环境  

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
#用户接口  

##添加用户  
http://127.0.0.1:8789/adduser?username=user001&bankname=bankA&usertype=1  
参数：username  用户名
	 bankname  开始时使用的bank，这个和用户能使用哪些货币无关，只是代表，这个用户注册时来自bankname所在的场子
	 usertype  1 普通用户，2 exchanger
备注：username是全局唯一的，注册时bankname指定的bank必须已经存在，否则注册失败
返回值：用户的信息

```
```
##查询用户  
http://127.0.0.1:8789/getuser?username=user001
参数：username 用户名

返回值：用户的信息
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


