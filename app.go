package main

import (
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "net"
    "os"
    "strconv"
    "time"
)

const (
    layout = "02/Jan/2006:15:04:05 -0700"
)

type LogData struct {
    TimeLocal            string  `json:"time_local"`
    RemoteAddr           string  `json:"remote_addr"`
    RequestMethod        string  `json:"request_method"`
    RequestURI           string  `json:"request_uri"`
    Referrer             string  `json:"referrer"`
    UserAgent            string  `json:"useragent"`
    Host                 string  `json:"host"`
    BytesSent            string  `json:"bytes_sent"`
    Status               string  `json:"status"`
    UpstreamResponseTime string  `json:"upstream_response_time"`
    RequestTime          string  `json:"request_time"`
}

func getEnv(key string) string {
    value, exists := os.LookupEnv(key)
    if !exists {
        log.Fatalf("Environment variable %s not set", key)
    }
    return value
}

func main() {
    host := getEnv("HOST")
    logFile := getEnv("LOG_FILE")
    logSpeedStr := getEnv("LOG_SPEED")
    logSpeed, err := strconv.ParseFloat(logSpeedStr, 64)
    if err != nil {
        log.Fatalf("Invalid log speed: %v", err)
    }
    logInterval := time.Duration(float64(time.Second) / logSpeed)

    methods := []string{"HEAD", "GET", "POST"}
    uris := []string{"/index.html", "/uri", "/test", "/boy", "/buy", "/basket", "/afternoon.html", "/", "/mail", "/index"}
    userAgents := []string{"Mozilla/5.0 (iPad; CPU OS 8_3 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) GSA/11.1.66360 Mobile/12F69 Safari/600.1.4", "Mozilla/5.0 (Linux; Android 5.0.2; SM-A500FU Build/LRX22G) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.83 Mobile Safari/537.36", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.106 Safari/537.36", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.73 Safari/537.36", "Mozilla/5.0 (Windows NT 6.1; Trident/7.0; CCWOW; rv:11.0) like Gecko", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.73 Safari/537.36 OPR/34.0.2036.36", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.106 Safari/537.36", "Mozilla/5.0 (Windows NT 6.0) yi; AppleWebKit/345667.12221 (KHTML, like Gecko) Chrome/23.0.1271.26 Safari/453667.1221", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.106 Safari/537.36", "Mozilla/5.0 (iPad; CPU OS 8_4_1 like Mac OS X) AppleWebKit/537.51.1 (KHTML, like Gecko) GSA/4.1.0.31802 Mobile/12H321 Safari/9537.53"}

    ticker := time.NewTicker(logInterval)
    for range ticker.C {
        logData := LogData{
            TimeLocal:            time.Now().Format(layout),
            RemoteAddr:           getRandomIP(),
            RequestMethod:        getRandomString(methods),
            RequestURI:           getRandomString(uris),
            Referrer:             "www." + getRandomString([]string{"example.com", "example.net", "example.org", "random.com", "random.net", "random.org"}) + "/",
            UserAgent:            getRandomString(userAgents),
            Host:                 host,
            BytesSent:            "585",
            Status:               "200",
            UpstreamResponseTime: fmt.Sprintf("%.3f", rand.Float64()/2),
            RequestTime:          fmt.Sprintf("%.3f", rand.Float64()/2),
        }

        bytes, err := json.Marshal(logData)
        if err != nil {
            log.Fatal(err)
        }
        logLine := string(bytes) + "\n"

        file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
            log.Fatal(err)
        }
        if _, err := file.Write([]byte(logLine)); err != nil {
            log.Fatal(err)
        }
        if err := file.Close(); err != nil {
            log.Fatal(err)
        }
    }
}

func getRandomIP() string {
    b := make([]byte, 4)
    _, err := rand.Read(b)
    if err != nil {
        log.Fatal(err)
    }
    ip := net.IPv4(b[0], b[1], b[2], b[3])
    return ip.String()
}

func getRandomString(s []string) string {
    return s[rand.Intn(len(s))]
}
