package main

import _ "github.com/lib/pq"

import(
	"fmt"
	"os"
	"database/sql"

	"github.com/F0hor/config"
	"github.com/F0hor/database"
)

type state struct {
	cfg *config.Config
	db *database.Queries
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("You must enter desired command name")
		os.Exit(1)
	}

	conf := config.Read()
	db, err := sql.Open("postgres", conf.DbUrl)
	dbQueries := database.New(db)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	myState := state{
		cfg: &conf,
		db: dbQueries,
	}

	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}
	initCmds(&cmds)

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	err = cmds.run(&myState, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initCmds(cmds *commands) {
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAggregation)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", handlerFollowing)
	cmds.register("unfollow", handlerUnfollow)
}
