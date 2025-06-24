package orm

import "reflect"

type RelationMeta struct {
	Relation        string
	MainTable       string
	MainType        reflect.Type
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

type MapRelationMeta = map[string]RelationMeta
