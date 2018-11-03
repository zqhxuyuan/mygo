package main

import (
	"fmt"
	"sort"
	"sync"
)

type KV struct {
	count int
	keys  []string
	hash  map[string]interface{}
	lock  sync.RWMutex
}

// 添加kv键值对
func (this *KV) Set(k string, v interface{}) {
	this.lock.Lock()
	if _, ok := this.hash[k]; !ok {
		this.keys = append(this.keys, k)
		sort.Strings(this.keys)
		this.count++
	}
	this.hash[k] = v
	this.lock.Unlock()
}

// 获取数据长度
func (this *KV) Count() int {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.count
}

// 由key检索value
func (this *KV) Get(k string) (interface{}, bool) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	v, ok := this.hash[k]
	return v, ok
}

// 根据key排序，返回有序的vaule切片
func (this *KV) Values() []interface{} {
	this.lock.RLock()
	defer this.lock.RUnlock()
	vals := make([]interface{}, this.count)
	for i := 0; i < this.count; i++ {
		vals[i] = this.hash[this.keys[i]]
	}
	return vals
}

func main() {
	fmt.Println("kv...")
}
