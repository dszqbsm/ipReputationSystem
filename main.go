package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// 数据准备
	url := "https://api.abuseipdb.com/api/v2/check"
	queryParams := "ipAddress=118.25.6.39&maxAgeInDays=90"
	headers := map[string]string{
		"Accept":       "application/json",
		"Key":          "e721369ffba7384c0b6c98b1a46e0b29a5b236b7cdfc84983f8e6bf3fae03c9567c06ea0b80871e2",
		"Content-Type": "application/json",
	}

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}
	req.URL.RawQuery = queryParams    // 添加参数
	for key, value := range headers { // 添加请求头
		req.Header.Set(key, value)
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 处理响应
	body, err := ioutil.ReadAll(resp.Body) // 读取响应
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}
	var decodeResponse map[string]interface{}
	err = json.Unmarshal(body, &decodeResponse) // 反序列化JSON
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	// 打印响应
	output, err := json.MarshalIndent(decodeResponse, "", "  ")
	if err != nil {
		fmt.Printf("Error formatting JSON: %v\n", err)
		return
	}

	fmt.Println(string(output))
}
