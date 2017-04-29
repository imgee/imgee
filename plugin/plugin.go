// Package plugin provides ...
package plugin

import (
	"errors"

	"github.com/chzyer/readline"
	"github.com/imgee/imgee/config"
	"github.com/imgee/imgee/resize"
)

var (
	plugins    = make(map[string]Plugin)
	completers []readline.PrefixCompleterInterface
	cmds       []string
)

type Plugin interface {
	Exec(args string) error
	Cmds() []string
	Command() string
}

// register plugins
func init() {
	Register(config.Conf)
	Register(&resize.Resize{})
}

// register plugins
func Register(p Plugin) {
	if _, ok := plugins[p.Command()]; ok {
		panic("the plugin already registered")
	} else {
		plugins[p.Command()] = p
	}

	cmds = append(cmds, p.Command())
	pc := readline.PcItem(p.Command())
	if p.Cmds != nil {
		cs := make([]readline.PrefixCompleterInterface, len(p.Cmds()))
		for i, v := range p.Cmds() {
			cs[i] = readline.PcItem(v)
		}
		pc.SetChildren(cs)
	}
	completers = append(completers, pc)
}

// register plugin by name
func RegisterByCommand(c string, f func(args string) error) {
	Register(&emptyPlugin{
		command: c,
		exec:    f,
	})
}

// call plugins
func Call(cmd, args string) error {
	p, ok := plugins[cmd]
	if !ok {
		return errors.New("Invalid command please try help")
	}

	return p.Exec(args)
}

// cmds
func Cmds() []string {
	return cmds
}

// completer
func Completers() []readline.PrefixCompleterInterface {
	return completers
}

// empty plugin
type emptyPlugin struct {
	command string
	exec    func(args string) error
}

func (p *emptyPlugin) Exec(args string) error { return p.exec(args) }
func (p *emptyPlugin) Cmds() []string         { return nil }
func (p *emptyPlugin) Command() string        { return p.command }
