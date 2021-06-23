package types

import (
	"encoding/json"
	"reflect"

	"github.com/xdrm-io/aicra/validator"
	"github.com/xdrm-io/articles-api/model"
)

// ArticlesType used for input validation
type ArticlesType struct{}

// GoType returns the `[]model.Article` type
func (ArticlesType) GoType() reflect.Type {
	return reflect.TypeOf([]model.Article{})
}

// Validator for articles values
func (ArticlesType) Validator(typename string, avail ...validator.Type) validator.ValidateFunc {
	if typename != "[]article" {
		return nil
	}

	return func(value interface{}) (interface{}, bool) {
		switch cast := value.(type) {

		case []*model.Article:
			values := make([]model.Article, 0)
			for _, ptr := range cast {
				values = append(values, *ptr)
			}
			return values, true

		case []model.Article:
			return cast, true

		case string:
		case []byte:
			asBytes := []byte(cast)
			var receiver []model.Article
			err := json.Unmarshal(asBytes, &receiver)
			return receiver, err == nil

		}

		// unknown type
		return []model.Article{}, false
	}
}
