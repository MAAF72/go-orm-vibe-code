package orm

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/iancoleman/strcase"
)

var registry sync.Map

func ParseOrGetMeta(model Model) (*ModelMeta, error) {
	modelType := reflect.TypeOf(model)

	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	if res, ok := registry.Load(modelType); ok {
		return res.(*ModelMeta), nil
	}

	res := ModelMeta{
		TableName:     model.TableName(),
		ModelType:     modelType,
		FieldMetas:    &MapFieldMeta{},
		RelationMetas: &MapRelationMeta{},
	}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		if subRes, _ := ParseRelationMeta(modelType, field, i); subRes != nil {
			(*res.RelationMetas)[field.Name] = subRes
		}

		if subRes, _ := ParseFieldMeta(modelType, field, i); subRes != nil {
			(*res.FieldMetas)[field.Name] = subRes
		}

	}

	registry.Store(modelType, &res)

	return &res, nil
}

func ParseRelationMeta(modelType reflect.Type, field reflect.StructField, fieldIdx int) (*RelationMeta, error) {
	ormTag := field.Tag.Get("orm")
	if ormTag == "" || ormTag == "-" {
		return nil, nil
	}

	fieldType := field.Type
	fieldTypeElem := fieldType.Elem()

	// only for pointer of struct or slice of struct
	if (fieldType.Kind() != reflect.Pointer && fieldType.Kind() != reflect.Slice) || (fieldTypeElem.Kind() != reflect.Struct) {
		return nil, nil
	}

	res := RelationMeta{}

	if iface, ok := reflect.New(fieldTypeElem).Interface().(Model); ok {
		res.AssocTable = iface.TableName()
		res.AssocType = fieldTypeElem
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
			res.Relation = val
		case "assocField":
			res.AssocField = val
		case "mainField":
			res.MainField = val
		}
	}

	{
		mainStructField, ok := modelType.FieldByName(res.MainField)
		if !ok {
			return nil, fmt.Errorf("orm: mainField '%s' specified in tag for %s.%s does not exist in model %s", res.MainField, modelType.Name(), field.Name, modelType.Name())
		}
		res.MainFieldIndex = mainStructField.Index

		assocStructField, ok := res.AssocType.FieldByName(res.AssocField)
		if !ok {
			return nil, fmt.Errorf("orm: assocField '%s' specified in tag for %s.%s does not exist in model %s", res.AssocField, modelType.Name(), field.Name, res.AssocType.Name())
		}
		res.AssocFieldIndex = assocStructField.Index

		res.GetMainField = func(m Model) any {
			v := reflect.ValueOf(m).Elem()
			f := v.FieldByIndex(res.MainFieldIndex)

			if f.IsZero() {
				return nil
			}

			if f.Kind() == reflect.Pointer {
				f = f.Elem()
			}

			return f.Interface()
		}

		res.GetAssocField = func(v reflect.Value) any {
			f := v.FieldByIndex(res.AssocFieldIndex)

			if f.IsZero() {
				return nil
			}

			if f.Kind() == reflect.Pointer {
				f = f.Elem()
			}

			return f.Interface()
		}
	}

	res.Attach = func(mainModelValue, assocModelValue reflect.Value) {
		field := mainModelValue.Field(fieldIdx)
		if field.CanSet() {
			if field.Kind() == reflect.Slice {
				field.Set(reflect.Append(field, assocModelValue))
			} else if field.Kind() == reflect.Ptr {
				field.Set(assocModelValue.Addr())
			}
		}
	}

	return &res, nil
}

func ParseFieldMeta(modelType reflect.Type, field reflect.StructField, fieldIdx int) (*FieldMeta, error) {
	fieldType := field.Type

	if (fieldType.Kind() == reflect.Struct) ||
		(fieldType.Kind() == reflect.Pointer && fieldType.Elem().Kind() == reflect.Struct) ||
		(fieldType.Kind() == reflect.Slice && fieldType.Elem().Kind() == reflect.Struct) {
		return nil, nil
	}

	res := FieldMeta{}

	if ormTag := field.Tag.Get("orm"); ormTag != "" && ormTag != "-" {
		parts := strings.SplitSeq(ormTag, ";")

		for j := range parts {
			subParts := strings.Split(j, ":")

			if len(subParts) != 2 {
				continue
			}

			key, val := subParts[0], subParts[1]

			switch key {
			case "column":
				res.DatabaseName = val
			case "type":
				res.DatabaseType = val
			}
		}
	}

	res.FieldName = field.Name
	res.FieldType = field.Type.Name()

	if res.DatabaseName == "" {
		res.DatabaseName = strcase.ToSnake(res.FieldName)
	}

	if res.DatabaseType == "" {
		res.DatabaseType = "TODO"
	}

	res.GetField = func(m Model) any {
		v := reflect.ValueOf(m).Elem()
		f := v.Field(fieldIdx)

		if f.IsZero() {
			return nil
		}

		return f.Interface()
	}

	return &res, nil
}
