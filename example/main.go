package main

import (
	"context"
	"flag"

	"github.com/fatelei/cmcc-point/pkg"
)

func main() {
	var ip string
	var privateKey string
	var mobile string
	var cmd string
	var code string
	var sessionID string

	flag.StringVar(&ip, "ip", "127.0.0.1", "ip")
	flag.StringVar(&privateKey, "rsa", "", "rsa key")
	flag.StringVar(&mobile, "mobile", "", "mobile")
	flag.StringVar(&cmd, "cmd", "", "cmd")
	flag.StringVar(&code, "code", "", "code")
	flag.StringVar(&sessionID, "session", "", "session")
	flag.Parse()
	if len(mobile) == 0 {
		println("mobile is requred")
		return
	}

	ctl := cmcc.NewCmcc(ip, privateKey)

	switch cmd {
	case "sms_code":
		err := ctl.SendSmsCode(context.Background(), mobile)
		if err != nil {
			panic(err)
		}
	case "login":
		_, err := ctl.LoginMall(context.Background(), mobile, code)
		if err != nil {
			panic(err)
		}
	case "point":
		point, err := ctl.GetPoints(context.Background(), mobile, sessionID)
		if err != nil {
			panic(err)
		}
		println(point)
	}
}
