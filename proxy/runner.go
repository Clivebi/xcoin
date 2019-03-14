package proxy

import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"io/ioutil"
)

// Response  返回结构体
type Response struct {
	ErrorMessage string                 `json:"errmsg"`
	TxID         string                 `json:"txid"`
	TxValidCode  string                 `json:"valid_code"`
	Payload      map[string]interface{} `json:"data"`
}

func getResponse(payload string, txID string, Code string, err error) []byte {
	var data map[string]interface{}
	json.Unmarshal([]byte(payload), &data)
	if data == nil {
		data = map[string]interface{}{}
	}
	rsp := &Response{
		ErrorMessage: "sucess",
		TxID:         txID,
		TxValidCode:  Code,
		Payload:      data,
	}
	if err != nil {
		rsp.ErrorMessage = err.Error()
	}
	buf, _ := json.Marshal(rsp)
	return buf
}

type request struct {
	function string
	args     []string
	result   chan []byte
}

//AppRunner 合约调用器
type AppRunner struct {
	queue   chan *request
	client  *channel.Client
	sdk     *fabsdk.FabricSDK
	conf    *appConfig
	confile string
}

//NewAppRunner 创建合约调用器
func NewAppRunner(confile string) (*AppRunner, error) {
	o := &AppRunner{
		queue:   make(chan *request, 64),
		client:  nil,
		sdk:     nil,
		confile: confile,
	}
	err := o.initClient(confile)
	if err != nil {
		return nil, err
	}
	go o.messageloop()
	return o, nil
}

//Close close the apprunner
func (o *AppRunner) Close() {
	o.queue <- nil
	if o.sdk != nil {
		o.sdk.Close()
	}
}

// SendRequest call chaincode and wait the result
func (o *AppRunner) SendRequest(function string, args []string) []byte {
	req := &request{
		args:     args,
		function: function,
		result:   make(chan []byte),
	}
	o.queue <- req
	ret := <-req.result
	close(req.result)
	return ret
}

func (o *AppRunner) messageloop() {
	for {
		req := <-o.queue
		if req == nil {
			break
		}
		req.result <- o.callCC(req.function, req.args)
	}
}

func (o *AppRunner) initClient(confile string) error {
	buf, err := ioutil.ReadFile(confile)
	if err != nil {
		return err
	}
	conf := &appConfig{}
	err = json.Unmarshal(buf, conf)
	if err != nil {
		return err
	}
	o.conf = conf
	buf, err = ioutil.ReadFile(conf.ConfigFile)
	if err != nil {
		return err
	}
	opt := config.FromFile(conf.ConfigFile)
	if opt == nil {
		return errors.New("load yaml config  failed")
	}
	sdk, err := fabsdk.New(opt)
	if err != nil {
		return errors.New("Failed to create new SDK:" + err.Error())
	}
	o.sdk = sdk
	clientChannelContext := sdk.ChannelContext(conf.ChannelID, fabsdk.WithUser(conf.OrgAdmin), fabsdk.WithOrg(conf.OrgName))
	client, err := channel.New(clientChannelContext)
	if err != nil {
		o.sdk.Close()
		o.sdk = nil
		return errors.New("Failed to create new channel client: " + err.Error())
	}
	o.client = client
	return nil
}

func (o *AppRunner) callCC(function string, args []string) []byte {
	if o.client == nil {
		err := o.initClient(o.confile)
		if err != nil {
			return getResponse("", "", "", err)
		}
	}
	txArgs := make([][]byte, len(args))
	for i, v := range args {
		txArgs[i] = []byte(v)
	}
	rep, err := o.client.Execute(channel.Request{ChaincodeID: o.conf.ChainCode, Fcn: function, Args: txArgs},
		channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		return getResponse("", "", "", err)
	}
	return getResponse(string(rep.Payload), string(rep.TransactionID), rep.TxValidationCode.String(), nil)
}
