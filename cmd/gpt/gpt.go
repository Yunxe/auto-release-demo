package main

import (
	"auto-release-demo/pkg/gpt"
	"fmt"
	"os"
)

func main() {
	args := os.Args

	// 打印所有参数
	fmt.Println("Number of arguments:", len(args))
	fmt.Println("Arguments:", args)

	// 假设第二个参数是 cURL 响应体
	if len(args) >= 2 {
		curlResponse := args[1]
		gpt.NewGPT(curlResponse)
		fmt.Println("cURL Response:", curlResponse)
	}
}
