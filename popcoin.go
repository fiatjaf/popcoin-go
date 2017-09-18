package popcoin

import (
	"net/http"
	"time"

	"gopkg.in/jmcvetta/napping.v3"
)

const base = "https://popcoin.ws/api"

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
	_, err := c.Get(base+"/ping", nil, &r, &werr)
	if err != nil {
		return r, err
	}
	if werr.Status == "error" {
		return r, werr
	}
	// pretty.Log(r)
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
	_, err := c.Post(base+"/identify", struct {
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

// Spend consumes some amount of credits from an user account
// and returns the current balances for that user.
//
// https://paper.dropbox.com/doc/Popcoin-API-gdXBTirKxRGeJHgKXOJfl#:uid=374371276125238259264252&h2=POST-Spend
func (c Client) Spend(user string, amount float64, desc string) (SpendResponse, error) {
	r := SpendResponse{}
	werr := Error{}
	_, err := c.Post(base+"/spend", struct {
		User        string  `json:"user"`
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}{user, amount, desc}, &r, &werr)
	if err != nil {
		return r, err
	}
	if werr.Status == "error" {
		return r, werr
	}
	return r, nil
}

type SpendResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	Balances struct {
		Available string `json:"available"`
		Current   string `json:"current"`
	} `json:"dev"`
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
	_, err := c.Get(base+"/spend", &params, &r, &werr)
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
		Id          string  `json:"_id"`
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
		SpentAt     string  `json:"spent_at"`
	} `json:"spends"`
}

type Error struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (err Error) Error() string {
	return err.Message
}
