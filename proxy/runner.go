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

func response(rsp string, err error) string {
	ret := make(map[string]string)
	if err != nil {
		ret["error"] = err.Error()
	} else {
		ret["error"] = "ok"
		ret["data"] = rsp
	}
	buf, _ := json.Marshal(ret)
	return string(buf)
}

type request struct {
	args   []string
	result chan string
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

func (o *AppRunner) SendRequest(args []string) string {
	req := &request{
		args:   args,
		result: make(chan string),
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
		req.result <- response(o.callCC(req.args))
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

func (o *AppRunner) callCC(args []string) (string, error) {
	if len(args) != 3 {
		return "", errors.New("invalid parameters")
	}
	if o.client == nil {
		err := o.initClient()
		if err != nil {
			return "", err
		}
	}
	txArgs := [][]byte{[]byte(args[1]), []byte(args[2])}
	rep, err := o.client.Execute(channel.Request{ChaincodeID: o.conf.ChainCode, Fcn: args[0], Args: txArgs},
		channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		return "", errors.New("failed to call " + args[0] + " error:" + err.Error())
	}
	return string(rep.Payload), nil
}
