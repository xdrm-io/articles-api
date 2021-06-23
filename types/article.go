package types

import (
	"encoding/json"
	"reflect"

	"github.com/xdrm-io/aicra/validator"
	"github.com/xdrm-io/articles-api/model"
)

// ArticleType used for input validation
type ArticleType struct{}

// GoType returns the `model.Article` type
func (ArticleType) GoType() reflect.Type {
	return reflect.TypeOf(model.Article{})
}

// Validator for article values
func (ArticleType) Validator(typename string, avail ...validator.Type) validator.ValidateFunc {
	if typename != "article" {
		return nil
	}

	return func(value interface{}) (interface{}, bool) {
		switch cast := value.(type) {

		case *model.Article:
			return *cast, true

		case model.Article:
			return cast, true

		case string:
		case []byte:
			asBytes := []byte(cast)
			var receiver model.Article
			err := json.Unmarshal(asBytes, &receiver)
			return receiver, err == nil

		}

		// unknown type
		return model.Article{}, false
	}
}
