package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

type Commander interface {
	Run() error
}

type spec struct {
	field        reflect.StructField
	help         string
	required     bool
	positional   bool
	defaultValue string
}

type Command struct {
	name        string
	short       string
	long        string
	specs       []*spec
	commands    []*Command
	args        []string
	helpCommand *Command
	commander   Commander
	parent      *Command
}

func (c *Command) addHelpCmd() {
	if c.helpCommand != nil {
		return
	}
	help := &Command{
		name:   "help",
		parent: c,
	}
	c.addCommand(help)
}

func NewRoot(c Commander) *Command {
	name := reflect.TypeOf(c).Elem().Name()
	name = strings.ToLower(name)
	cmd := &Command{
		name:      name,
		commander: c,
	}
	return cmd
}

// FindCommand do not use
func (c *Command) FindCommand(args []string) (*Command, []string) {
	var innerfind func(*Command, []string) (*Command, []string)

	innerfind = func(c *Command, innerArgs []string) (*Command, []string) {
		if len(innerArgs) == 0 {
			return c, innerArgs
		}
		if len(c.commands) == 0 {
			return c, innerArgs
		}
		for i := 0; i < len(c.commands); i++ {
			if innerArgs[0] == c.commands[i].name {
				return innerfind(c.commands[i], innerArgs[1:])
			}
		}
		return c, innerArgs
	}

	cmdFound, a := innerfind(c, args)
	return cmdFound, a
}

func (c *Command) HasParent() bool {
	return c.parent != nil
}

func (c *Command) HasSubCommands() bool {
	return len(c.commands) > 0
}

func (c *Command) Root() *Command {
	if c.HasParent() {
		return c.Parent().Root()
	}
	return c
}

func (c *Command) Parent() *Command {
	return c.parent
}

func (c *Command) addCommand(cmds ...*Command) {
	for _, cmd := range cmds {
		cmd.parent = c
		c.commands = append(c.commands, cmd)
	}
}

func (c *Command) AddCommand(cmder Commander) *Command {
	name := reflect.TypeOf(cmder).Elem().Name()
	name = strings.ToLower(name)
	cmd := &Command{
		name:      name,
		commander: cmder,
		parent:    c,
	}
	c.addCommand(cmd)
	return cmd
}

func (c *Command) execute(args []string) error {
	err := c.commander.Run()
	return err
}

func (c *Command) Execute(args []string) error {
	if c.HasParent() {
		return c.Root().Execute(args)
	}
	args = args[1:]

	var flags []string
	cmd, flags := c.FindCommand(args)
	err := cmd.execute(flags)
	if err != nil {
		return err
	}
	return nil
}

// 以下は main 処理

// Tkm is a tkm private command
type Tkm struct {
}

func (c Tkm) Run() error {
	fmt.Println("tkm private command!")
	return nil
}

// Hello return hello your name
type Hello struct {
	// 君の名は。
	Name string `help:"君の名は。"`
}

func (c Hello) Run() error {
	fmt.Printf("Hello %s\n", c.Name)
	return nil
}

// Hoge print Fuga
type Hoge struct {
	Fuga string `default:"fuga"`
}

func (c Hoge) Run() error {
	fmt.Println("hoge", c.Fuga)
	return nil
}

func main() {
	tkm := &Tkm{}
	root := NewRoot(tkm)

	hello := &Hello{}
	hoge := &Hoge{}

	helloCmd := root.AddCommand(hello)
	helloCmd.AddCommand(hoge)

	err := root.Execute(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
