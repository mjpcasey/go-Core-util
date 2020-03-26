package utils

import (
	"fmt"
	"sync"
)

type WaitGroup struct {
	sync.WaitGroup
}

func CreateWaitGroup() *WaitGroup {
	return &WaitGroup{}
}

func (w *WaitGroup) Warp(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}

type WaitGroupWithException struct {
	sync.WaitGroup
	quit chan error
	fns  []func() error
	sync.RWMutex
}

func CreateWaitGroupWithException() *WaitGroupWithException {
	return &WaitGroupWithException{
		quit: make(chan error, 1),
	}
}

func (w *WaitGroupWithException) Warp(fn func() error) {
	w.append(fn)
}

func (w *WaitGroupWithException) Wait() (err error) {

	go func() {
		w.run()
	}()

	return <-w.quit
}

func (w *WaitGroupWithException) append(fn func() error) {
	// 锁的作用主要在于避免执行的同时继续添加fn
	w.Lock()
	w.fns = append(w.fns, fn)
	w.Unlock()
}

func (w *WaitGroupWithException) run() {
	// 锁的作用主要在于避免执行的同时继续添加fn
	w.RLock()
	defer w.RUnlock()
	// 遍历fns，异步调用方法
	for id, fn := range w.fns {
		w.WaitGroup.Add(1)
		go func(f func() error, id int) {
			err := f()
			// 其中一个方法调用错误，即退出当前锁组, 将错误信息返回给quit信号
			if err != nil {
				w.quit <- fmt.Errorf("waitgroup error, fns id: %d, err message: %+v", id, err)
			}

			w.WaitGroup.Done()
		}(fn, id)
	}

	w.WaitGroup.Wait()
	w.quit <- nil
}
