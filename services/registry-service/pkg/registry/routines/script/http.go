package script

import (
	"bytes"
	"encoding/json"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"net/http"
)

type JsonError struct {
	Msg string `json:"msg"`
}

func doRequest(method, uri string, data map[string]interface{}, headers map[string]interface{}) otto.Value {
	client := &http.Client{}

	var (
		req *http.Request
		err error
	)

	if data == nil {
		req, err = http.NewRequest(method, uri, nil)
		if err != nil {
			result, _ := otto.ToValue(&JsonError{Msg: err.Error()})
			return result
		}
	} else {
		b, err := json.Marshal(data)
		if err != nil {
			result, _ := otto.ToValue(&JsonError{Msg: err.Error()})
			return result
		}
		req, err = http.NewRequest(method, uri, bytes.NewBuffer(b))
		if err != nil {
			result, _ := otto.ToValue(&JsonError{Msg: err.Error()})
			return result
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		result, _ := otto.ToValue(&JsonError{Msg: err.Error()})
		return result
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		result, _ := otto.ToValue(&JsonError{Msg: err.Error()})
		return result
	}
	resp.Body.Close()

	result, _ := otto.ToValue(string(body))

	return result
}

var HttpGet = func(call otto.FunctionCall) otto.Value {
	uri := call.Argument(0).String()
	return doRequest(http.MethodGet, uri, nil, nil)
}

var HttpPost = func(call otto.FunctionCall) otto.Value {
	uri := call.Argument(0).String()
	obj := call.Argument(1).Object()
	data, err := obj.Value().Export()
	if err != nil {
		panic(err)
	}

	mappedData, ok := data.(map[string]interface{})
	if !ok {
		panic("can't convert data to map")
	}
	return doRequest(http.MethodPost, uri, mappedData, nil)
}

var HttpDelete = func(call otto.FunctionCall) otto.Value {
	uri := call.Argument(0).String()
	return doRequest(http.MethodDelete, uri, nil, nil)
}

var HttpPut = func(call otto.FunctionCall) otto.Value {
	uri := call.Argument(0).String()
	obj := call.Argument(1).Object()
	data, err := obj.Value().Export()
	if err != nil {
		panic(err)
	}

	mappedData, ok := data.(map[string]interface{})
	if !ok {
		panic("can't convert data to map")
	}

	return doRequest(http.MethodPut, uri, mappedData, nil)
}
