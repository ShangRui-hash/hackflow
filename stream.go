package hackflow

import (
	"sync"
)

//Stream 流 用于将一个管道变成多个管道，或者管道变成一个管道，或者将多个管道变成一个管道
type Stream struct {
	src     []chan string
	dst     []chan string
	filters []func(string) string
}

//NewStream 创建一个新的流
func NewStream() *Stream {
	return &Stream{
		filters: []func(string) string{
			defaultFilter,
		},
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

//AddDst 添加一个输出管道
func (f *Stream) AddDst(dst chan string) {
	f.dst = append(f.dst, dst)
}

//GetDst 获取输出流
func (s *Stream) GetDst() []chan string {
	var wg sync.WaitGroup
	for _, srcCh := range s.src {
		wg.Add(1)
		for _, dstCh := range s.dst {
			go func(srcCh, dstCh chan string) {
				defer wg.Done()
				for line := range srcCh {
					for _, filter := range s.filters {
						line = filter(line)
					}
					dstCh <- line
				}
			}(srcCh, dstCh)
		}
	}
	go func() {
		logger.Debug("wait for all goroutine")
		wg.Wait()
		for _, dstCh := range s.dst {
			close(dstCh)
		}
		logger.Debug("all goroutine done")
	}()
	return s.dst
}
