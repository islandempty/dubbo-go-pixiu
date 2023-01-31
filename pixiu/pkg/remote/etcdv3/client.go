/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package etcdv3

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
	"github.com/coreos/etcd/clientv3"
	"github.com/go-errors/errors"
	model2 "github.com/nacos-group/nacos-sdk-go/model"
	perrors "github.com/pkg/errors"
	"strings"
	"time"
)

type EtcdV3Client struct {
	etcdV3Client *clientv3.Client
}

func (client *EtcdV3Client) GetAllServicesInfo() (model2.ServiceList, error) {

	//return client.namingClient.GetAllServicesInfo(param)
}

func (client *EtcdV3Client) SelectInstances() ([]model2.Instance, error) {
	//return client.namingClient.SelectInstances(param)
}

func (client *EtcdV3Client) Subscribe() error {
	return
}

func (client *EtcdV3Client) Unsubscribe() error {
	return client.etcdV3Client.Delete()
}

func NewEtcdv3Client(config *model.RemoteConfig) (*EtcdV3Client, error) {
	timeout, err := time.ParseDuration(config.Timeout)
	if err != nil {
		return nil, errors.Errorf("Incorrect timeout configuration: %s", config.Timeout)
	}
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(config.Address, ","),
		DialTimeout: timeout,
	})

	if err != nil {
		return nil, perrors.WithMessagef(err, "etcd client create error")
	}
	return &EtcdV3Client{etcdV3Client: client}, nil
}
