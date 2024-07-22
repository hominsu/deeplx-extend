// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: deeplx/v1/translate.proto

package v1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on TranslateRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *TranslateRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on TranslateRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// TranslateRequestMultiError, or nil if none found.
func (m *TranslateRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *TranslateRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Token

	if all {
		switch v := interface{}(m.GetPayload()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, TranslateRequestValidationError{
					field:  "Payload",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, TranslateRequestValidationError{
					field:  "Payload",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetPayload()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return TranslateRequestValidationError{
				field:  "Payload",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return TranslateRequestMultiError(errors)
	}

	return nil
}

// TranslateRequestMultiError is an error wrapping multiple validation errors
// returned by TranslateRequest.ValidateAll() if the designated constraints
// aren't met.
type TranslateRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TranslateRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TranslateRequestMultiError) AllErrors() []error { return m }

// TranslateRequestValidationError is the validation error returned by
// TranslateRequest.Validate if the designated constraints aren't met.
type TranslateRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TranslateRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TranslateRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TranslateRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TranslateRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TranslateRequestValidationError) ErrorName() string { return "TranslateRequestValidationError" }

// Error satisfies the builtin error interface
func (e TranslateRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTranslateRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TranslateRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TranslateRequestValidationError{}

// Validate checks the field values on TranslateRequest_Payload with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *TranslateRequest_Payload) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on TranslateRequest_Payload with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// TranslateRequest_PayloadMultiError, or nil if none found.
func (m *TranslateRequest_Payload) ValidateAll() error {
	return m.validate(true)
}

func (m *TranslateRequest_Payload) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Text

	// no validation rules for SourceLang

	// no validation rules for TargetLang

	if len(errors) > 0 {
		return TranslateRequest_PayloadMultiError(errors)
	}

	return nil
}

// TranslateRequest_PayloadMultiError is an error wrapping multiple validation
// errors returned by TranslateRequest_Payload.ValidateAll() if the designated
// constraints aren't met.
type TranslateRequest_PayloadMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TranslateRequest_PayloadMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TranslateRequest_PayloadMultiError) AllErrors() []error { return m }

// TranslateRequest_PayloadValidationError is the validation error returned by
// TranslateRequest_Payload.Validate if the designated constraints aren't met.
type TranslateRequest_PayloadValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TranslateRequest_PayloadValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TranslateRequest_PayloadValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TranslateRequest_PayloadValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TranslateRequest_PayloadValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TranslateRequest_PayloadValidationError) ErrorName() string {
	return "TranslateRequest_PayloadValidationError"
}

// Error satisfies the builtin error interface
func (e TranslateRequest_PayloadValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTranslateRequest_Payload.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TranslateRequest_PayloadValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TranslateRequest_PayloadValidationError{}
