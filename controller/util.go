package controller

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

//封装请求库
func Request(method string,url string,body io.Reader,h map[string]string) (int, []byte){
	client := &http.Client{}
	switch method {
	case "get","GET":
		method = "GET"
	case "post","POST":
		method = "POST"
	case "put","PUT":
		method = "PUT"
	case "delete","DELETE":
		method = "POST"
	}
	r,err := http.NewRequest(method,url,body)
	if err != nil {
		fmt.Println("http错误",err)
	}

	r.Header.Add("Host", "")
	r.Header.Add("Connection","keep-alive")
	r.Header.Add("Accept-Encoding","gzip, deflate, br")
	r.Header.Add("Content-Length","0")
	r.Header.Add("Cache-Control","no-cache")
	for k,v := range h{
		r.Header.Add(k,v)
	}
	response, err := client.Do(r)
	defer response.Body.Close()

	var data []byte = nil

	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ := gzip.NewReader(response.Body)
		for {
			buf := make([]byte, 1024)
			n, err := reader.Read(buf)
			if err != nil && err != io.EOF {
				panic(err)
			}
			if n == 0 {
				break
			}
			data = append(data,buf...)
		}
	default:
		data, _ = ioutil.ReadAll(response.Body)
	}
	return response.StatusCode,data
}