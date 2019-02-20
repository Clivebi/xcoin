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

type Response struct {
	Status       int    `json:"status"`
	ErrorMessage string `json:"errmsg"`
	TxID         string `json:"txid"`
	TxValidCode  string `json:"valid_code"`
	Payload      string `json:"data"`
}

func getResponse(payload string, txID string, Code string, err error) []byte {
	rsp := &Response{
		Status:       200,
		ErrorMessage: "sucess",
		TxID:         txID,
		TxValidCode:  Code,
		Payload:      payload,
	}
	if err != nil {
		rsp.Status = 500
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

type AppRunner struct {
	queue  chan *request
	client *channel.Client
	sdk    *fabsdk.FabricSDK
	conf   *AppConfig
}

func NewAppRunner() (*AppRunner, error) {
	o := &AppRunner{
		queue:  make(chan *request, 64),
		client: nil,
		sdk:    nil,
	}
	err := o.initClient()
	if err != nil {
		return nil, err
	}
	go o.messageloop()
	return o, nil
}

func (o *AppRunner) Close() {
	o.queue <- nil
	if o.sdk != nil {
		o.sdk.Close()
	}
}

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

func (o *AppRunner) initClient() error {
	buf, err := ioutil.ReadFile("./runner.conf")
	if err != nil {
		return err
	}
	conf := &AppConfig{}
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
	var err error = nil
	if o.client == nil {
		err = o.initClient()
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
		return getResponse("", "", "", errors.New("failed to call "+args[0]+" error:"+err.Error()))
	}
	return getResponse(string(rep.Payload), string(rep.TransactionID), rep.TxValidationCode.String(), nil)
}
