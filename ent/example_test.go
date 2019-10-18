// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"log"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
)

// dsn for the database. In order to run the tests locally, run the following command:
//
//	 ENT_INTEGRATION_ENDPOINT="root:pass@tcp(localhost:3306)/test?parseTime=True" go test -v
//
var dsn string

func ExampleUser() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the user's edges.
	u0 := client.User.
		Create().
		SetName("string").
		SetLastLogin(time.Now()).
		SaveX(ctx)
	log.Println("user created:", u0)
	ua1 := client.UserAccount.
		Create().
		SetName("string").
		SetPasswd("string").
		SetEmail("string").
		SetCreatedAt(time.Now()).
		SaveX(ctx)
	log.Println("useraccount created:", ua1)

	// create user vertex with its edges.
	u := client.User.
		Create().
		SetName("string").
		SetLastLogin(time.Now()).
		AddFriends(u0).
		SetAccount(ua1).
		SaveX(ctx)
	log.Println("user created:", u)

	// query edges.
	u0, err = u.QueryFriends().First(ctx)
	if err != nil {
		log.Fatalf("failed querying friends: %v", err)
	}
	log.Println("friends found:", u0)

	ua1, err = u.QueryAccount().First(ctx)
	if err != nil {
		log.Fatalf("failed querying account: %v", err)
	}
	log.Println("account found:", ua1)

	// Output:
}
func ExampleUserAccount() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the useraccount's edges.

	// create useraccount vertex with its edges.
	ua := client.UserAccount.
		Create().
		SetName("string").
		SetPasswd("string").
		SetEmail("string").
		SetCreatedAt(time.Now()).
		SaveX(ctx)
	log.Println("useraccount created:", ua)

	// query edges.

	// Output:
}