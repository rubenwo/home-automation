package script

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v.(string))
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

	if len(call.ArgumentList) > 1 {
		obj := call.Argument(1).Object()
		fmt.Println(obj)
		data, err := obj.Value().Export()
		fmt.Println(data)
		if err != nil {
			panic(err)
		}

		headers, ok := data.(map[string]interface{})
		if !ok {
			panic("can't convert headers to map")
		}
		return doRequest(http.MethodGet, uri, nil, headers)
	}

	return doRequest(http.MethodGet, uri, nil, nil)
}

var HttpPost = func(call otto.FunctionCall) otto.Value {
	uri := call.Argument(0).String()
	obj := call.Argument(1).Object()
	fmt.Println(obj)
	data, err := obj.Value().Export()
	fmt.Println(data)
	if err != nil {
		panic(err)
	}

	mappedData, ok := data.(map[string]interface{})
	if !ok {
		panic("can't convert data to map")
	}
	fmt.Println(mappedData)

	if len(call.ArgumentList) > 2 {
		obj := call.Argument(2).Object()
		fmt.Println(obj)
		data, err := obj.Value().Export()
		fmt.Println(data)
		if err != nil {
			panic(err)
		}

		headers, ok := data.(map[string]interface{})
		if !ok {
			panic("can't convert headers to map")
		}
		return doRequest(http.MethodPost, uri, mappedData, headers)
	}

	return doRequest(http.MethodPost, uri, mappedData, nil)
}

var HttpDelete = func(call otto.FunctionCall) otto.Value {
	uri := call.Argument(0).String()
	if len(call.ArgumentList) > 1 {
		obj := call.Argument(1).Object()
		fmt.Println(obj)
		data, err := obj.Value().Export()
		fmt.Println(data)
		if err != nil {
			panic(err)
		}

		headers, ok := data.(map[string]interface{})
		if !ok {
			panic("can't convert headers to map")
		}
		return doRequest(http.MethodDelete, uri, nil, headers)
	}
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

	if len(call.ArgumentList) > 2 {
		obj := call.Argument(2).Object()
		fmt.Println(obj)
		data, err := obj.Value().Export()
		fmt.Println(data)
		if err != nil {
			panic(err)
		}

		headers, ok := data.(map[string]interface{})
		if !ok {
			panic("can't convert headers to map")
		}
		return doRequest(http.MethodPut, uri, mappedData, headers)
	}

	return doRequest(http.MethodPut, uri, mappedData, nil)
}
