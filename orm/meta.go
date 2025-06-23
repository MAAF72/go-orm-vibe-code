package orm

import "reflect"

type RelationMeta struct {
	Relation     string
	MainTable    string
	MainType     reflect.Type
	MainField    string
	GetMainField func(Model) any
	// MainKey    string
	AssocTable    string
	AssocType     reflect.Type
	AssocField    string
	GetAssocField func(reflect.Value) any
	// AssocKey   string
	Attach func(reflect.Value, reflect.Value) // attach $2 to $1
}

type MapRelationMeta = map[string]RelationMeta
