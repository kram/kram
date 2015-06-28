// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

package libraries

import (
	"github.com/kram/kram/src/types"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Library_HTTP struct {
	Library

	request *http.Request
	verb    string
	uri     string
	body    io.Reader
	headers map[string]string
}

func (self Library_HTTP) Instance() (types.Lib, string) { return &Library_HTTP{}, self.Type() }
func (self Library_HTTP) Type() string                  { return "HTTP" }
func (self Library_HTTP) M_Type() *types.Class          { return self.String(self.Type()) }

func (self *Library_HTTP) M_Request(params []*types.Class) *types.Class {

	// Default parameters
	verb := "GET"
	uri := "http://localhost"

	// Parse params
	self.Params(params, &verb, &uri)

	self.verb = verb
	self.uri = uri
	self.headers = make(map[string]string)
	self.body = strings.NewReader("")

	return self.Null()
}

func (self *Library_HTTP) M_Header(params []*types.Class) *types.Class {
	key := ""
	value := ""

	self.Params(params, &key, &value)

	self.headers[key] = value

	return self.Null()
}

func (self *Library_HTTP) M_Body(params []*types.Class) *types.Class {
	body := ""

	self.Params(params, &body)

	self.body = strings.NewReader(body)

	return self.Null()
}

func (self *Library_HTTP) M_Do(params []*types.Class) *types.Class {

	// Init Go net/http library
	req, err := http.NewRequest(self.verb, self.uri, self.body)

	self.Error(err)

	req.Header.Set("User-Agent", "Gus/HTTP Library 0.1")

	for key, val := range self.headers {
		req.Header.Set(key, val)
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	self.Error(err)

	response, err := ioutil.ReadAll(resp.Body)

	self.Error(err)

	return self.String(string(response))
}

func (self Library_HTTP) Error(err error) {
	if err != nil {
		log.Panicf("HTTP::Request() encountered an error: %s", err)
	}
}
