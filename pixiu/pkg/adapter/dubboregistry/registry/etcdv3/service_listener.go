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
	"dubbo.apache.org/dubbo-go/v3/common"
	common2 "github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/common"
	"github.com/coreos/etcd/clientv3"
	"sync"
)

type serviceListener struct {
	url    *common.URL
	client *clientv3.Client
	//instanceMap map[string]nacosModel.Instance
	cacheLock sync.Mutex

	exit            chan struct{}
	wg              sync.WaitGroup
	adapterListener common2.RegistryEventListener
}

// WatchAndHandle todo WatchAndHandle is useless for service listener
func (z *serviceListener) WatchAndHandle() {
	panic("implement me")
}

// newEtcdv3SrvListener creates a new etcd service listener
func newEtcdv3SrvListener(url *common.URL, client *clientv3.Client, adapterListener common2.RegistryEventListener) *serviceListener {
	return &serviceListener{
		url:             url,
		client:          client,
		exit:            make(chan struct{}),
		adapterListener: adapterListener,
		//instanceMap:     map[string]nacosModel.Instance{},
	}
}

// Close closes this listener
func (el *serviceListener) Close() {
	close(el.exit)
	el.wg.Wait()
}
