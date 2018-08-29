package proxy

import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"io/ioutil"
	"log"
	"strconv"
	"time"
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
	queue chan *request
}

func NewAppRunner() *AppRunner {
	o := &AppRunner{
		queue: make(chan *request, 64),
	}
	go o.messageloop()
	return o
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

func (o *AppRunner) callCC(args []string) (string, error) {
	if len(args) != 3 {
		return "", errors.New("invalid parameters")
	}
	buf, err := ioutil.ReadFile("./runner.conf")
	if err != nil {
		return "", err
	}
	conf := &AppConfig{}
	err = json.Unmarshal(buf, conf)
	if err != nil {
		return "", err
	}
	buf, err = ioutil.ReadFile(conf.ConfigFile)
	if err != nil {
		return "", err
	}
	opt := config.FromFile(conf.ConfigFile)
	if opt == nil {
		return "", errors.New("load yaml config  failed")
	}
	return o.callInvoke(opt, conf, args)
}

func (o *AppRunner) callInvoke(configOpt core.ConfigProvider, conf *AppConfig, args []string) (string, error) {
	sdk, err := fabsdk.New(configOpt)
	if err != nil {
		log.Println("Failed to create new SDK: ", err)
		return "", err
	}
	defer sdk.Close()

	clientChannelContext := sdk.ChannelContext(conf.ChannelID, fabsdk.WithUser(conf.OrgAdmin), fabsdk.WithOrg(conf.OrgName))
	client, err := channel.New(clientChannelContext)
	if err != nil {
		log.Println("Failed to create new channel client: ", err)
		return "", err
	}
	eventID := strconv.Itoa(time.Now().Unix())
	// Register chaincode event (pass in channel which receives event details when the event is complete)
	reg, notifier, err := client.RegisterChaincodeEvent(conf.ChainCode, eventID)
	if err != nil {
		log.Println("Failed to register cc event:", err)
		return "", err
	}
	defer client.UnregisterChaincodeEvent(reg)

	txArgs := [][]byte{[]byte(args[1]), []byte(args[2])}
	rep, err := client.Execute(channel.Request{ChaincodeID: conf.ChainCode, Fcn: args[0], Args: txArgs},
		channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		log.Printf("Failed to call:%s %s\n", args[0], err)
		return "", err
	}
	log.Printf("Received CC event: %#v\n", rep)

	select {
	case ccEvent := <-notifier:
		log.Printf("Received CC event: %#v\n", ccEvent)
		return string(ccEvent.Payload), nil
	case <-time.After(time.Second * 20):
		log.Printf("Did NOT receive CC event for eventId(%s)\n", eventID)
	}
	return "", errors.New("read result time out")
}
