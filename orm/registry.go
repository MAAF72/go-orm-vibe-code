package orm

import (
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

		// only for pointer of struct or slice of struct
		if (field.Type.Kind() != reflect.Pointer && field.Type.Kind() != reflect.Slice) || (field.Type.Elem().Kind() != reflect.Struct) {
			continue
		}

		subRes := RelationMeta{}

		subRes.MainTable = model.TableName()
		subRes.MainType = modelType

		if iface, ok := reflect.New(field.Type.Elem()).Elem().Addr().Interface().(Model); ok {
			subRes.AssocTable = iface.TableName()
			subRes.AssocType = field.Type.Elem()
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

		subRes.GetMainField = func(m Model) any {
			v := reflect.ValueOf(m).Elem()

			return v.FieldByName(subRes.MainField).Interface()
		}

		subRes.GetAssocField = func(v reflect.Value) any {
			return v.FieldByName(subRes.AssocField).Interface()
		}

		subRes.Attach = func(primary, foreign reflect.Value) {
			field := primary.Field(i)

			if field.CanSet() {
				if field.Kind() == reflect.Slice {
					field.Set(reflect.Append(field, foreign))
				} else {
					field.Set(foreign.Addr())
				}
			}
		}

		res[field.Name] = subRes
	}

	registry.Store(modelType, &res)

	return &res, nil
}
