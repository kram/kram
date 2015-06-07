// Copyright (c) 2015 The Gus Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

package libraries

import (
	"github.com/zegl/Gus/src/types"
	"io/ioutil"
	"log"
	"net/http"
)

type Library_HTTP struct {
	Library
}

func (self Library_HTTP) Instance() (types.Lib, string) { return &Library_HTTP{}, self.Type() }
func (self Library_HTTP) Type() string                  { return "HTTP" }
func (self Library_HTTP) M_Type() *types.Class          { return self.String(self.Type()) }

func (self Library_HTTP) M_Request(params []*types.Class) *types.Class {
	if len(params) != 2 {
		log.Panic("HTTP::Request() expects exactly 2 parameters")
	}

	// Gus parameters
	verb := params[0].ToString()
	uri := params[1].ToString()

	// Init Go net/http library
	client := &http.Client{}
	req, err := http.NewRequest(verb, uri, nil)

	self.Error(err)

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
