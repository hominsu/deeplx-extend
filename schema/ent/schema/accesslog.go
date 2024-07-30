package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// AccessLog holds the schema definition for the AccessLog entity.
type AccessLog struct {
	ent.Schema
}

// Annotations of the AccessLog.
func (AccessLog) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "access_logs",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
	}
}

// Fields of the AccessLog.
func (AccessLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.Int64("user_id").Optional(),
		field.String("ip"),
		field.String("country_name").Optional(),
		field.String("country_code").Optional(),
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
