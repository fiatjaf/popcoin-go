// popcoin is a wrapper for the https://popcoin.ws/ REST API.
package popcoin

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"gopkg.in/jmcvetta/napping.v3"
)

var BASE = "https://popcoin.ws/api"

// NewClient creates a new Popcoin client.
//
// The Client itself is just a wrapper over the token you provide here.
// It is threadsafe and you can create one for your entire app.
// The token can be found at your Popcoin dashboard.
//
// https://paper.dropbox.com/doc/Popcoin-API-gdXBTirKxRGeJHgKXOJfl#:uid=366194628984948534431988&h2=Authentication
func NewClient(token string) *Client {
	header := http.Header{}
	header.Set("Accept", "application/json")
	header.Set("Content-Type", "application/json")
	header.Set("Authorization", "Bearer "+token)

	return &Client{
		&napping.Session{
			Header: &header,
		},
	}
}

type Client struct {
	*napping.Session
}

// Ping returns basic details about your account.
//
// You can use it to test your token.
//
// https://paper.dropbox.com/doc/Popcoin-API-gdXBTirKxRGeJHgKXOJfl#:uid=067916279625319691097338&h2=Ping
func (c Client) Ping() (PingResponse, error) {
	r := PingResponse{}
	werr := Error{}
	_, err := c.Get(BASE+"/ping", nil, &r, &werr)
	if err != nil {
		return r, err
	}
	if werr.Status == "error" {
		return r, werr
	}
	return r, nil
}

type PingResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Dev     struct {
		Referral  string `json:"referral"`
		Email     string `json:"email"`
		CreatedAt string `json:"created_at"`
	} `json:"dev"`
}

// Identify creates or updates an user with an email address.
//
// The user id can be defined at your discretion.
//
// https://paper.dropbox.com/doc/Popcoin-API-gdXBTirKxRGeJHgKXOJfl#:uid=094648638904648461880497&h2=Identify
func (c Client) Identify(user, email string) (IdentifyResponse, error) {
	r := IdentifyResponse{}
	werr := Error{}
	_, err := c.Post(BASE+"/identify", struct {
		User  string `json:"user"`
		Email string `json:"email"`
	}{user, email}, &r, &werr)
	if err != nil {
		return r, err
	}
	if werr.Status == "error" {
		return r, werr
	}
	return r, nil
}

type IdentifyResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// GetUser retrieves a user data. Email and balance.
//
// https://paper.dropbox.com/doc/Popcoin-API-gdXBTirKxRGeJHgKXOJfl#:uid=196172651128797776297907&h2=GET-Identify
func (c Client) GetUser(user string) (UserResponse, error) {
	r := UserResponse{}
	werr := Error{}

	params := url.Values{
		"user": []string{user},
	}
	_, err := c.Get(BASE+"/identify", &params, &r, &werr)
	if err != nil {
		return r, err
	}
	if werr.Status == "error" {
		return r, werr
	}
	return r, nil
}

type UserResponse struct {
	Status string `json:"status"`
	User   struct {
		Key      string   `json:"key"`
		Email    string   `json:"email"`
		Balances Balances `json:"balances"`
		Link     string   `json:"link"`
	} `json:"user"`
}

// Spend consumes some amount of credits from an user account
// and returns the current balances for that user.
//
// https://paper.dropbox.com/doc/Popcoin-API-gdXBTirKxRGeJHgKXOJfl#:uid=374371276125238259264252&h2=POST-Spend
func (c Client) Spend(user string, amount float64, desc string) (SpendResponse, error) {
	r := SpendResponse{}
	werr := Error{}
	_, err := c.Post(BASE+"/spend", struct {
		User        string      `json:"user"`
		Amount      humbleFloat `json:"amount"`
		Description string      `json:"description"`
	}{user, humbleFloat(amount), desc}, &r, &werr)
	if err != nil {
		return r, err
	}
	if werr.Status == "error" {
		return r, werr
	}
	return r, nil
}

type humbleFloat float64

func (p humbleFloat) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("%.4f", p)
	return []byte(s[0 : len(s)-1]), nil
}

type SpendResponse struct {
	Status   string   `json:"status"`
	Message  string   `json:"message"`
	Balances Balances `json:"balances"`
}

// ListSpends returns a list of spends from an user for a period.
//
// https://paper.dropbox.com/doc/Popcoin-API-gdXBTirKxRGeJHgKXOJfl#:uid=086943532382217583794225&h2=GET-Spend
func (c Client) ListSpends(user string, gte, lte time.Time) (ListSpendsResponse, error) {
	r := ListSpendsResponse{}
	werr := Error{}
	params := napping.Params{
		"user": user,
		"gte":  gte.Format("2006-01-02"),
		"lte":  lte.Format("2006-01-02"),
	}.AsUrlValues()
	_, err := c.Get(BASE+"/spend", &params, &r, &werr)
	if err != nil {
		return r, err
	}
	if werr.Status == "error" {
		return r, werr
	}
	return r, nil
}

type ListSpendsResponse struct {
	Status string `json:"status"`
	Spends []struct {
		Id          string    `json:"_id"`
		Amount      float64   `json:"amount"`
		Description string    `json:"description"`
		SpentAt     time.Time `json:"spent_at"`
	} `json:"spends"`
}

type Error struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Balances struct {
	Available string `json:"available"`
	Current   string `json:"current"`
}

func (err Error) Error() string {
	return err.Message
}
