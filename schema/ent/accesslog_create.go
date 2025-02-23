// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/oio-network/deeplx-extend/schema/ent/accesslog"
	"github.com/oio-network/deeplx-extend/schema/ent/user"
)

// AccessLogCreate is the builder for creating a AccessLog entity.
type AccessLogCreate struct {
	config
	mutation *AccessLogMutation
	hooks    []Hook
}

// SetCreatedAt sets the "created_at" field.
func (alc *AccessLogCreate) SetCreatedAt(t time.Time) *AccessLogCreate {
	alc.mutation.SetCreatedAt(t)
	return alc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (alc *AccessLogCreate) SetNillableCreatedAt(t *time.Time) *AccessLogCreate {
	if t != nil {
		alc.SetCreatedAt(*t)
	}
	return alc
}

// SetUserID sets the "user_id" field.
func (alc *AccessLogCreate) SetUserID(i int64) *AccessLogCreate {
	alc.mutation.SetUserID(i)
	return alc
}

// SetNillableUserID sets the "user_id" field if the given value is not nil.
func (alc *AccessLogCreate) SetNillableUserID(i *int64) *AccessLogCreate {
	if i != nil {
		alc.SetUserID(*i)
	}
	return alc
}

// SetIP sets the "ip" field.
func (alc *AccessLogCreate) SetIP(s string) *AccessLogCreate {
	alc.mutation.SetIP(s)
	return alc
}

// SetCountryName sets the "country_name" field.
func (alc *AccessLogCreate) SetCountryName(s string) *AccessLogCreate {
	alc.mutation.SetCountryName(s)
	return alc
}

// SetNillableCountryName sets the "country_name" field if the given value is not nil.
func (alc *AccessLogCreate) SetNillableCountryName(s *string) *AccessLogCreate {
	if s != nil {
		alc.SetCountryName(*s)
	}
	return alc
}

// SetCountryCode sets the "country_code" field.
func (alc *AccessLogCreate) SetCountryCode(s string) *AccessLogCreate {
	alc.mutation.SetCountryCode(s)
	return alc
}

// SetNillableCountryCode sets the "country_code" field if the given value is not nil.
func (alc *AccessLogCreate) SetNillableCountryCode(s *string) *AccessLogCreate {
	if s != nil {
		alc.SetCountryCode(*s)
	}
	return alc
}

// SetID sets the "id" field.
func (alc *AccessLogCreate) SetID(i int64) *AccessLogCreate {
	alc.mutation.SetID(i)
	return alc
}

// SetOwnerUserID sets the "owner_user" edge to the User entity by ID.
func (alc *AccessLogCreate) SetOwnerUserID(id int64) *AccessLogCreate {
	alc.mutation.SetOwnerUserID(id)
	return alc
}

// SetNillableOwnerUserID sets the "owner_user" edge to the User entity by ID if the given value is not nil.
func (alc *AccessLogCreate) SetNillableOwnerUserID(id *int64) *AccessLogCreate {
	if id != nil {
		alc = alc.SetOwnerUserID(*id)
	}
	return alc
}

// SetOwnerUser sets the "owner_user" edge to the User entity.
func (alc *AccessLogCreate) SetOwnerUser(u *User) *AccessLogCreate {
	return alc.SetOwnerUserID(u.ID)
}

// Mutation returns the AccessLogMutation object of the builder.
func (alc *AccessLogCreate) Mutation() *AccessLogMutation {
	return alc.mutation
}

// Save creates the AccessLog in the database.
func (alc *AccessLogCreate) Save(ctx context.Context) (*AccessLog, error) {
	alc.defaults()
	return withHooks(ctx, alc.sqlSave, alc.mutation, alc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (alc *AccessLogCreate) SaveX(ctx context.Context) *AccessLog {
	v, err := alc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (alc *AccessLogCreate) Exec(ctx context.Context) error {
	_, err := alc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (alc *AccessLogCreate) ExecX(ctx context.Context) {
	if err := alc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (alc *AccessLogCreate) defaults() {
	if _, ok := alc.mutation.CreatedAt(); !ok {
		v := accesslog.DefaultCreatedAt()
		alc.mutation.SetCreatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (alc *AccessLogCreate) check() error {
	if _, ok := alc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "AccessLog.created_at"`)}
	}
	if _, ok := alc.mutation.IP(); !ok {
		return &ValidationError{Name: "ip", err: errors.New(`ent: missing required field "AccessLog.ip"`)}
	}
	return nil
}

func (alc *AccessLogCreate) sqlSave(ctx context.Context) (*AccessLog, error) {
	if err := alc.check(); err != nil {
		return nil, err
	}
	_node, _spec := alc.createSpec()
	if err := sqlgraph.CreateNode(ctx, alc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != _node.ID {
		id := _spec.ID.Value.(int64)
		_node.ID = int64(id)
	}
	alc.mutation.id = &_node.ID
	alc.mutation.done = true
	return _node, nil
}

func (alc *AccessLogCreate) createSpec() (*AccessLog, *sqlgraph.CreateSpec) {
	var (
		_node = &AccessLog{config: alc.config}
		_spec = sqlgraph.NewCreateSpec(accesslog.Table, sqlgraph.NewFieldSpec(accesslog.FieldID, field.TypeInt64))
	)
	if id, ok := alc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := alc.mutation.CreatedAt(); ok {
		_spec.SetField(accesslog.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := alc.mutation.IP(); ok {
		_spec.SetField(accesslog.FieldIP, field.TypeString, value)
		_node.IP = value
	}
	if value, ok := alc.mutation.CountryName(); ok {
		_spec.SetField(accesslog.FieldCountryName, field.TypeString, value)
		_node.CountryName = value
	}
	if value, ok := alc.mutation.CountryCode(); ok {
		_spec.SetField(accesslog.FieldCountryCode, field.TypeString, value)
		_node.CountryCode = value
	}
	if nodes := alc.mutation.OwnerUserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   accesslog.OwnerUserTable,
			Columns: []string{accesslog.OwnerUserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.UserID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// AccessLogCreateBulk is the builder for creating many AccessLog entities in bulk.
type AccessLogCreateBulk struct {
	config
	err      error
	builders []*AccessLogCreate
}

// Save creates the AccessLog entities in the database.
func (alcb *AccessLogCreateBulk) Save(ctx context.Context) ([]*AccessLog, error) {
	if alcb.err != nil {
		return nil, alcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(alcb.builders))
	nodes := make([]*AccessLog, len(alcb.builders))
	mutators := make([]Mutator, len(alcb.builders))
	for i := range alcb.builders {
		func(i int, root context.Context) {
			builder := alcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*AccessLogMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, alcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, alcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil && nodes[i].ID == 0 {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int64(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, alcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (alcb *AccessLogCreateBulk) SaveX(ctx context.Context) []*AccessLog {
	v, err := alcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (alcb *AccessLogCreateBulk) Exec(ctx context.Context) error {
	_, err := alcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (alcb *AccessLogCreateBulk) ExecX(ctx context.Context) {
	if err := alcb.Exec(ctx); err != nil {
		panic(err)
	}
}
