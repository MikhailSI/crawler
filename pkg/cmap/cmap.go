package cmap

import (
	"crypto/md5"
	"sync"
)

type CMap struct {
	sync.RWMutex
	mp map[[16]byte]bool
}

func NewCMap() *CMap {
	return &CMap{
		mp: make(map[[16]byte]bool),
	}
}

func (cm *CMap) CheckAdd(url string) bool {
	cm.Lock()
	defer cm.Unlock()

	key := md5.Sum([]byte(url))
	_, ok := cm.mp[key]
	if !ok {
		cm.mp[key] = true
		return false
	}

	return true
}
