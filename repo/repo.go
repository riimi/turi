package repo

import (
	"context"
	"mq/academy/ent"
)

type GameRepo interface {
	WithTx(ctx context.Context, fn func(gr GameRepo) error) error
	CreateUser(ctx context.Context, name string) (*ent.User, error)
	CreateUserAccount(ctx context.Context, u *ent.User, passwd, email string) (*ent.UserAccount, error)
	DeleteUser(ctx context.Context, u *ent.User) error
	SelectUserByName(ctx context.Context, name string) (*ent.User, error)
	SelectUsersByName(ctx context.Context, names ...string) ([]*ent.User, error)
	UserAccount(ctx context.Context, u *ent.User) (*ent.UserAccount, error)
	UserFriends(ctx context.Context, u *ent.User) ([]*ent.User, error)
	AddFriends(ctx context.Context, u *ent.User, fs []*ent.User) error
	AddFriendsByName(ctx context.Context, u *ent.User, names ...string) error
	UserMaybeFriends(ctx context.Context, u *ent.User) ([]*ent.User, error)
}

type contextKey struct{}

func NewContext(parent context.Context, repo GameRepo) context.Context {
	return context.WithValue(parent, contextKey{}, repo)
}

func FromContext(ctx context.Context) GameRepo {
	return ctx.Value(contextKey{}).(GameRepo)
}
