package account

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/Optum/dce/pkg/arn"
	validation "github.com/go-ozzo/ozzo-validation"
)

// We don't use the internal errors package here because validation will rewrite it anyways
// Just spit out errors and turn them into validation errors inside the appropriate functions

var validateAdminRoleArn = []validation.Rule{
	validation.NotNil.Error("must be a string"),
}

var validateID = []validation.Rule{
	validation.NotNil.Error("must be a string"),
	validation.Match(regexp.MustCompile("^[0-9]{12}$")).Error("must be a string with 12 digits"),
}

var validateInt64 = []validation.Rule{
	validation.NotNil.Error("must be an epoch timestamp"),
}

var validatePrincipalRoleArn = []validation.Rule{
	validation.NilOrNotEmpty.Error("must be an ARN or empty"),
}

var validatePrincipalPolicyHash = []validation.Rule{
	validation.NilOrNotEmpty.Error("must be a hash or empty"),
}

var validateStatus = []validation.Rule{
	validation.NotNil.Error("must be a valid account status"),
}

var validateMetadata = []validation.Rule{
	validation.By(validateMetadataContent),
}

func isNil(value interface{}) error {
	if !reflect.ValueOf(value).IsNil() {
		return errors.New("must be empty")
	}
	return nil
}

func isNilOrUsableAdminRole(am Manager) validation.RuleFunc {
	return func(value interface{}) error {
		if !reflect.ValueOf(value).IsNil() {
			a, _ := value.(*arn.ARN)
			err := am.ValidateAccess(a)
			if err != nil {
				return errors.New("must be an admin role arn that can be assumed")
			}
		}
		return nil
	}
}

func isAccountNotLeased(value interface{}) error {
	s, _ := value.(*Status)
	if s.String() == StatusLeased.String() {
		return errors.New("must not be leased")
	}
	return nil
}

func validateMetadataContent(value interface{}) error {
	if value == nil {
		return nil
	}
	
	metadata, ok := value.(map[string]interface{})
	if !ok {
		return errors.New("metadata must be a map[string]interface{}")
	}
	
	const maxKeys = 50
	const maxKeyLength = 128
	const maxValueLength = 1024
	
	if len(metadata) > maxKeys {
		return fmt.Errorf("metadata cannot have more than %d keys", maxKeys)
	}
	
	for key, val := range metadata {
		if len(key) > maxKeyLength {
			return fmt.Errorf("metadata key '%s' exceeds maximum length of %d characters", key, maxKeyLength)
		}
		
		if strings.Contains(key, "<") || strings.Contains(key, ">") || strings.Contains(key, "&") {
			return fmt.Errorf("metadata key '%s' contains invalid characters", key)
		}
		
		if strings.HasPrefix(key, "_") || strings.HasPrefix(key, "$") {
			return fmt.Errorf("metadata key '%s' cannot start with reserved characters (_ or $)", key)
		}
		
		switch v := val.(type) {
		case string:
			if len(v) > maxValueLength {
				return fmt.Errorf("metadata value for key '%s' exceeds maximum length of %d characters", key, maxValueLength)
			}
			if strings.Contains(v, "<script") || strings.Contains(v, "javascript:") || strings.Contains(v, "data:") {
				return fmt.Errorf("metadata value for key '%s' contains potentially dangerous content", key)
			}
		case bool, int, int64, float64:
		case nil:
		default:
			return fmt.Errorf("metadata value for key '%s' must be string, number, boolean, or null", key)
		}
	}
	
	return nil
}
