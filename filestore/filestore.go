/*
Copyright 2014 Jiang Le

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package filestore

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Provider interface {
	Init(config string) (err error)
	PutFile(localFileUrl string, remoteFileUrl string) (url string, err error)
	Get(filename string) (data []byte, err error)
	Delete(filename string) (files int, err error)
	List(path string) (entries []Entry, err error)
}

type Entry struct {
	Fname    string `json:"fname"`
	Hash     string `json:"hash"`
	Fsize    int64  `json:"fsize"`
	PutTime  int64  `json:"putTime"`
	MimeType string `json:"mimeType"`
}

type managerConfig struct {
	UrlPrefix      string `json:"urlPrefix"`
	ProviderConfig string `json:"providerConfig"`
}

type Manager struct {
	Provider
	config *managerConfig
}

var provides = make(map[string]Provider)

// Register makes a session provide available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, provide Provider) {
	if provide == nil {
		panic("filestore: Register provide is nil")
	}
	if _, dup := provides[name]; dup {
		panic("filestore: Register called twice for provider " + name)
	}
	provides[name] = provide
}

func NewManager(provideName, config string) (*Manager, error) {
	provider, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("filestore: unknown provide %q (forgotten import?)", provideName)
	}
	cf := new(managerConfig)
	err := json.Unmarshal([]byte(config), cf)
	if err != nil {
		return nil, err
	}
	if cf.UrlPrefix == "" {
		return nil, errors.New("filestore: no urlPrefix provided")
	}
	err = provider.Init(cf.ProviderConfig)
	if err != nil {
		return nil, err
	}

	return &Manager{
		provider,
		cf,
	}, nil
}
