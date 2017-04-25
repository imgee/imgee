// Package main provides ...
package main

import (
	"os"
	"strings"

	"github.com/imgee/imgee/client"
	"github.com/imgee/imgee/config"
	"github.com/imgee/imgee/httpd"
	"github.com/imgee/imgee/plugin"
)

const ver = "v0.0.1"

var (
	req  = make(chan string, 1)
	next = make(chan struct{}, 1)

	cli   *client.Client
	eArgs = os.Args
	noIf  = true
	args  string
)

// init
func init() {
	// load configuration
	config.Init(ver)

	// register imgee help
	plugin.RegisterByCommand("help", help)
	// register imgee version
	plugin.RegisterByCommand("version", version)

	// httpd
	if len(eArgs) == 1 {
		// init client
		cli = client.Init(req, next)
		go cli.Run()

		// start web server
		go httpd.Run()

		noIf = false
	}
}

func main() {
	// perform command
	if noIf {
		cmd := eArgs[1]
		args = strings.Join(eArgs[2:], " ")
		if err := plugin.Call(cmd, args); err != nil {
			println(err.Error())
		}
		return
	}

	// wait command
LOOP:
	for {
		select {
		case request, ok := <-req:
			if !ok {
				break LOOP
			}
			if len(request) == 0 {
				cli.Next()
				continue
			}

			subReq := client.CmdRgx().FindStringSubmatch(request)
			if len(subReq) == 0 {
				println("syntax error")
				cli.Next()
				continue
			}
			args = strings.TrimSpace(subReq[2])
			cmd := strings.TrimSpace(subReq[1])
			if err := plugin.Call(cmd, args); err != nil {
				println(err.Error())
			}

			cli.Next()
		}
	}
}

// imgee help
const usage = `
      ***** TRY IT WITHOUT ANYTHING TO HAVE INTERFACE *****
Usage:
      mylg [command] [args...]

      Available commands:

      resize					resizing image
      version					shows imgee version

Example:
      imgee resize -w 200 -o static/image/a.zip test.jpg
`

func help(args string) error {
	if noIf {
		println(usage)
	} else {
		cli.Help()
	}
	return nil
}

func version(args string) error {
	println(ver)
	return nil
}
