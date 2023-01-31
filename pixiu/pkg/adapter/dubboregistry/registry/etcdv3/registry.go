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
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/common"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/registry"
	baseRegistry "github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/registry/base"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
	"github.com/coreos/etcd/clientv3"
	"github.com/pkg/errors"
	"strings"
	"time"
)

func init() {
	registry.SetRegistry(constant.EtcdV3, newEtcdV3Registry)
}

type EtcdV3Registry struct {
	*baseRegistry.BaseRegistry
	etcdV3Listeners map[registry.RegisteredType]registry.Listener
	client          *clientv3.Client
}

var _ registry.Registry = new(EtcdV3Registry)

func newEtcdV3Registry(regConfig model.Registry, adapterListener common.RegistryEventListener) (registry.Registry, error) {
	timeout, err := time.ParseDuration(regConfig.Timeout)
	if err != nil {
		return nil, errors.Errorf("Incorrect timeout configuration: %s", regConfig.Timeout)
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(regConfig.Address, ","),
		DialTimeout: timeout,
	})
	if err != nil {
		return nil, errors.Errorf("Initialize etcd client failed: %s", err.Error())
	}

	etcdV3Registry := &EtcdV3Registry{
		client:          client,
		etcdV3Listeners: make(map[registry.RegisteredType]registry.Listener),
	}
	etcdV3Registry.etcdV3Listeners[registry.RegisteredTypeInterface] = newEtcdV3IntfListener(client, etcdV3Registry, &regConfig, adapterListener)

	baseReg := baseRegistry.NewBaseRegistry(etcdV3Registry, adapterListener)
	etcdV3Registry.BaseRegistry = baseReg
	return baseReg, nil
}

func (e *EtcdV3Registry) DoSubscribe() error {
	intfListener, ok := e.etcdV3Listeners[registry.RegisteredTypeInterface]
	if !ok {
		return errors.New("Listener for interface level registration does not initialized")
	}
	go intfListener.WatchAndHandle()
	return nil
}

func (e *EtcdV3Registry) DoUnsubscribe() error {
	panic("implement me")
}
