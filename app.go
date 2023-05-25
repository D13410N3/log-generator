package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "math/rand"
    "net"
    "os"
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

func main() {
    host := os.Getenv("HOST")
    logFile := os.Getenv("LOG_FILE")
    // Log speed is in logs per second
    logSpeed, err := time.ParseDuration(os.Getenv("LOG_SPEED") + "s")
    if err != nil {
        log.Fatalf("invalid log speed: %v", err)
    }

    methods := []string{"HEAD", "GET", "POST"}
    uris := []string{"/index.html", "/uri", "/test"}
    userAgents := []string{"UA-1", "UA-2", "Test-UA", "Mozilla"}

    ticker := time.NewTicker(logSpeed)
    for range ticker.C {
        logData := LogData{
            TimeLocal:            time.Now().Format(layout),
            RemoteAddr:           getRandomIP(),
            RequestMethod:        getRandomString(methods),
            RequestURI:           getRandomString(uris),
            Referrer:             "www." + getRandomString([]string{"example.com", "example.net", "example.org"}) + "/",
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
        err = ioutil.WriteFile(logFile, []byte(logLine), 0644)
        if err != nil {
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
