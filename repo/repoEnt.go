package repo

import (
	"context"
	"github.com/pkg/errors"
	"log"
	"mq/academy/ent"
	"mq/academy/ent/user"
	"time"
)

var userNil = errors.New("user pointer is nil")

type GameRepoEnt struct {
	client *ent.Client
}

func NewRepo() GameRepo {
	client, err := ent.Open("sqlite3", "file:ent.db?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}
	//defer client.Close()

	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	if err := Gen(ctx, client); err != nil {
		log.Fatal(err)
	}
	//if err := Traverse(ctx, client); err != nil {
	//	log.Fatal(err)
	//}
	return &GameRepoEnt{client: client}
}

func (repo *GameRepoEnt) WithTx(ctx context.Context, fn func(repo GameRepo) error) error {
	tx, err := repo.client.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()

	if err := fn(&GameRepoEnt{client: tx.Client()}); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = errors.Wrapf(err, "rollback got error: %v", rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commit got error")
	}
	return nil
}

func (repo *GameRepoEnt) CreateUser(ctx context.Context, name string) (*ent.User, error) {
	return repo.client.User.Create().SetName(name).SetLastLogin(time.Now()).Save(ctx)
}

func (repo *GameRepoEnt) CreateUserAccount(ctx context.Context, u *ent.User, passwd, email string) (*ent.UserAccount, error) {
	if u == nil {
		return nil, userNil
	}
	return repo.client.UserAccount.Create().
		SetName(u.Name).SetEmail(email).SetPasswd(passwd).
		SetCreatedAt(time.Now()).
		SetOwner(u).
		Save(ctx)
}

func (repo *GameRepoEnt) DeleteUser(ctx context.Context, u *ent.User) error {
	if u == nil {
		return userNil
	}
	return repo.client.User.DeleteOne(u).Exec(ctx)
}

func (repo *GameRepoEnt) SelectUserByName(ctx context.Context, name string) (*ent.User, error) {

	return repo.client.User.Query().Where(user.Name(name)).Only(ctx)
}

func (repo *GameRepoEnt) SelectUsersByName(ctx context.Context, names ...string) ([]*ent.User, error) {
	return repo.client.User.Query().Where(user.NameIn(names...)).All(ctx)
}

func (repo *GameRepoEnt) UserAccount(ctx context.Context, u *ent.User) (*ent.UserAccount, error) {
	if u == nil {
		return nil, userNil
	}
	return u.QueryAccount().Only(ctx)
}

func (repo *GameRepoEnt) UserFriends(ctx context.Context, u *ent.User) ([]*ent.User, error) {
	if u == nil {
		return nil, userNil
	}
	us, err := u.QueryFriends().All(ctx)
	if err != nil {
		return nil, err
	}
	ret := make([]*ent.User, len(us))
	for i, u := range us {
		ret[i] = u
	}
	return ret, nil
}

func (repo *GameRepoEnt) AddFriends(ctx context.Context, u *ent.User, fs []*ent.User) error {
	if u == nil {
		return userNil
	}
	return u.Update().AddFriends(fs...).Exec(ctx)
}

func (repo *GameRepoEnt) AddFriendsByName(ctx context.Context, u *ent.User, names ...string) error {
	if u == nil {
		return userNil
	}
	fs, err := repo.SelectUsersByName(ctx, names...)
	if err != nil {
		return err
	}
	return u.Update().AddFriends(fs...).Exec(ctx)
}

func (repo *GameRepoEnt) UserMaybeFriends(ctx context.Context, u *ent.User) ([]*ent.User, error) {
	if u == nil {
		return nil, userNil
	}
	us, err := u.QueryFriends().QueryFriends().
		//Where(user.Not(user.Name(u.Name))).
		//Where(user.Not(user.HasFriendsWith(user.Name(u.Name)))).
		Where(user.And(
			user.Not(user.Name(u.Name)),
			user.Not(user.HasFriendsWith(user.Name(u.Name))),
		)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	ret := make([]*ent.User, len(us))
	for i, u := range us {
		ret[i] = u
	}
	return ret, nil
}

func Gen(ctx context.Context, client *ent.Client) error {
	now := time.Now()
	coco, err := client.User.Create().SetName("Coco").SetLastLogin(now).Save(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}
	dan, err := client.User.Create().SetName("Dan").SetLastLogin(now).AddFriends(coco).Save(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}
	alex, err := client.User.Create().SetName("Alex").SetLastLogin(now).AddFriends(coco).Save(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}
	momo, err := client.User.Create().SetName("Momo").SetLastLogin(now).AddFriends(coco).AddFriends(alex).Save(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}
	dodo, err := client.User.Create().SetName("Dodo").SetLastLogin(now).AddFriends(alex).Save(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}
	log.Print("user created: ", coco, dan, alex, momo, dodo)
	return nil
}

func Traverse(ctx context.Context, client *ent.Client) error {
	friendOfFriend, err := client.User.Query().
		Where(user.Name("Alex")).QueryFriends().QueryFriends().
		Where(user.Not(user.Name("Alex"))).
		//Where(user.Not(user.HasFriendsWith(user.Name("Alex")))).
		All(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to query")
	}

	friend, err := client.User.Query().
		Where(user.HasFriendsWith(user.Name("Alex"))).
		All(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to query")
	}
	temp := make(map[string]*ent.User)
	for _, ff := range friendOfFriend {
		temp[ff.Name] = ff
	}
	for _, f := range friend {
		delete(temp, f.Name)
	}
	maybeFriend := make([]*ent.User, 0, len(temp))
	for _, mf := range temp {
		maybeFriend = append(maybeFriend, mf)
	}
	return nil
}

func Traverse2(ctx context.Context, client *ent.Client) error {
	_, err := client.User.Query().
		Where(user.Name("Alex")).QueryFriends().QueryFriends().
		Where(user.And(
			user.Not(user.Name("Alex")),
			user.Not(user.HasFriendsWith(user.Name("Alex"))),
		)).
		All(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to query")
	}
	return nil
}
