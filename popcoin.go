package popcoin

import (
	"net/http"

	"github.com/kr/pretty"

	"gopkg.in/jmcvetta/napping.v3"
)

const base = "https://popcoin.ws/api"

func NewClient(token string) *Client {
	header := http.Header{}
	header.Set("Accept", "application/json")
	header.Set("Content-Type", "application/json")
	header.Set("Authorization", "Bearer "+token)

	return &Client{
		&napping.Session{
			Log:    true,
			Header: &header,
		},
	}
}

type Client struct {
	*napping.Session
}

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
	pretty.Log(r)
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

func (c Client) Spend(user string, amount float64) (SpendResponse, error) {
	r := SpendResponse{}
	werr := Error{}
	_, err := c.Post(base+"/spend", struct {
		User   string  `json:"user"`
		Amount float64 `json:"amount"`
	}{user, amount}, &r, &werr)
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

type Error struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (err Error) Error() string {
	return err.Message
}