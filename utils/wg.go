package utils

import "sync"

type WgFunc func() interface{}

type Wg struct {
	data  chan interface{}
	wg    *sync.WaitGroup
	funcs []WgFunc
}

func NewWg() *Wg {
	return &Wg{
		data: make(chan interface{}),
		wg:   &sync.WaitGroup{},
	}
}

// Add 添加函数
func (w *Wg) Add(f WgFunc) {
	w.funcs = append(w.funcs, f)
}

func (w *Wg) do() {
	for _, f := range w.funcs {
		w.wg.Add(1)
		go func(f WgFunc) {
			defer w.wg.Done()
			w.data <- f()
		}(f)
	}
}

// Range 执行并发函数
func (w *Wg) Range(f func(v interface{})) {
	w.do()
	go func() {
		defer close(w.data)
		w.wg.Wait()
	}()
	for i := range w.data {
		f(i)
	}
}
