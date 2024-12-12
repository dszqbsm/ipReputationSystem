package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 从AbuseIPDB数据库查询指定IP的信誉分数并返回
func abuseipdbCheck(ip string) (float64, error) {
	// 数据准备
	url := "https://api.abuseipdb.com/api/v2/check"
	queryParams := fmt.Sprintf("ipAddress=%s&maxAgeInDays=90", ip)
	headers := map[string]string{
		"Accept":       "application/json",
		"Key":          "e721369ffba7384c0b6c98b1a46e0b29a5b236b7cdfc84983f8e6bf3fae03c9567c06ea0b80871e2",
		"Content-Type": "application/json",
	}

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return -1, err
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
		return -1, err
	}
	defer resp.Body.Close()

	// 处理响应
	body, err := ioutil.ReadAll(resp.Body) // 读取响应
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return -1, err
	}
	var decodeResponse map[string]interface{}
	err = json.Unmarshal(body, &decodeResponse) // 反序列化JSON
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return -1, err
	}

	// 提取abuseConfidenceScore
	data, ok := decodeResponse["data"].(map[string]interface{})
	if !ok {
		fmt.Println("Error: data field is not a map")
		return -1, err
	}
	abuseConfidenceScore, ok := data["abuseConfidenceScore"].(float64)
	if !ok {
		fmt.Println("Error: abuseConfidenceScore field is not a float64")
		return -1, err
	}
	return abuseConfidenceScore, nil
}

func main() {
	r := gin.Default()
	// 请求示例：http://127.0.0.1:8090/?ip=127.0.0.1
	r.GET("/", func(c *gin.Context) {
		ip := c.Query("ip")
		score, err := abuseipdbCheck(ip)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"score": score})
		}
	})
	r.Run(":8090")
}
