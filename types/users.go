package types

import (
	"encoding/json"
	"reflect"

	"github.com/xdrm-io/aicra/validator"
	"github.com/xdrm-io/articles-api/model"
)

// UsersType used for input validation
type UsersType struct{}

// GoType returns the `[]model.User` type
func (UsersType) GoType() reflect.Type {
	return reflect.TypeOf([]model.User{})
}

// Validator for users values
func (UsersType) Validator(typename string, avail ...validator.Type) validator.ValidateFunc {
	if typename != "[]user" {
		return nil
	}

	return func(value interface{}) (interface{}, bool) {
		switch cast := value.(type) {

		case []*model.User:
			values := make([]model.User, 0)
			for _, ptr := range cast {
				values = append(values, *ptr)
			}
			return values, true

		case []model.User:
			return cast, true

		case string:
		case []byte:
			asBytes := []byte(cast)
			var receiver []model.User
			err := json.Unmarshal(asBytes, &receiver)
			return receiver, err == nil

		}

		// unknown type
		return []model.User{}, false
	}
}
