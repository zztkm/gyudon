package gyudon

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Commander interface {
	Run([]string) error
}

type spec struct {
	field reflect.StructField
	name  string
	help  string
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

func parseCommand(c Commander) (*Command, error) {
	t := reflect.TypeOf(c)

	if t.Kind() != reflect.Ptr {
		return nil, errors.New("Commander must be pointers to struct")
	}

	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return nil, errors.New("Commander must be pointers to struct")
	}

	n := t.NumField()

	name := t.Name()
	name = strings.ToLower(name)
	cmd := &Command{
		name:      name,
		commander: c,
	}

	for i := 0; i < n; i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue // Export されていないフィールドに対して何もしない
		}
		help, _ := field.Tag.Lookup("help")
		s := &spec{
			field: field,
			name:  field.Name,
			help:  help,
		}
		cmd.specs = append(cmd.specs, s)
	}

	return cmd, nil
}

func NewCommand(c Commander) (*Command, error) {
	return parseCommand(c)
}

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
	cmd, _ := parseCommand(cmder)
	c.addCommand(cmd)
	return cmd
}

func (c *Command) execute(args []string) error {
	err := c.commander.Run(args)
	return err
}

func (c *Command) Execute(args []string) error {
	if c.HasParent() {
		return c.Root().Execute(args)
	}
	args = args[1:]

	var flags []string
	cmd, flags := c.FindCommand(args)
	for i := 0; i < len(cmd.specs); i++ {
		fmt.Println(cmd.specs[i].name)
		fmt.Println(cmd.specs[i].help)
	}

	for i := 0; i < len(flags); i++ {
		for j := 0; j < len(cmd.specs); j++ {
			flag := strings.TrimLeft(flags[i], "-")
			if flag == strings.ToLower(cmd.specs[j].name) {
				fmt.Println("true")
			}
		}
	}

	err := cmd.execute(flags)
	if err != nil {
		return err
	}
	return nil
}
