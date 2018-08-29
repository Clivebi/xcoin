
checkout https://github.com/Clivebi/xcoin 到
$GOAPTH/src/github.com/hyperledger/

切换到这个目录：
$GOAPTH/src/github.com/hyperledger/xcoin/run-network
执行 byfn.sh -m up
完成后
执行测试用例
docker exec cli ./scripts/one.sh

清除
byfn.sh -m down



http 通信协议

1、生成rsa 公私钥对 http://xxxxxx/keygen.do

2、调用合约功能 http://xxxxxx/call.do,请求参数格式为json
请求参数 
{
	"func":"函数名",
	"args":"参数",
	"signature":"对参数的签名，使用rsa私钥对参数字符串签名所得",
}

其中参数为 json字符串，根据函数不同，格式也不同，目前支持的函数有
adduser		注册用户
{
	"pubkey":"公钥",
	"timestamp":int类型，时间戳
}
返回的信息中有个用户ID，需要自己保存，后续的其它接口需要使用

getuser		获取用户账户信息
{
	"callid":"客户ID",
	"id":"要查询的账户信息的目标ID",
	"timestamp":int类型，时间戳
}
upgradeuser	升级账户，只有root账户才能调用此接口，第一个注册的账户为root账户
调用此接口，用户被升级成seller
{
	"callid":"客户ID",
	"id":"要升级的客户ID",
	"timestamp":int类型，时间戳,
	"limit":int类型，此seller可以最多使用的币，超过这个值，将会不能转出
}
send 		转账
{
	"callid":"客户ID",
	"toid":"要给转账的ID",
	"timestamp":int类型，时间戳,
	"coin":int类型，要转出多少币
}