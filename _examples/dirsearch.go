package main

import (
	"fmt"

	"github.com/ShangRui-hash/hackflow"
)

func main() {
	targetCh := make(chan string, 1024)
	targetCh <- "https://www.lenovo.com"
	close(targetCh)
	resultCh, err := hackflow.GoDirsearch(&hackflow.GoDirsearchRunConfig{
		TargetCh:            targetCh,
		RoutineCount:        100,
		Dictionary:          "./dict.txt",
		StatusCodeBlackList: []int{404},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	for result := range resultCh {
		fmt.Println(result)
	}
}
