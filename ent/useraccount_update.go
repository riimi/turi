// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"mq/academy/ent/predicate"
	"mq/academy/ent/user"
	"mq/academy/ent/useraccount"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
)

// UserAccountUpdate is the builder for updating UserAccount entities.
type UserAccountUpdate struct {
	config
	name         *string
	passwd       *string
	email        *string
	createdAt    *time.Time
	owner        map[int]struct{}
	clearedOwner bool
	predicates   []predicate.UserAccount
}

// Where adds a new predicate for the builder.
func (uau *UserAccountUpdate) Where(ps ...predicate.UserAccount) *UserAccountUpdate {
	uau.predicates = append(uau.predicates, ps...)
	return uau
}

// SetName sets the name field.
func (uau *UserAccountUpdate) SetName(s string) *UserAccountUpdate {
	uau.name = &s
	return uau
}

// SetPasswd sets the passwd field.
func (uau *UserAccountUpdate) SetPasswd(s string) *UserAccountUpdate {
	uau.passwd = &s
	return uau
}

// SetEmail sets the email field.
func (uau *UserAccountUpdate) SetEmail(s string) *UserAccountUpdate {
	uau.email = &s
	return uau
}

// SetCreatedAt sets the createdAt field.
func (uau *UserAccountUpdate) SetCreatedAt(t time.Time) *UserAccountUpdate {
	uau.createdAt = &t
	return uau
}

// SetOwnerID sets the owner edge to User by id.
func (uau *UserAccountUpdate) SetOwnerID(id int) *UserAccountUpdate {
	if uau.owner == nil {
		uau.owner = make(map[int]struct{})
	}
	uau.owner[id] = struct{}{}
	return uau
}

// SetOwner sets the owner edge to User.
func (uau *UserAccountUpdate) SetOwner(u *User) *UserAccountUpdate {
	return uau.SetOwnerID(u.ID)
}

// ClearOwner clears the owner edge to User.
func (uau *UserAccountUpdate) ClearOwner() *UserAccountUpdate {
	uau.clearedOwner = true
	return uau
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (uau *UserAccountUpdate) Save(ctx context.Context) (int, error) {
	if len(uau.owner) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"owner\"")
	}
	if uau.clearedOwner && uau.owner == nil {
		return 0, errors.New("ent: clearing a unique edge \"owner\"")
	}
	return uau.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (uau *UserAccountUpdate) SaveX(ctx context.Context) int {
	affected, err := uau.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (uau *UserAccountUpdate) Exec(ctx context.Context) error {
	_, err := uau.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uau *UserAccountUpdate) ExecX(ctx context.Context) {
	if err := uau.Exec(ctx); err != nil {
		panic(err)
	}
}

func (uau *UserAccountUpdate) sqlSave(ctx context.Context) (n int, err error) {
	selector := sql.Select(useraccount.FieldID).From(sql.Table(useraccount.Table))
	for _, p := range uau.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = uau.driver.Query(ctx, query, args, rows); err != nil {
		return 0, err
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return 0, fmt.Errorf("ent: failed reading id: %v", err)
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return 0, nil
	}

	tx, err := uau.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		builder = sql.Update(useraccount.Table).Where(sql.InInts(useraccount.FieldID, ids...))
	)
	if value := uau.name; value != nil {
		builder.Set(useraccount.FieldName, *value)
	}
	if value := uau.passwd; value != nil {
		builder.Set(useraccount.FieldPasswd, *value)
	}
	if value := uau.email; value != nil {
		builder.Set(useraccount.FieldEmail, *value)
	}
	if value := uau.createdAt; value != nil {
		builder.Set(useraccount.FieldCreatedAt, *value)
	}
	if !builder.Empty() {
		query, args := builder.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if uau.clearedOwner {
		query, args := sql.Update(useraccount.OwnerTable).
			SetNull(useraccount.OwnerColumn).
			Where(sql.InInts(user.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(uau.owner) > 0 {
		for _, id := range ids {
			eid := keys(uau.owner)[0]
			query, args := sql.Update(useraccount.OwnerTable).
				Set(useraccount.OwnerColumn, eid).
				Where(sql.EQ(useraccount.FieldID, id).And().IsNull(useraccount.OwnerColumn)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(uau.owner) {
				return 0, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"owner\" %v already connected to a different \"UserAccount\"", keys(uau.owner))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
}

// UserAccountUpdateOne is the builder for updating a single UserAccount entity.
type UserAccountUpdateOne struct {
	config
	id           int
	name         *string
	passwd       *string
	email        *string
	createdAt    *time.Time
	owner        map[int]struct{}
	clearedOwner bool
}

// SetName sets the name field.
func (uauo *UserAccountUpdateOne) SetName(s string) *UserAccountUpdateOne {
	uauo.name = &s
	return uauo
}

// SetPasswd sets the passwd field.
func (uauo *UserAccountUpdateOne) SetPasswd(s string) *UserAccountUpdateOne {
	uauo.passwd = &s
	return uauo
}

// SetEmail sets the email field.
func (uauo *UserAccountUpdateOne) SetEmail(s string) *UserAccountUpdateOne {
	uauo.email = &s
	return uauo
}

// SetCreatedAt sets the createdAt field.
func (uauo *UserAccountUpdateOne) SetCreatedAt(t time.Time) *UserAccountUpdateOne {
	uauo.createdAt = &t
	return uauo
}

// SetOwnerID sets the owner edge to User by id.
func (uauo *UserAccountUpdateOne) SetOwnerID(id int) *UserAccountUpdateOne {
	if uauo.owner == nil {
		uauo.owner = make(map[int]struct{})
	}
	uauo.owner[id] = struct{}{}
	return uauo
}

// SetOwner sets the owner edge to User.
func (uauo *UserAccountUpdateOne) SetOwner(u *User) *UserAccountUpdateOne {
	return uauo.SetOwnerID(u.ID)
}

// ClearOwner clears the owner edge to User.
func (uauo *UserAccountUpdateOne) ClearOwner() *UserAccountUpdateOne {
	uauo.clearedOwner = true
	return uauo
}

// Save executes the query and returns the updated entity.
func (uauo *UserAccountUpdateOne) Save(ctx context.Context) (*UserAccount, error) {
	if len(uauo.owner) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"owner\"")
	}
	if uauo.clearedOwner && uauo.owner == nil {
		return nil, errors.New("ent: clearing a unique edge \"owner\"")
	}
	return uauo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (uauo *UserAccountUpdateOne) SaveX(ctx context.Context) *UserAccount {
	ua, err := uauo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return ua
}

// Exec executes the query on the entity.
func (uauo *UserAccountUpdateOne) Exec(ctx context.Context) error {
	_, err := uauo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uauo *UserAccountUpdateOne) ExecX(ctx context.Context) {
	if err := uauo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (uauo *UserAccountUpdateOne) sqlSave(ctx context.Context) (ua *UserAccount, err error) {
	selector := sql.Select(useraccount.Columns...).From(sql.Table(useraccount.Table))
	useraccount.ID(uauo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = uauo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		ua = &UserAccount{config: uauo.config}
		if err := ua.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into UserAccount: %v", err)
		}
		id = ua.ID
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("UserAccount with id: %v", uauo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one UserAccount with the same id: %v", uauo.id)
	}

	tx, err := uauo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		builder = sql.Update(useraccount.Table).Where(sql.InInts(useraccount.FieldID, ids...))
	)
	if value := uauo.name; value != nil {
		builder.Set(useraccount.FieldName, *value)
		ua.Name = *value
	}
	if value := uauo.passwd; value != nil {
		builder.Set(useraccount.FieldPasswd, *value)
		ua.Passwd = *value
	}
	if value := uauo.email; value != nil {
		builder.Set(useraccount.FieldEmail, *value)
		ua.Email = *value
	}
	if value := uauo.createdAt; value != nil {
		builder.Set(useraccount.FieldCreatedAt, *value)
		ua.CreatedAt = *value
	}
	if !builder.Empty() {
		query, args := builder.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if uauo.clearedOwner {
		query, args := sql.Update(useraccount.OwnerTable).
			SetNull(useraccount.OwnerColumn).
			Where(sql.InInts(user.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(uauo.owner) > 0 {
		for _, id := range ids {
			eid := keys(uauo.owner)[0]
			query, args := sql.Update(useraccount.OwnerTable).
				Set(useraccount.OwnerColumn, eid).
				Where(sql.EQ(useraccount.FieldID, id).And().IsNull(useraccount.OwnerColumn)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return nil, rollback(tx, err)
			}
			if int(affected) < len(uauo.owner) {
				return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"owner\" %v already connected to a different \"UserAccount\"", keys(uauo.owner))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return ua, nil
}