package util

import "sync"

var mtx sync.Mutex

// TODO 实现双重判断避免递归
func Synchronized(fn func()) {
	mtx.Lock()
	defer mtx.Unlock()
	fn()
}