package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// UserAccount holds the schema definition for the UserAccount entity.
type UserAccount struct {
	ent.Schema
}

// Fields of the UserAccount.
func (UserAccount) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("passwd"),
		field.String("email"),
		field.Time("createdAt"),
	}
}

// Edges of the UserAccount.
func (UserAccount) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).Ref("account").Unique().Required(),
	}
}
