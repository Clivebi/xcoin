总则:  
1、所有接口不区分post或者get，但是post只能使用标准FORM格式（ application/x-www-form-urlencoded）或者application/json. 
2、返回结果均是标准化的JSON  
3、所有参数，若无特别说明，字符集默认为:所有字符数字加下划线，使用其它字符可能造成无法预知的逻辑错误. 
4、因为没有签名保证，此系统只能运行在内网环境  

banker操作接口. 
创建. 
createbank?bankname=bankA&currency=USD&chip=chipA&exchanger=ex_A. 
参数：bankname 	银行名. 
 	 currency 	法币. 
 	 chip 	  	Token名字  
 	 exchanger 	和x管保持一致性，我们这里目前没有意义. 
备注：所有参数都要求是全局唯一的，也就是banname全局唯一，currency全局唯一，chip全局唯一，exchanger全局唯一. 
返回值：创建后，银行的具体信息，json格式. 


查询
getbank?bankname=bankA
参数：bankname 	银行名
返回值：银行的目前信息

调整（预支）资金池
userChangeThreshold?bankname=bankA&add_threshold=66
参数：bankname 银行名
	 add_threshold 增加值
备注：预支资金池影响issue的功能，issue将会使用资金池里面的资金，如果资金不够，issue将会失败，刚创建的银行，资金池为0，需要进行首次调整才能正常使用issue功能
返回值：操作完成后银行的信息

用户接口

添加用户
addmember?username=mem_test&bankname=bank_test&usertype=1&father=jqq

查询用户


其它功能
法币转Token（issue）

充值法币

提现法币

花费Token

法币/Token转账


