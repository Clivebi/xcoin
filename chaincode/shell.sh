
#安装链码 peer0.org1.example.com
export CORE_PEER_ADDRESS=peer0.org1.example.com:7051
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
peer chaincode install -n xcoin -v 1.0 -l golang -p github.com/hyperledger/xcoin/chaincode


#安装链码 peer1.org1.example.com
export CORE_PEER_ADDRESS=peer1.org1.example.com:7051
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
peer chaincode install -n xcoin -v 1.0 -l golang -p github.com/hyperledger/xcoin/chaincode


#安装链码 peer0.org2.example.com
export CORE_PEER_ADDRESS=peer0.org2.example.com:7051
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
peer chaincode install -n xcoin -v 1.0 -l golang -p github.com/hyperledger/xcoin/chaincode


#安装链码 peer1.org2.example.com
export CORE_PEER_ADDRESS=peer1.org2.example.com:7051
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
peer chaincode install -n xcoin -v 1.0 -l golang -p github.com/hyperledger/xcoin/chaincode


#实例化链码 peer1.org2.example.com
export CORE_PEER_ADDRESS=peer1.org2.example.com:7051
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n xcoin -l golang  -v 1.0 -c '{"Args":["init","debug"]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"


#执行函数
export CORE_PEER_ADDRESS=peer0.org2.example.com:7051
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n mycc --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses peer0.org2.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["adduser","{\"pubkey\":\"MIIBCgKCAQEA1QVZzTjW/uq6tk/Oct6rO/HyPZB5+xSt2n7bJy7XPFRNbfeKV0gGYeHJy8ctbDCN8TxU2evxrPr5QaXcJsOHBoJdo9MLaxwirz5bT2Ctom7W2hIfUVcafzTvbRtpAZCkS+ZwGjn/u3/gsJqF0HUHZlmobGL9JxF0BF8vqD/x7VD0qaPlwPQNRk3cyuywIR/a1kAjxiXjWUrmnFvqpNZwO5P0+/KwQ3D8v/PS+s1ZaG1SaFqPSHm9CZp7wknin3/0prLPxnVmtUv3lsKP6lv4Vc6i3OSlBWBIlSIXG9GjDgCcAHvDZnixgoa6jv7Zif2qYfOTqwL/rTSgWPGRm3oWCwIDAQAB\",\"timestamp\":100}","signature"]}'

#执行函数
export CORE_PEER_ADDRESS=peer1.org2.example.com:7051
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
peer chaincode query -C mychannel -n xcoin -c '{"Args":["query","a"]}'












peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n mycc --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses peer0.org2.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["adduser","{\"pubkey\":\"MIIBCgKCAQEA1QVZzTjW/uq6tk/Oct6rO/HyPZB5+xSt2n7bJy7XPFRNbfeKV0gGYeHJy8ctbDCN8TxU2evxrPr5QaXcJsOHBoJdo9MLaxwirz5bT2Ctom7W2hIfUVcafzTvbRtpAZCkS+ZwGjn/u3/gsJqF0HUHZlmobGL9JxF0BF8vqD/x7VD0qaPlwPQNRk3cyuywIR/a1kAjxiXjWUrmnFvqpNZwO5P0+/KwQ3D8v/PS+s1ZaG1SaFqPSHm9CZp7wknin3/0prLPxnVmtUv3lsKP6lv4Vc6i3OSlBWBIlSIXG9GjDgCcAHvDZnixgoa6jv7Zif2qYfOTqwL/rTSgWPGRm3oWCwIDAQAB\",\"timestamp\":100}","signature"]}'

peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n mycc --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses peer0.org2.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["adduser","{\"pubkey\":\"MIIBCgKCAQEAzYVkwedEM9pYZbGi6rSoDopr+R/rb3UDq40WXKFb1ZbVuwXXIWWvOlSGbNvSZ0CvAANNjYbY7DcPwILlXm/7AT0L43ZJIgTYHYOZYVVCXWwNZBWcOfnh63iwvMSQ47ml17q64fR2Cpy9JHK28k3IRM3of7NxqdbCoOUm2qqPtJ3y+Yo3tlSCywS+EhFmj9Ukf7yPA4rgmcCsybHn+i9wT19uA6oZmz0VgObgPEtaHzEPGFEtfh/LJVTRsxl3WMqaNFutDBDWTDtg4ieVmb8aX2gVNC26sJTaklzPWodfK+fZpVqYC9LgbLCURDVTmAUqGfGqqmdyGlDIEF+mfVrXwwIDAQAB\",\"timestamp\":100}","signature"]}'


peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n mycc --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses peer0.org2.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["adduser","{\"pubkey\":\"MIIBCgKCAQEAsn+iNcrVLL4qEHGbXwW+zit08N0FzP54hP8x+EMhYbVUNE06ui7/8j2pm6WvjB3KXVdkL/CHSNad4UB9+asFeL5hCTZhmoBRaDZ13yqhtICHlukzemFNyrfE0LZJc4RRdQHa6eMcKpo/TJKCNHoBlESIX9QLmrNrr6GQ93obrx5FxSlP/iPPII+e1dNpTB7j7Lo/PdscsYnA0N1KRWHdqsrsYUM7sPwCRe8DoA+bPbZ0VA17HCoga6z68cpe3K+r+99uUbp1zYfpUCrg/mgDGJYNGDDlfqHZbDBkWBdeT7TZ9RUkhp20NBaY6oQ7+CDRYf/gRR1Gq8sjYAzEHog7GQIDAQAB\",\"timestamp\":100}","signature"]}'


peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n mycc --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses peer0.org2.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["upgradeuser","{\"callid\":\"350051298d9ddd47bc0028e0f3d3fa2a\",\"id\":\"7fcf6391cea96aaed60069168f23b953\",\"limit\":1000000,\"timestamp\":100}","signature"]}'


peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n mycc --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses peer0.org2.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["getuser","{\"callid\":\"7fcf6391cea96aaed60069168f23b953\",\"id\":\"7fcf6391cea96aaed60069168f23b953\",\"timestamp\":100}","signature"]}'

peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n mycc --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses peer0.org2.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["send","{\"callid\":\"7fcf6391cea96aaed60069168f23b953\",\"toid\":\"2e99cf4e88e97162be67e663ebf79476\",\"coin\":100,\"timestamp\":100}","signature"]}'

peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n mycc --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses peer0.org2.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["getuser","{\"callid\":\"2e99cf4e88e97162be67e663ebf79476\",\"id\":\"2e99cf4e88e97162be67e663ebf79476\",\"timestamp\":100}","signature"]}'

