package main

import (
	"fmt"
	"log"
	"os"

	"github.com/zztkm/gyudon"
)

// Tkm is a tkm private command
type Tkm struct {
}

func (c Tkm) Run(args []string) error {
	fmt.Println("tkm private command!")
	return nil
}

// Hello return hello your name
type Hello struct {
	// 君の名は。
	Name string `help:"君の名は。"`
}

func (c Hello) Run(args []string) error {
	fmt.Printf("Hello %s\n", c.Name)
	return nil
}

// Hoge print Fuga
type Hoge struct {
	Fuga string `help:"fugafuga" default:"fuga"`
}

func (c Hoge) Run(args []string) error {
	fmt.Println("hoge", c.Fuga)
	return nil
}

func main() {
	tkm := &Tkm{}
	root, err := gyudon.NewCommand(tkm)
	if err != nil {
		log.Fatal(err)
	}

	hello := &Hello{}
	hoge := &Hoge{}

	root.AddCommand(hello)
	root.AddCommand(hoge)

	err = root.Execute(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
