// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/oio-network/deeplx-extend/schema/ent/accesslog"
	"github.com/oio-network/deeplx-extend/schema/ent/user"
)

// AccessLog is the model entity for the AccessLog schema.
type AccessLog struct {
	config `json:"-"`
	// ID of the ent.
	ID int64 `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UserID holds the value of the "user_id" field.
	UserID int64 `json:"user_id,omitempty"`
	// IP holds the value of the "ip" field.
	IP string `json:"ip,omitempty"`
	// CountryName holds the value of the "country_name" field.
	CountryName string `json:"country_name,omitempty"`
	// CountryCode holds the value of the "country_code" field.
	CountryCode string `json:"country_code,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the AccessLogQuery when eager-loading is set.
	Edges        AccessLogEdges `json:"edges"`
	selectValues sql.SelectValues
}

// AccessLogEdges holds the relations/edges for other nodes in the graph.
type AccessLogEdges struct {
	// OwnerUser holds the value of the owner_user edge.
	OwnerUser *User `json:"owner_user,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// OwnerUserOrErr returns the OwnerUser value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e AccessLogEdges) OwnerUserOrErr() (*User, error) {
	if e.OwnerUser != nil {
		return e.OwnerUser, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: user.Label}
	}
	return nil, &NotLoadedError{edge: "owner_user"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*AccessLog) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case accesslog.FieldID, accesslog.FieldUserID:
			values[i] = new(sql.NullInt64)
		case accesslog.FieldIP, accesslog.FieldCountryName, accesslog.FieldCountryCode:
			values[i] = new(sql.NullString)
		case accesslog.FieldCreatedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the AccessLog fields.
func (al *AccessLog) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case accesslog.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			al.ID = int64(value.Int64)
		case accesslog.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				al.CreatedAt = value.Time
			}
		case accesslog.FieldUserID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field user_id", values[i])
			} else if value.Valid {
				al.UserID = value.Int64
			}
		case accesslog.FieldIP:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field ip", values[i])
			} else if value.Valid {
				al.IP = value.String
			}
		case accesslog.FieldCountryName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field country_name", values[i])
			} else if value.Valid {
				al.CountryName = value.String
			}
		case accesslog.FieldCountryCode:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field country_code", values[i])
			} else if value.Valid {
				al.CountryCode = value.String
			}
		default:
			al.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the AccessLog.
// This includes values selected through modifiers, order, etc.
func (al *AccessLog) Value(name string) (ent.Value, error) {
	return al.selectValues.Get(name)
}

// QueryOwnerUser queries the "owner_user" edge of the AccessLog entity.
func (al *AccessLog) QueryOwnerUser() *UserQuery {
	return NewAccessLogClient(al.config).QueryOwnerUser(al)
}

// Update returns a builder for updating this AccessLog.
// Note that you need to call AccessLog.Unwrap() before calling this method if this AccessLog
// was returned from a transaction, and the transaction was committed or rolled back.
func (al *AccessLog) Update() *AccessLogUpdateOne {
	return NewAccessLogClient(al.config).UpdateOne(al)
}

// Unwrap unwraps the AccessLog entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (al *AccessLog) Unwrap() *AccessLog {
	_tx, ok := al.config.driver.(*txDriver)
	if !ok {
		panic("ent: AccessLog is not a transactional entity")
	}
	al.config.driver = _tx.drv
	return al
}

// String implements the fmt.Stringer.
func (al *AccessLog) String() string {
	var builder strings.Builder
	builder.WriteString("AccessLog(")
	builder.WriteString(fmt.Sprintf("id=%v, ", al.ID))
	builder.WriteString("created_at=")
	builder.WriteString(al.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("user_id=")
	builder.WriteString(fmt.Sprintf("%v", al.UserID))
	builder.WriteString(", ")
	builder.WriteString("ip=")
	builder.WriteString(al.IP)
	builder.WriteString(", ")
	builder.WriteString("country_name=")
	builder.WriteString(al.CountryName)
	builder.WriteString(", ")
	builder.WriteString("country_code=")
	builder.WriteString(al.CountryCode)
	builder.WriteByte(')')
	return builder.String()
}

// AccessLogs is a parsable slice of AccessLog.
type AccessLogs []*AccessLog
