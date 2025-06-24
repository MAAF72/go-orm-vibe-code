package orm

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/georgysavva/scany/sqlscan"
	"github.com/gookit/goutil/dump"
	"github.com/iancoleman/strcase"
)

func Select[T any](ctx context.Context, db *sql.DB) ([]T, error) {
	modelInstance := any(new(T)).(Model)

	query := "SELECT * FROM " + modelInstance.TableName()

	dump.Println(query)

	var results []T
	err := sqlscan.Select(ctx, db, &results, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	if len(results) == 0 {
		return results, nil
	}

	err = parseMeta(ctx, db, results...)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func Get[T any](ctx context.Context, db *sql.DB, id int) (*T, error) {
	modelInstance := any(new(T)).(Model)

	query := "SELECT * FROM " + modelInstance.TableName() + " WHERE id = ?"

	dump.Println(query)

	var result T
	tmp := []T{result}

	err := sqlscan.Get(ctx, db, &tmp[0], query, id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	err = parseMeta(ctx, db, tmp...)
	if err != nil {
		return nil, err
	}

	result = tmp[0]

	return &result, nil
}

func parseMeta[T any](ctx context.Context, db *sql.DB, results ...T) error {
	modelInstance := any(new(T)).(Model)

	meta, err := ParseOrGetMeta(modelInstance)
	if err != nil {
		return fmt.Errorf("failed to parse metadata for %s: %w", modelInstance.TableName(), err)
	}

	if meta == nil || len(*meta) == 0 {
		fmt.Println("nil meta")
		return nil
	}

	for _, rel := range *meta {
		mainIDs := make([]string, 0)
		mapMainField := make(map[any][]*T, 0)

		for i := range results {
			if key := rel.GetMainField(any(&results[i]).(Model)); key != nil {
				mapMainField[key] = append(mapMainField[key], &results[i])
				mainIDs = append(mainIDs, fmt.Sprintf("%v", key))
			}
		}

		resultRel := reflect.New(reflect.SliceOf(rel.AssocType))
		queryRel := "SELECT * FROM " + rel.AssocTable + " WHERE " + strcase.ToSnake(rel.AssocField) + " IN (" + strings.Join(mainIDs, ",") + ")"

		dump.Println(queryRel)

		err := sqlscan.Select(ctx, db, resultRel.Interface(), queryRel)
		if err != nil {
			return fmt.Errorf("query failed: %w", err)
		}

		resultRelElem := resultRel.Elem()
		for i := 0; i < resultRelElem.Len(); i++ {
			iElem := resultRelElem.Index(i)

			key := rel.GetAssocField(iElem)
			if key == nil {
				continue
			}

			for j := 0; j < len(mapMainField[key]); j++ {
				rel.Attach(reflect.ValueOf(mapMainField[key][j]).Elem(), iElem)
			}
		}
	}

	return nil
}
