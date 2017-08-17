package main

import "github.com/ajg/form"


type User struct {
	Name         string            `form:"name"`
	Email        string            `form:"email"`
}

func PostUser(url string, u User) error {
	var c http.Client
	_, err := c.PostForm(url, form.EncodeToValues(u))
	return err
}

func main() {
gw:= User{Name: "rakesh", Email: "abc"}
PostUser("https://127.0.0.1/api/v1", gw)
}
