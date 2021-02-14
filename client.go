package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func main() {
	values := url.Values{}
	values.Set("server_name", "test")
	values.Add("log_name", "temperature")
	values.Add("ip_address", "127.0.0.1")
	values.Add("mac_address", "hogeeee")
	values.Add("other_info", "koreha test")
	values.Add("value", "30.2")
	t := time.Now()
	values.Add("date_point", t.Format("2006-01-02 15:04:05"))
	req, err := http.NewRequest(
		"POST",
		"http://127.0.0.1:3000/api/v1/server_logs/create",
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	var buf []byte
	_, err = resp.Body.Read(buf)
	if err != nil {
		panic(err)
	}

	err = resp.Body.Close()

	if err != nil {
		panic(err)
	}

	fmt.Println(string(buf))
	fmt.Println("finished.")
}
