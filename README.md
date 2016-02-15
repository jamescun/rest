Rest
====

[![GoDoc](https://godoc.org/github.com/jamescun/rest?status.svg)](https://godoc.org/github.com/jamescun/rest) [![License](https://img.shields.io/badge/license-BSD-blue.svg)](LICENSE)


Rest makes the creation of client libraries for REST-like APIs simpler. It is not intended to be used directly, but serve as the foundation of a client library which implements the required models.

    go get github.com/jamescun/rest


Example
-------

```go
// define a new JSON API exposed at https://api.example.org
api, _ := rest.New("http://api.example.org", rest.EncoderJSON{})

// add authorization header to all requests
api.Before(func(r *http.Request) {
	r.Header.Set("Authorization", "Token supersecretpassword")
})

// define requests that can be made against API
ListUsers := api.Request("GET", "/1/users")
CreateUser := api.Request("POST", "/1/users")


// create a user with our above api
// equivelant to raw "POST /1/users" with body
// response of request will be unmarshaled from JSON into user variable
var user User
err := CreateUser(&user, User{Name: "James"})
if err != nil {
	fmt.Println("error: could not create user:", err)
}
fmt.Printf("info: created user: %+v\n", user)

// list all users with our above api
// equivelant to raw "GET /1/users?limit=10&filter=James"
// response of request will be unmarshaled from JSON into users variable
var users []User
err = ListUsers(&users, rest.Limit(10), rest.Param("filter", "James"))
if err != nil {
	fmt.Println("error: could not list users:", err)
}
fmt.Printf("info: user list: %+v\n", users)
```

See the [examples](examples) directory for more.
