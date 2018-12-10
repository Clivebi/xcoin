#!/bin/bash
export CORE_PEER_ADDRESS=peer0.org2.example.com:7051
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/xcoin/run-network/crypto-config/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/xcoin/run-network/crypto-config/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp


ORDERER_CA=/opt/gopath/src/github.com/hyperledger/xcoin/run-network/crypto-config/ordererOrganizations/example.com/msp/tlscacerts/tlsca.example.com-cert.pem
PEER0_ORG1_CA=/opt/gopath/src/github.com/hyperledger/xcoin/run-network/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
PEER0_ORG2_CA=/opt/gopath/src/github.com/hyperledger/xcoin/run-network/crypto-config/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt


peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile $ORDERER_CA -C mychannel -n mycc --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles $PEER0_ORG1_CA --peerAddresses peer0.org2.example.com:7051 --tlsRootCertFiles $PEER0_ORG2_CA -c '{"Args":["adduser","{\"pubkey\":\"MIIBCgKCAQEA1QVZzTjW/uq6tk/Oct6rO/HyPZB5+xSt2n7bJy7XPFRNbfeKV0gGYeHJy8ctbDCN8TxU2evxrPr5QaXcJsOHBoJdo9MLaxwirz5bT2Ctom7W2hIfUVcafzTvbRtpAZCkS+ZwGjn/u3/gsJqF0HUHZlmobGL9JxF0BF8vqD/x7VD0qaPlwPQNRk3cyuywIR/a1kAjxiXjWUrmnFvqpNZwO5P0+/KwQ3D8v/PS+s1ZaG1SaFqPSHm9CZp7wknin3/0prLPxnVmtUv3lsKP6lv4Vc6i3OSlBWBIlSIXG9GjDgCcAHvDZnixgoa6jv7Zif2qYfOTqwL/rTSgWPGRm3oWCwIDAQAB\",\"timestamp\":100}","signature"]}'
sleep 3
peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile $ORDERER_CA -C mychannel -n mycc --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles $PEER0_ORG1_CA --peerAddresses peer0.org2.example.com:7051 --tlsRootCertFiles $PEER0_ORG2_CA -c '{"Args":["adduser","{\"pubkey\":\"MIIBCgKCAQEAzYVkwedEM9pYZbGi6rSoDopr+R/rb3UDq40WXKFb1ZbVuwXXIWWvOlSGbNvSZ0CvAANNjYbY7DcPwILlXm/7AT0L43ZJIgTYHYOZYVVCXWwNZBWcOfnh63iwvMSQ47ml17q64fR2Cpy9JHK28k3IRM3of7NxqdbCoOUm2qqPtJ3y+Yo3tlSCywS+EhFmj9Ukf7yPA4rgmcCsybHn+i9wT19uA6oZmz0VgObgPEtaHzEPGFEtfh/LJVTRsxl3WMqaNFutDBDWTDtg4ieVmb8aX2gVNC26sJTaklzPWodfK+fZpVqYC9LgbLCURDVTmAUqGfGqqmdyGlDIEF+mfVrXwwIDAQAB\",\"timestamp\":100}","signature"]}'
sleep 3

peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile $ORDERER_CA -C mychannel -n mycc --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles $PEER0_ORG1_CA --peerAddresses peer0.org2.example.com:7051 --tlsRootCertFiles $PEER0_ORG2_CA -c '{"Args":["adduser","{\"pubkey\":\"MIIBCgKCAQEAsn+iNcrVLL4qEHGbXwW+zit08N0FzP54hP8x+EMhYbVUNE06ui7/8j2pm6WvjB3KXVdkL/CHSNad4UB9+asFeL5hCTZhmoBRaDZ13yqhtICHlukzemFNyrfE0LZJc4RRdQHa6eMcKpo/TJKCNHoBlESIX9QLmrNrr6GQ93obrx5FxSlP/iPPII+e1dNpTB7j7Lo/PdscsYnA0N1KRWHdqsrsYUM7sPwCRe8DoA+bPbZ0VA17HCoga6z68cpe3K+r+99uUbp1zYfpUCrg/mgDGJYNGDDlfqHZbDBkWBdeT7TZ9RUkhp20NBaY6oQ7+CDRYf/gRR1Gq8sjYAzEHog7GQIDAQAB\",\"timestamp\":100}","signature"]}'

sleep 3
peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile $ORDERER_CA -C mychannel -n mycc --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles $PEER0_ORG1_CA --peerAddresses peer0.org2.example.com:7051 --tlsRootCertFiles $PEER0_ORG2_CA -c '{"Args":["upgradeuser","{\"callid\":\"350051298d9ddd47bc0028e0f3d3fa2a\",\"id\":\"7fcf6391cea96aaed60069168f23b953\",\"limit\":1000000,\"timestamp\":100}","signature"]}'

sleep 3
peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile $ORDERER_CA -C mychannel -n mycc --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles $PEER0_ORG1_CA --peerAddresses peer0.org2.example.com:7051 --tlsRootCertFiles $PEER0_ORG2_CA -c '{"Args":["getuser","{\"callid\":\"7fcf6391cea96aaed60069168f23b953\",\"id\":\"7fcf6391cea96aaed60069168f23b953\",\"timestamp\":100}","signature"]}'
sleep 3
peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile $ORDERER_CA -C mychannel -n mycc --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles $PEER0_ORG1_CA --peerAddresses peer0.org2.example.com:7051 --tlsRootCertFiles $PEER0_ORG2_CA -c '{"Args":["send","{\"callid\":\"7fcf6391cea96aaed60069168f23b953\",\"toid\":\"2e99cf4e88e97162be67e663ebf79476\",\"coin\":100,\"timestamp\":100}","signature"]}'
sleep 3
peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile $ORDERER_CA -C mychannel -n mycc --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles $PEER0_ORG1_CA --peerAddresses peer0.org2.example.com:7051 --tlsRootCertFiles $PEER0_ORG2_CA -c '{"Args":["getuser","{\"callid\":\"2e99cf4e88e97162be67e663ebf79476\",\"id\":\"2e99cf4e88e97162be67e663ebf79476\",\"timestamp\":100}","signature"]}'

