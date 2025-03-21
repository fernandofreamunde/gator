package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/fernandofreamunde/gator/internal/commands"
	"github.com/fernandofreamunde/gator/internal/config"
	"github.com/fernandofreamunde/gator/internal/database"
)

func main() {
	fmt.Println("Hello Gator!")
	s := config.State{}
	c := config.Read() // maybe return errors? and detect and exit with os.Exit(1)
	s.Config = &c

	db, err := sql.Open("postgres", c.DbUrl)
	if err != nil {
		fmt.Println("Error: Failed to connect to the database:", err)
	}

	dbQueries := database.New(db)
	s.Db = dbQueries

	cmds := &commands.Commands{
		Registry: make(map[string]func(*config.State, commands.Command) error),
	}

	cmds.Register("login", handlerLogin)
	cmds.Register("register", handlerRegister)
	cmds.Register("reset", handlerReset)
	cmds.Register("users", middlewareLoggedIn(handlerUsers))
	cmds.Register("agg", handlerAgg)
	cmds.Register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.Register("feeds", handlerListFeeds)
	cmds.Register("follow", middlewareLoggedIn(handlerFollowFeed))
	cmds.Register("following", middlewareLoggedIn(handlerFollowing))
	cmds.Register("unfollow", middlewareLoggedIn(handlerUnfollow))

	if len(os.Args) < 2 {
		fmt.Printf("Error: not enough arguments")
		os.Exit(1)
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	cmd := commands.Command{
		Name: cmdName,
		Args: cmdArgs,
	}

	err = cmds.Run(&s, cmd)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func middlewareLoggedIn(handler func(s *config.State, cmd commands.Command, user database.User) error) func(*config.State, commands.Command) error {

	return func(s *config.State, cmd commands.Command) error {
		ctx := context.Background()
		username := sql.NullString{String: s.Config.CurrentUserName, Valid: true}

		user, err := s.Db.GetUserByName(ctx, username)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}

}

func handlerLogin(s *config.State, cmd commands.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("username is requiered")
	}

	ctx := context.Background()
	name := sql.NullString{String: cmd.Args[0], Valid: true}
	existingUser, _ := s.Db.GetUserByName(ctx, name)
	if !existingUser.Name.Valid {
		return fmt.Errorf("user does not exist")
	}
	err := s.Config.SetUser(cmd.Args[0])

	if err != nil {
		return err
	}
	fmt.Println("user has been sucessfuly set")

	return nil
}

func handlerRegister(s *config.State, cmd commands.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("username is requiered")
	}
	ctx := context.Background()
	name := sql.NullString{String: cmd.Args[0], Valid: true}
	existingUser, _ := s.Db.GetUserByName(ctx, name)
	if existingUser.Name.Valid {
		return fmt.Errorf("user already exists")
	}

	user, err := s.Db.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		return err
	}

	fmt.Println("user has been sucessfuly created")

	err = s.Config.SetUser(user.Name.String)
	if err != nil {
		return err
	}
	fmt.Println("user has been sucessfuly set")

	return nil
}

func handlerReset(s *config.State, cmd commands.Command) error {
	ctx := context.Background()

	err := s.Db.ResetUsers(ctx)
	if err != nil {
		return err
	}
	fmt.Println("users have been sucessfuly reset")

	return nil
}

func handlerUsers(s *config.State, cmd commands.Command, user database.User) error {
	ctx := context.Background()

	users, err := s.Db.GetUsers(ctx)
	if err != nil {
		return err
	}

	for _, loopUser := range users {
		if user.Name.String == loopUser.Name.String {

			fmt.Println("* ", loopUser.Name.String, "(current)")
			continue
		}
		fmt.Println("* ", loopUser.Name.String)
	}
	return nil
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response %s", err)
	}

	var feed RSSFeed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response %s", err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	return &feed, nil
}

func handlerAgg(s *config.State, cmd commands.Command) error {

	if len(cmd.Args) < 1 {
		return fmt.Errorf("time between reqs requiered")
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Println("Collecting feeds every", timeBetweenReqs)

	//	feed, err := fetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *config.State) {
	nftf, _ := s.Db.GetNextFeedToFetch(context.Background())

	rss, _ := fetchFeed(context.Background(), nftf.Url.String)

	for _, item := range rss.Channel.Item {
		pubat, err := time.Parse("2006-01-02T15:04:05Z", item.PubDate)
		if err != nil {
			pubat = time.Now()
		}

		_, err = s.Db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			FeedID:      uuid.NullUUID{UUID: nftf.ID, Valid: true},
			Title:       sql.NullString{String: item.Title, Valid: true},
			Url:         sql.NullString{String: item.Link, Valid: true},
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: pubat,
		})

		if err != nil {
			if strings.Contains(err.Error(), "posts_url_key") {
				continue
			}
			fmt.Println(err)
		}
	}

	s.Db.MarkFeedFetched(context.Background(), nftf.ID)
}

func handlerAddFeed(s *config.State, cmd commands.Command, user database.User) error {

	if len(cmd.Args) < 2 {
		return fmt.Errorf("feed name and url are requiered")
	}

	ctx := context.Background()
	name := sql.NullString{String: cmd.Args[0], Valid: true}
	url := sql.NullString{String: cmd.Args[1], Valid: true}

	newFeed, err := s.Db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    uuid.NullUUID{UUID: user.ID, Valid: true},
	})

	if err != nil {
		return err
	}

	_, err = s.Db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    uuid.NullUUID{UUID: user.ID, Valid: true},
		FeedID:    uuid.NullUUID{UUID: newFeed.ID, Valid: true},
	})

	if err != nil {
		return err
	}
	fmt.Println(newFeed)

	return nil
}

func handlerListFeeds(s *config.State, cmd commands.Command) error {

	ctx := context.Background()
	feeds, err := s.Db.GetFeeds(ctx)

	if err != nil {
		return err
	}
	fmt.Println(feeds)
	for _, feed := range feeds {
		fmt.Println(feed.Name)
		fmt.Println(feed.Url)
		fmt.Println(feed.Username)
		fmt.Println("---")
	}

	return nil
}

func handlerFollowing(s *config.State, cmd commands.Command, user database.User) error {

	ctx := context.Background()
	feedFolows, err := s.Db.GetUserFeedFollows(ctx, uuid.NullUUID{UUID: user.ID, Valid: true})
	if err != nil {
		return err
	}

	fmt.Println("Feed Follows:")
	for _, follow := range feedFolows {
		fmt.Println(follow.Name.String)
		fmt.Println(follow.Name_2.String)
		fmt.Println("---")
	}

	return nil
}

func handlerFollowFeed(s *config.State, cmd commands.Command, user database.User) error {

	if len(cmd.Args) < 1 {
		return fmt.Errorf("feed url is requiered")
	}

	url := sql.NullString{String: cmd.Args[0], Valid: true}

	ctx := context.Background()

	feed, err := s.Db.GetFeedByUrl(ctx, url)
	if err != nil {
		return err
	}

	follow, err := s.Db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    uuid.NullUUID{UUID: user.ID, Valid: true},
		FeedID:    uuid.NullUUID{UUID: feed.ID, Valid: true},
	})

	if err != nil {
		return err
	}
	fmt.Println(follow.FeedName.String)
	fmt.Println(follow.UserName.String)

	return nil
}

func handlerUnfollow(s *config.State, cmd commands.Command, user database.User) error {

	if len(cmd.Args) < 1 {
		return fmt.Errorf("feed url is requiered")
	}

	url := sql.NullString{String: cmd.Args[0], Valid: true}
	ctx := context.Background()

	feed, err := s.Db.GetUserFeedFollowByUrl(ctx, database.GetUserFeedFollowByUrlParams{
		UserID: uuid.NullUUID{UUID: user.ID, Valid: true},
		Url:    url,
	})
	if err != nil {
		return err
	}

	err = s.Db.DeleteUserFeedFollow(ctx, feed.ID)

	if err != nil {
		return err
	}

	return nil
}
