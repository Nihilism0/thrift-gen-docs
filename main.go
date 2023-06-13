package main

import (
	"github.com/cloudwego/thriftgo/plugin"
	"io/ioutil"
	"os"
	"thrift-gen-docs/generate"
)

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		println("Failed to get input:", err.Error())
		os.Exit(1)
	}

	req, err := plugin.UnmarshalRequest(data)
	if err != nil {
		println("Failed to unmarshal request:", err.Error())
		os.Exit(1)
	}

	os.Exit(generate.Run(req))
}
