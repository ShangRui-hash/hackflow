package hackflow

import (
	"sync"
)

type Stream struct {
	src     []chan string
	dst     chan string
	filters []func(string) string
}

//NewStream 创建一个新的流
func NewStream() *Stream {
	return &Stream{
		filters: []func(string) string{
			defaultFilter,
		},
		dst: make(chan string, 1024),
	}
}

//defaultFilter 默认过滤器
func defaultFilter(line string) string {
	return line
}

//SetFilter 设置过滤器
func (f *Stream) AddFilter(filter func(string) string) {
	f.filters = append(f.filters, filter)
}

//AddSource 添加一个源
func (f *Stream) AddSrc(src chan string) {
	f.src = append(f.src, src)
}

//GetDst 获取输出流
func (s *Stream) GetDst() chan string {
	var wg sync.WaitGroup
	for i, srcCh := range s.src {
		wg.Add(1)
		go func(i int, srcCh chan string) {
			defer wg.Done()
			for line := range srcCh {
				for _, filter := range s.filters {
					line = filter(line)
				}
				s.dst <- line
			}
		}(i, srcCh)
	}
	go func() {
		logger.Debug("wait for all goroutine")
		wg.Wait()
		close(s.dst)
		logger.Debug("all goroutine done")
	}()
	return s.dst
}
