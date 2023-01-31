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
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/springcloud/servicediscovery"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
	etcdv3 "github.com/dubbogo/gost/database/kv/etcd/v3"
	"github.com/nacos-group/nacos-sdk-go/clients/cache"
	"github.com/nacos-group/nacos-sdk-go/vo"
	perrors "github.com/pkg/errors"
	"sync"
)

type etcdv3ServiceDiscovery struct {
	targetService []string
	//descriptor    string
	client      *etcdv3.Client
	config      *model.RemoteConfig
	listener    servicediscovery.ServiceEventListener
	instanceMap map[string]servicediscovery.ServiceInstance

	cacheLock       sync.Mutex
	callbackFlagMap cache.ConcurrentMap
	//done      chan struct{}
}

func (e *etcdv3ServiceDiscovery) Subscribe() error {
	if e.client == nil {
		return perrors.New("etcd naming client stopped")
	}

	serviceNames := e.listener.GetServiceNames()

	for _, serviceName := range serviceNames {

	}
}

func (e *etcdv3ServiceDiscovery) Unsubscribe() error {
	if e.client == nil {
		return perrors.New("etcd naming client stopped")
	}

	serviceNames := e.listener.GetServiceNames()

	for _, serviceName := range serviceNames {
		subscribeParam := &vo.SubscribeParam{
			ServiceName:       serviceName,
			GroupName:         e.config.Group,
			Clusters:          []string{"DEFAULT"},
			SubscribeCallback: n.Callback,
		}
		_ = e.client.Unsubscribe(subscribeParam)
	}
	return nil
}