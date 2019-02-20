

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


项目说明
chaincode 
合约代码，主要逻辑实现在这里面

run-network
测试网络配置和脚本等

proxy
API编程接口，客户端软件通过http/https 和proxy通信，之后proxy调用合约功能，chaincode运行在hyperledger的docker之中，proxy运行在机器上（非docker）