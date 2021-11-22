package main

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/sirupsen/logrus"
)

func main() {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.CheckRetry = func(ctx *retryablehttp.RequestCtx, resp *retryablehttp.Response, err error) (bool, error) {
		text := ioutil.ReadAll(resp.Body)
		if strings.Contains(string(text), "error") {
			return true, fmt.Errorf("error")
		}
		if resp != nil && resp.StatusCode == 429 {
			return true, nil
		}
		return false, nil
	}
	standardClient := retryClient.StandardClient() // *http.Client
	resp, err := standardClient.Get("https://xdu.databankes.cn")
	if err != nil {
		logrus.Error("standardClient.Get failed,err:", err)
		return
	}
	text, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("ioutil.ReadAll failed,err:", err)
		return
	}
	fmt.Println(text)
}
