package literaladapter

import (
	"context"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

// Auth is the authentication information for the literal.club API.
type Auth struct {
	Email    string `env:"LITERAL_EMAIL"`
	Password string `env:"LITERAL_PASSWORD"`
}

const literalURL = "https://literal.club/graphql/"

func login() (*graphql.Client, error) {
	var auth Auth
	if err := env.Parse(&auth); err != nil {
		return nil, err
	}

	client := graphql.NewClient(literalURL, http.DefaultClient)
	m := loginM{}
	if err := client.Mutate(context.Background(), &m, map[string]interface{}{
		"email":    graphql.String(auth.Email),
		"password": graphql.String(auth.Password),
	}); err != nil {
		return nil, err
	}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: string(m.Login.Token)},
	)
	cli := oauth2.NewClient(context.Background(), src)
	return graphql.NewClient(literalURL, cli), nil
}

// CurrentlyReading retrieves the currently reading list.
func CurrentlyReading() ([]Book, error) {
	client, err := login()
	if err != nil {
		return nil, err
	}

	q := readingQ{}
	if err := client.Query(context.Background(), &q, nil); err != nil {
		return nil, err
	}

	var books []Book
	for _, rs := range q.MyReadingStates {
		if rs.ReadingState.Status != "IS_READING" {
			continue
		}
		book := rs.ReadingState.Book
		books = append(books, book)
	}
	return books, nil
}

// Book is a book.
type Book struct {
	Slug        graphql.String
	Title       graphql.String
	Subtitle    graphql.String
	Description graphql.String
	Authors     []Author
}

// Author is an author.
type Author struct {
	Name graphql.String
}

type loginM struct {
	Login struct {
		Token graphql.String
	} `graphql:"login(email: $email, password: $password)"`
}

type readingQ struct {
	MyReadingStates []struct {
		ReadingState struct {
			Status graphql.String
			Book   Book
		} `graphql:"... on ReadingState"`
	}
}
