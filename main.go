package main

import(
	"fmt"
	"os"

	"github.com/F0hor/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("You must enter desired command name")
		os.Exit(1)
	}

	conf := config.Read()
	myState := state{
		cfg: &conf,
	}

	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}
	initCmds(&cmds)

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	err := cmds.run(&myState, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initCmds(cmds *commands) {
	cmds.register("login", handlerLogin)
}

