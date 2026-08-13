package main

import(
	"fmt"
	"errors"
	"time"
	"context"

	"github.com/google/uuid"

	"github.com/F0hor/database"
)

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handlers[cmd.name]
	if !ok {
		fmt.Printf("Command %s does not exist", cmd.name)
	}

	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("The login command expects one argument (username)")
	}

	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("%s is now loged in\n", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("The register command expects one argument (username)")
	}

	tNow := time.Now()
	user, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID: uuid.New(),
			CreatedAt: tNow,
			UpdatedAt: tNow,
			Name: cmd.args[0],
		},
	)
	if err != nil {
		return err
	}

	fmt.Printf("New user with name %v was created at %v with new uuid: %v\n", user.Name, user.CreatedAt, user.ID)

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("%s is now loged in\n", cmd.args[0])
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return errors.New("The reset command expects zero arguments")
	}

	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("Database has been reset")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return errors.New("The users command expects zero arguments")
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, u := range users {
		if u.Name == s.cfg.User {
			fmt.Printf("* %v (current)\n", u.Name)
		} else {
			fmt.Printf("* %v\n", u.Name)
		}
	}

	return nil
}

func handlerAggregation(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return errors.New("The users command expects zero arguments")
	}

	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Println(feed)

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return errors.New("The addfeed command expects two arguments")
	}

	ctx := context.Background()

	user, err := s.db.GetUser(ctx, s.cfg.User)
	if err != nil {
		return errors.New("Failed to retieve current user from database")
	}

	tNow := time.Now()
	feed, err := s.db.CreateFeed(
		ctx,
		database.CreateFeedParams{
			ID: uuid.New(),
			CreatedAt: tNow,
			UpdatedAt: tNow,
			Name: cmd.args[0],
			Url: cmd.args[1],
			UserID: user.ID,
		},
	)
	if err != nil {
		fmt.Errorf("Failed to create new feed. Reason:\n%v\n", err)
	}

	fmt.Println(feed)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return errors.New("The feeds command expects zero arguments")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to retieve data from database:\n%v\n", err)
	}

	uName := ""
	for _, f := range feeds {
		if uName != f.UserName {
			uName = f.UserName
			fmt.Printf("%v:\n", uName)
		}
		fmt.Printf(" - %v -> %v\n", f.Name, f.Url)
	}

	return nil
}

