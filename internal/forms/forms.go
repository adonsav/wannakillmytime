package forms

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/url"
	"strings"
)

// Form creates a custom form struct, embeds a url.values object
type Form struct {
	url.Values
	Errors errors
}

var validate = validator.New()

// New initializes the form struct
func New(data url.Values) *Form {
	return &Form{
		Values: data,
		Errors: errors(map[string][]string{}),
	}
}

// Required checks for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be empty")
		}
	}
}

// Has checks if form field is present and not empty
func (f *Form) Has(field string) bool {
	x := f.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field cannot be empty")
		return false
	}
	return true
}

// Valid returns true if there are no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// MinLength checks for field input minimum length
func (f *Form) MinLength(field string, length int) bool {
	fld := f.Get(field)
	if len(fld) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

// IsEmail checks for valid email address using https://github.com/go-playground/validator
func (f *Form) IsEmail(field string) {
	errs := validate.Var(f.Get(field), "email")
	if errs != nil {
		f.Errors.Add(field, "Invalid email address")
	}
}
