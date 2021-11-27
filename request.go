package hackflow

import (
	"net/http"
	"strings"
)

//GetRequest 用来将一个URLCh通道转换成一个RequestCh通道
func GetRequest(URLCh chan string) chan *http.Request {
	requestCh := make(chan *http.Request, 1024)
	go func() {
		for url := range URLCh {
			input(requestCh, url)
		}
		close(requestCh)
	}()
	return requestCh
}

//inputURL 将url输入到管道中
func input(requestCh chan *http.Request, URL string) {
	if !strings.HasPrefix(URL, "http://") && !strings.HasPrefix(URL, "https://") {
		input(requestCh, "http://"+URL)
		input(requestCh, "https://"+URL)
	} else {
		request, err := http.NewRequest("GET", URL, nil)
		if err == nil {
			requestCh <- request
		}
	}
}
