package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

var NumberCheck = regexp.MustCompile("[0-9]")
var NumberValid = regexp.MustCompile(`(^\+[0-9]{2}|^\+[0-9]{2}\(0\)|^\(\+[0-9]{2}\)\(0\)|^00[0-9]{2}|^0)([0-9]{9}$|[0-9\-\s]{10}$)`)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This Password is too small (minimum is %d characters)", d))
	}
}

func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) != d {
		f.Errors.Add(field, fmt.Sprintf("Check value length", d))
	}
}

func (f *Form) MobileNumberCheck(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		f.Errors.Add(field, "This number is invalid")
	}
}

func (f *Form) MobileCountryCheckCode(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		f.Errors.Add(field, "Invalid phone number")
	}
}
func (f *Form) MobileCheckPref(field string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !strings.HasPrefix(value, "+") {
		f.Errors.Add(field, "Please check number. Must start with +")
	}
}
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
