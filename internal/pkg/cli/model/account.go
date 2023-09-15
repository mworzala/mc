package model

import "fmt"

type Account struct {
	UUID     string
	Username string
}

func (a *Account) String() string {
	return a.Username
}

type LoginPrompt struct {
	Url  string
	Code string
}

func (p *LoginPrompt) String() string {
	return fmt.Sprintf("%s %s", p.Url, p.Code)
}

type AccessToken struct {
	Username string
	UUID     string
	Token    string
}

func (t *AccessToken) String() string {
	return t.Token
}
