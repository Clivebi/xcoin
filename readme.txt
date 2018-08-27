
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