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
	Args        struct{}
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
		name: "help",
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
func (c *Command) FindCommand(args []string) *Command {
	if len(c.commands) == 0 {
		return nil
	}
	if len(args) == 0 {
		return nil
	}
	for _, cmd := range c.commands {
		if cmd.name == args[0] {
			return cmd
		}
	}
	return nil
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

	cmd := c.FindCommand(args)
	if cmd == nil {
		return c.execute(args)
	}
	cmd.execute(args)
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

func main() {
	tkm := &Tkm{}
	root := NewRoot(tkm)

	hello := Hello{}

	root.AddCommand(hello)

	err := root.Execute(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
