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
	common2 "github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/common"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/registry"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/remoting/zookeeper"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
	"github.com/coreos/etcd/clientv3"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"sync"
	"time"
)

const (
	MaxFailTimes = 2
	ConnDelay    = 3 * time.Second
)

var _ registry.Listener = new(etcdV3IntfListener)

type etcdV3IntfListener struct {
	exit            chan struct{}
	client          *clientv3.Client
	regConf         *model.Registry
	reg             *EtcdV3Registry
	wg              sync.WaitGroup
	addr            string
	adapterListener common2.RegistryEventListener
	//serviceInfoMap  map[string]*serviceInfo
}

func newEtcdV3IntfListener(client *clientv3.Client, reg *EtcdV3Registry, regConf *model.Registry, adapterListener common2.RegistryEventListener) registry.Listener {
	return &etcdV3IntfListener{
		exit:            make(chan struct{}),
		client:          client,
		regConf:         regConf,
		reg:             reg,
		addr:            regConf.Address,
		adapterListener: adapterListener,
		//serviceInfoMap:  map[string]*serviceInfo{},
	}
}

func (e *etcdV3IntfListener) Close() {
	close(e.exit)
	e.wg.Wait()
}

func (e *etcdV3IntfListener) WatchAndHandle() {
	e.wg.Add(1)
	go e.watch()
}

func (e *etcdV3IntfListener) watch() {
	defer e.wg.Done()
	var (
		failTimes  int64 = 0
		delayTimer       = time.NewTimer(ConnDelay * time.Duration(failTimes))
	)
	defer delayTimer.Stop()
	for {
		serviceList, err := e.client.GetAllServicesInfo(vo.GetAllServiceInfoParam{
			GroupName: e.regConf.Group,
			NameSpace: e.regConf.Namespace,
			PageSize:  100,
		})
		// error handling
		if err != nil {
			failTimes++
			logger.Infof("watching etcd interface with error{%v}", err)
			// Exit the watch if root node is in error
			if err == zookeeper.ErrNilNod {
				logger.Errorf("watching etcd services got errNilNode,so exit listen")
				return
			}
			if failTimes > MaxFailTimes {
				logger.Errorf("Error happens on etcd exceed max fail times: %s,so exit listen", MaxFailTimes)
				return
			}
			delayTimer.Reset(ConnDelay * time.Duration(failTimes))
			<-delayTimer.C
			continue
		}
		failTimes = 0
		if err := e.updateServiceList(serviceList.Doms); err != nil {
			logger.Errorf("update service list failed %s", err)
		}
		time.Sleep(time.Second * 5)
	}
}
