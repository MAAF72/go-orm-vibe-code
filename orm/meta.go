package orm

import "reflect"

type ModelMeta struct {
	TableName     string
	ModelType     reflect.Type
	RelationMetas *MapRelationMeta
	FieldMetas    *MapFieldMeta
}

type RelationMeta struct {
	Relation        string
	MainField       string
	MainFieldIndex  []int
	GetMainField    func(Model) any
	AssocTable      string
	AssocType       reflect.Type
	AssocField      string
	AssocFieldIndex []int
	GetAssocField   func(reflect.Value) any
	Attach          func(reflect.Value, reflect.Value) // attach $2 to $1
}

type MapRelationMeta = map[string]*RelationMeta

type FieldMeta struct {
	FieldName    string
	FieldType    string
	GetField     func(Model) any
	DatabaseName string
	DatabaseType string
}

type MapFieldMeta = map[string]*FieldMeta
