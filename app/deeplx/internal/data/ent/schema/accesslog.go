package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// AccessLog holds the schema definition for the AccessLog entity.
type AccessLog struct {
	ent.Schema
}

// Fields of the AccessLog.
func (AccessLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.Int64("user_id").Optional(),
		field.String("ip"),
		field.String("country_name"),
		field.String("country_code"),
	}
}

// Mixin of the AccessLog.
func (AccessLog) Mixin() []ent.Mixin {
	return []ent.Mixin{
		CreateTimeMixin{},
	}
}

// Edges of the AccessLog.
func (AccessLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner_user", User.Type).
			Ref("access_logs").
			Unique().
			Field("user_id"),
	}
}
