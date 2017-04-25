// Package client provides ...
package client

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/chzyer/readline"
	"github.com/imgee/imgee/banner"
	"github.com/imgee/imgee/config"
	"github.com/imgee/imgee/plugin"
)

const (
	usage = `
Usage:
	The imgee tool, Processing images to fit the application.
	The vi/emacs mode, almost all basic features are supported. Press tab to see which options are available.

	resize                      image resizing with common interpolation methods.


	Please visit http://imgee.me/docs for more information
	`
)

type Client struct {
	next     chan struct{}
	cmd      chan<- string
	readline *readline.Instance
	prompt   string
}

func Init(cmd chan<- string, next chan struct{}) (cli *Client) {
	cli = &Client{
		cmd:    cmd,
		next:   next,
		prompt: "imgee",
	}

	// readline
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          Colorize("imgee"),
		HistoryFile:     "/tmp/.imgee.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		AutoComplete:    readline.NewPrefixCompleter(plugin.Completers()...),
	})
	if err != nil {
		panic(err)
	}
	cli.readline = rl

	// check update
	go checkUpdate()

	// print banner
	banner.Println(config.Conf.Version)

	return
}

// run
func (cli *Client) Run() {
	defer close(cli.cmd)

LOOP:
	for {
		line, err := cli.readline.Readline()
		if err != nil { // io.EOF, readline.ErrInterrupt
			switch err {
			case io.EOF:
				break LOOP
			case readline.ErrInterrupt:
				if len(line) == 0 {
					break LOOP
				}
				continue
			default:
				println(err.Error())
				break LOOP
			}
		}

		cli.cmd <- line
		if _, ok := <-cli.next; !ok {
			break
		}
	}
}

// SetPrompt set readline prompt and store it
func Colorize(p string) string {
	return fmt.Sprintf("\033[92m%s» \033[0m", strings.ToLower(p))
}
func (cli *Client) SetPrompt(p string) {
	cli.readline.SetPrompt(Colorize(p))
}

// Refresh prompt
func (cli *Client) Refresh() {
	cli.readline.Refresh()
}

// SetVim set mode to vim
func (cli *Client) SetVim() {
	if !cli.readline.IsVimMode() {
		cli.readline.SetVimMode(true)
		println("mode changed to vim")
	} else {
		println("mode already is vim")
	}
}

// SetEmacs set mode to emacs
func (cli *Client) SetEmacs() {
	if cli.readline.IsVimMode() {
		cli.readline.SetVimMode(false)
		println("mode changed to emacs")
	} else {
		println("mode already is emacs")
	}
}

func (cli *Client) Next() {
	cli.next <- struct{}{}
}

// Close the readline instance
func (cli *Client) Close(next chan struct{}) {
	cli.readline.Close()
}

// client help
func (cli *Client) Help() {
	println(usage)
}

// CmdRgx returns commands regex for validation
func CmdRgx() *regexp.Regexp {
	expr := fmt.Sprintf(`(%s)\s{0,1}(.*)`, strings.Join(plugin.Cmds(), "|"))
	re, _ := regexp.Compile(expr)
	return re
}

// checkUpdate checks if any new version is available
func checkUpdate() {

}
