package main

import (
	"bytes"
	"os"
	"time"

	"github.com/tendermint/abci/client"
	"github.com/tendermint/abci/types"
	common "github.com/tendermint/go-common"
	"github.com/tendermint/go-process"
)

//----------------------------------------

func StartApp(abciApp string) *process.Process {
	// Start the app
	//outBuf := NewBufferCloser(nil)
	proc, err := process.StartProcess("abci_app",
		"",
		"bash",
		[]string{"-c", abciApp},
		nil,
		os.Stdout,
	)
	if err != nil {
		panic("running abci_app: " + err.Error())
	}

	// TODO a better way to handle this?
	time.Sleep(time.Second)

	return proc
}

func StartClient(abciType string) abcicli.Client {
	// Start client
	client, err := abcicli.NewClient("tcp://127.0.0.1:46658", abciType, true)
	if err != nil {
		panic("connecting to abci_app: " + err.Error())
	}
	return client
}

func SetOption(client abcicli.Client, key, value string) {
	res := client.SetOptionSync(key, value)
	_, _, log := res.Code, res.Data, res.Log
	if res.IsErr() {
		panic(fmt.Sprintf("setting %v=%v: \nlog: %v", key, value, log))
	}
}

func Commit(client abcicli.Client, hashExp []byte) {
	res := client.CommitSync()
	_, data, log := res.Code, res.Data, res.Log
	if res.IsErr() {
		panic(fmt.Sprintf("committing %v\nlog: %v", log))
	}
	if !bytes.Equal(res.Data, hashExp) {
		panic(fmt.Sprintf("Commit hash was unexpected. Got %X expected %X",
			data, hashExp))
	}
}

func DeliverTx(client abcicli.Client, txBytes []byte, codeExp types.CodeType, dataExp []byte) {
	res := client.DeliverTxSync(txBytes)
	code, data, log := res.Code, res.Data, res.Log
	if code != codeExp {
		panic(fmt.Sprintf("DeliverTx response code was unexpected. Got %v expected %v. Log: %v",
			code, codeExp, log))
	}
	if !bytes.Equal(data, dataExp) {
		panic(fmt.Sprintf("DeliverTx response data was unexpected. Got %X expected %X",
			data, dataExp))
	}
}

func CheckTx(client abcicli.Client, txBytes []byte, codeExp types.CodeType, dataExp []byte) {
	res := client.CheckTxSync(txBytes)
	code, data, log := res.Code, res.Data, res.Log
	if res.IsErr() {
		panic(fmt.Sprintf("checking tx %X: %v\nlog: %v", txBytes, log))
	}
	if code != codeExp {
		panic(fmt.Sprintf("CheckTx response code was unexpected. Got %v expected %v. Log: %v",
			code, codeExp, log))
	}
	if !bytes.Equal(data, dataExp) {
		panic(fmt.Sprintf("CheckTx response data was unexpected. Got %X expected %X",
			data, dataExp))
	}
}
