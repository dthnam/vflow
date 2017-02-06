//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    memcache.go
//: details: TODO
//: author:  Mehrdad Arshad Rad
//: date:    02/01/2017
//:
//: Licensed under the Apache License, Version 2.0 (the "License");
//: you may not use this file except in compliance with the License.
//: You may obtain a copy of the License at
//:
//:     http://www.apache.org/licenses/LICENSE-2.0
//:
//: Unless required by applicable law or agreed to in writing, software
//: distributed under the License is distributed on an "AS IS" BASIS,
//: WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//: See the License for the specific language governing permissions and
//: limitations under the License.
//: ----------------------------------------------------------------------------
package ipfix

import (
	"encoding/binary"
	"hash/fnv"
	"net"
	"sync"
	"time"
)

var ShardNo = 32

type MemCache []*TemplatesShard

type Data struct {
	TemplateRecords TemplateRecords
	timestamp       int64
}

type TemplatesShard struct {
	template map[string]Data
	sync.RWMutex
}

func NewCache() MemCache {
	m := make(MemCache, ShardNo)
	for i := 0; i < ShardNo; i++ {
		m[i] = &TemplatesShard{template: make(map[string]Data)}
	}
	return m
}

func (m MemCache) getShard(id uint16, addr net.IP) (*TemplatesShard, []byte) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, id)
	key := append(addr, b...)

	hash := fnv.New32()
	hash.Write(key)
	return m[uint(hash.Sum32())%uint(ShardNo)], key
}

func (m *MemCache) insert(id uint16, addr net.IP, tr TemplateRecords) {
	shard, key := m.getShard(id, addr)
	shard.Lock()
	defer shard.Unlock()
	shard.template[string(key)] = Data{tr, time.Now().Unix()}
}

func (m *MemCache) retrieve(id uint16, addr net.IP) (TemplateRecords, bool) {
	shard, key := m.getShard(id, addr)
	shard.RLock()
	defer shard.RUnlock()
	v, ok := shard.template[string(key)]
	return v.TemplateRecords, ok
}

func (m *MemCache) remove(id int, addr string) {
	// TODO
}

func (m *MemCache) cleanup(id int, addr string) {
	// TODO
}

func (m *MemCache) dump(id int, addr string) {
	// TODO
}
