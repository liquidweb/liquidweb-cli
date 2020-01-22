/*
Copyright © LiquidWeb

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package validate

import (
	"fmt"
	"reflect"

	"github.com/spf13/cast"
)

func Validate(chk map[interface{}]string) error {

	for inputFieldValue, inputField := range chk {
		// inputField must be defined
		defined, shouldBeType, fieldVal := inputTypeDefined(inputField)
		if !defined {
			return fmt.Errorf("%w for input field [%+v] type [%s] is not valid", ValidationFailure, inputFieldValue, inputField)
		}

		// inputFieldValue must be of the correct type
		reflectValue := reflect.TypeOf(inputFieldValue).Name()
		if reflectValue != shouldBeType {
			return fmt.Errorf("%w for input field [%+v] type [%s] has an invalid type of [%s] wanted [%s]",
				ValidationFailure, inputFieldValue, inputField, reflectValue, shouldBeType)
		}

		// if there's a Validate method call it
		iface := fieldVal.Interface()
		if interfaceHasMethod(iface, "Validate") {
			if err := interfaceInputTypeValidate(iface, inputFieldValue); err != nil {
				return fmt.Errorf("%w for input field [%+v] %s", ValidationFailure, inputFieldValue, err)
			}
		}
	}

	return nil
}

func interfaceInputTypeValidate(iface, inputFieldValue interface{}) error {
	switch iface.(type) {
	case InputTypeUniqId:
		var obj InputTypeUniqId
		obj.UniqId = cast.ToString(inputFieldValue)
		if err := obj.Validate(); err != nil {
			return err
		}
	case InputTypeIP:
		var obj InputTypeIP
		obj.IP = cast.ToString(inputFieldValue)
		if err := obj.Validate(); err != nil {
			return err
		}
	case InputTypePositiveInt64:
		var obj InputTypePositiveInt64
		obj.PositiveInt64 = cast.ToInt64(inputFieldValue)
		if err := obj.Validate(); err != nil {
			return err
		}
	case InputTypePositiveInt:
		var obj InputTypePositiveInt
		obj.PositiveInt = cast.ToInt(inputFieldValue)
		if err := obj.Validate(); err != nil {
			return err
		}
	case InputTypeNonEmptyString:
		var obj InputTypeNonEmptyString
		obj.NonEmptyString = cast.ToString(inputFieldValue)
		if err := obj.Validate(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("bug: validation missing entry for %s", inputFieldValue)
	}

	return nil
}

func interfaceHasMethod(iface interface{}, methodName string) bool {
	ifaceVal := reflect.ValueOf(iface)

	if !ifaceVal.IsValid() {
		// not valid, so we already know its false
		return false
	}

	if ifaceVal.Type().Kind() != reflect.Ptr {
		ifaceVal = reflect.New(reflect.TypeOf(iface))
	}

	method := ifaceVal.MethodByName(methodName)

	if method.IsValid() {
		return true
	}

	return false
}

func inputTypeDefined(inputType string) (bool, string, reflect.Value) {
	var validTypes InputTypes

	err, fieldType, fieldVal := structHasField(validTypes, inputType)
	if err != nil {
		return false, fieldType, fieldVal
	}

	return true, fieldType, fieldVal
}

func structHasField(data interface{}, fieldName string) (error, string, reflect.Value) {
	dataVal := reflect.ValueOf(data)

	if !dataVal.IsValid() {
		return fmt.Errorf("failed fetching value for fieldName [%s]", fieldName), "",
			reflect.Value{}
	}

	if dataVal.Type().Kind() != reflect.Ptr {
		dataVal = reflect.New(reflect.TypeOf(data))
	}

	fieldVal := dataVal.Elem().FieldByName(fieldName)
	if !fieldVal.IsValid() {
		return fmt.Errorf("[%s] has no field [%s]", dataVal.Type(), fieldName), "", fieldVal
	}

	fieldValKindStr := fieldVal.Kind().String()

	if fieldValKindStr == "struct" {
		fieldValKindStr = fieldVal.Field(0).Kind().String()
	}

	return nil, fieldValKindStr, fieldVal
}
