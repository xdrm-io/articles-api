package types

import (
	"encoding/json"
	"reflect"

	"github.com/xdrm-io/aicra/validator"
	"github.com/xdrm-io/articles-api/model"
)

// UserType used for input validation
type UserType struct{}

// GoType returns the `model.User` type
func (UserType) GoType() reflect.Type {
	return reflect.TypeOf(model.User{})
}

// Validator for user values
func (UserType) Validator(typename string, avail ...validator.Type) validator.ValidateFunc {
	if typename != "user" {
		return nil
	}

	return func(value interface{}) (interface{}, bool) {
		switch cast := value.(type) {

		case *model.User:
			return *cast, true

		case model.User:
			return cast, true

		case string:
		case []byte:
			asBytes := []byte(cast)
			var receiver model.User
			err := json.Unmarshal(asBytes, &receiver)
			return receiver, err == nil

		}

		// unknown type
		return model.User{}, false
	}
}
