package orm

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

var registry sync.Map

func ParseOrGetMeta(model Model) (*MapRelationMeta, error) {
	modelType := reflect.TypeOf(model)

	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	if res, ok := registry.Load(modelType); ok {
		return res.(*MapRelationMeta), nil
	}

	res := MapRelationMeta{}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		ormTag := field.Tag.Get("orm")
		if ormTag == "" {
			continue
		}

		fieldType := field.Type
		fieldTypeElem := fieldType.Elem()

		// only for pointer of struct or slice of struct
		if (fieldType.Kind() != reflect.Pointer && fieldType.Kind() != reflect.Slice) || (fieldTypeElem.Kind() != reflect.Struct) {
			continue
		}

		subRes := RelationMeta{}

		subRes.MainTable = model.TableName()
		subRes.MainType = modelType

		if iface, ok := reflect.New(fieldTypeElem).Interface().(Model); ok {
			subRes.AssocTable = iface.TableName()
			subRes.AssocType = fieldTypeElem
		}

		parts := strings.SplitSeq(ormTag, ";")

		for j := range parts {
			subParts := strings.Split(j, ":")

			if len(subParts) != 2 {
				continue
			}

			key, val := subParts[0], subParts[1]

			switch key {
			case "relation":
				subRes.Relation = val
			case "assocField":
				subRes.AssocField = val
			case "mainField":
				subRes.MainField = val
			}
		}

		{
			mainStructField, ok := subRes.MainType.FieldByName(subRes.MainField)
			if !ok {
				return nil, fmt.Errorf("orm: mainField '%s' specified in tag for %s.%s does not exist in model %s", subRes.MainField, modelType.Name(), field.Name, subRes.MainType.Name())
			}
			subRes.MainFieldIndex = mainStructField.Index

			assocStructField, ok := subRes.AssocType.FieldByName(subRes.AssocField)
			if !ok {
				return nil, fmt.Errorf("orm: assocField '%s' specified in tag for %s.%s does not exist in model %s", subRes.AssocField, modelType.Name(), field.Name, subRes.AssocType.Name())
			}
			subRes.AssocFieldIndex = assocStructField.Index

			subRes.GetMainField = func(m Model) any {
				v := reflect.ValueOf(m).Elem()
				f := v.FieldByIndex(subRes.MainFieldIndex)

				if f.IsZero() {
					return nil
				}

				if f.Kind() == reflect.Pointer {
					f = f.Elem()
				}

				return f.Interface()
			}

			subRes.GetAssocField = func(v reflect.Value) any {
				f := v.FieldByIndex(subRes.AssocFieldIndex)

				if f.IsZero() {
					return nil
				}

				if f.Kind() == reflect.Pointer {
					f = f.Elem()
				}

				return f.Interface()
			}
		}

		subRes.Attach = func(mainModelValue, assocModelValue reflect.Value) {
			field := mainModelValue.Field(i)
			if field.CanSet() {
				if field.Kind() == reflect.Slice {
					field.Set(reflect.Append(field, assocModelValue))
				} else if field.Kind() == reflect.Ptr {
					field.Set(assocModelValue.Addr())
				}
			}
		}

		res[field.Name] = subRes
	}

	registry.Store(modelType, &res)

	return &res, nil
}
