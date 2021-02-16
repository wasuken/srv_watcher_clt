package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func getPacmanInfo() {

}

// ここのパクリじゃん。
// https://stackoverflow.com/questions/11356330/how-to-get-cpu-usage
func getCPUInfo() (idle, total uint64) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val
				if i == 4 {
					idle = val
				}
			}
			return
		}
	}
	return
}

func calcCPUPercent() float64 {
	idle0, total0 := getCPUInfo()
	time.Sleep(3 * time.Second)
	idle1, total1 := getCPUInfo()

	idleTicks := float64(idle1 - idle0)
	totalTicks := float64(total1 - total0)
	return 100 * (totalTicks - idleTicks) / totalTicks
}

func getSrvIfaceInfo() (error, string, string, string) {
	ifs, err := net.Interfaces()
	if err != nil {
		return err, "", "", ""
	}
	for _, ifr := range ifs[1:] {
		ips, err := ifr.Addrs()
		if err != nil {
			return err, "", "", ""
		}
		for _, ipa := range ips {
			if len(strings.Split(ipa.String(), ".")) == 4 {
				return nil, ifr.Name, ipa.String(), ifr.HardwareAddr.String()
			}
		}
	}
	// 別に取得できなくても良い。
	return nil, "", "", ""
}

func sendServer() {
	err, ifName, ipa, mac := getSrvIfaceInfo()
	if err != nil {
		panic(err)
	}
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	value := calcCPUPercent()

	values := url.Values{}
	values.Set("server_name", name)
	values.Add("log_name", "cpuInfo")
	values.Add("ip_address", ipa)
	values.Add("mac_address", mac)
	values.Add("other_info", ifName)
	values.Add("value", fmt.Sprintf("%f", value))
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
}

func main() {
	sendServer()
	fmt.Println("finished.")
}
