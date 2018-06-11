package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/medalon/stats_from_nginx/stats"
)

var format = "main"
var logFile string

func main() {
	// Read given file or from STDIN
	var logReader io.Reader
	logFile = "dummy"

	if logFile == "dummy" {
		logReader = strings.NewReader(`194.152.36.162 - - [26/Feb/2018:16:54:10 +0600] "GET /1px.png?stat=preroll&name=Beeline_PACMAN&act=click&pid=bQFrkd9E HTTP/1.1" 200 923 "https://www.super.kg/embed/video/183223?mesto=1&ismobile=1" "Mozilla/5.0 (Linux; Android 6.0.1; SM-G532F Build/MMB29T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Mobile Safari/537.36" "-"
		91.192.66.244 - - [27/Feb/2018:04:22:24 +0600] "GET /1px.png?stat=preroll&name=Beeline_PACMAN&act=show&pid=SeD3e3Ra HTTP/1.1" 200 923 "https://www.super.kg/embed/video/210837?mesto=0&ismobile=1" "Mozilla/5.0 (Linux; U; Android 6.0.1; en-US; ASUS_Z010DD Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/40.0.2214.89 UCBrowser/11.5.0.1015 Mobile Safari/537.36" "-"`)
	} else if logFile == "-" {
		logReader = os.Stdin
	} else {
		file, err := os.Open(logFile)
		if err != nil {
			panic(err)
		}
		logReader = file
		defer file.Close()
	}

	// Use nginx config file to extract format by the name
	var nginxConfig io.Reader
	nginxConfig = strings.NewReader(`
            http {
                log_format   main  '$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" "$http_x_forwarded_for"';
            }
		`)

	testRes := stats.ParseLogFile(logReader, nginxConfig, format)
	fmt.Println(testRes)
}
