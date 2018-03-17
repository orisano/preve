package preve

import (
	"reflect"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

var pullRequestActions = map[string]bool{
	"assigned":                 true,
	"unassigned":               true,
	"review_requested":         true,
	"review_requested_removed": true,
	"labeled":                  true,
	"unlabeled":                true,
	"opened":                   true,
	"edited":                   true,
	"closed":                   true,
	"reopened":                 true,
}

func IsValidPullRequestAction(action string) bool {
	return pullRequestActions[action]
}

func validatePullRequestAction(fl validator.FieldLevel) bool {
	f := fl.Field()
	switch k := f.Kind(); k {
	case reflect.Slice:
		if k := f.Elem().Kind(); k != reflect.String {
			panic("unsupported kind: " + k.String())
		}
		for i := 0; i < f.Len(); i++ {
			action := f.Index(i).String()
			if !IsValidPullRequestAction(action) {
				return false
			}
		}
		return true
	case reflect.String:
		return IsValidPullRequestAction(f.String())
	default:
		panic("unsupported kind: " + k.String())
	}
}

func validateRepository(fl validator.FieldLevel) bool {
	tokens := strings.Split(fl.Field().String(), "/")
	if len(tokens) != 2 {
		return false
	}
	return len(tokens[0]) > 0 && len(tokens[1]) > 0
}

func NewValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("pr_action", validatePullRequestAction)
	validate.RegisterValidation("repository", validateRepository)
	return validate
}
