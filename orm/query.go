package orm

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/georgysavva/scany/sqlscan"
	"github.com/iancoleman/strcase"
)

func Select[T Model](ctx context.Context, db *sql.DB) ([]T, error) {
	var m T
	query := "SELECT * FROM " + m.TableName()

	var results []T
	err := sqlscan.Select(ctx, db, &results, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	if meta, err := ParseOrGetMeta(m); err == nil && meta != nil {
		primaryIDs := make([]string, len(results))
		mapMainIDs := make(map[string]T)
		for i, x := range results {
			primaryIDs[i] = x.GetPK()
			mapMainIDs[primaryIDs[i]] = x
		}

		for _, rel := range *meta {
			mapMainField := make(map[any][]T, 0)

			for _, x := range results {
				key := rel.GetMainField(x)
				mapMainField[key] = append(mapMainField[key], x)
			}

			resultRel := reflect.New(reflect.SliceOf(rel.AssocType))
			queryRel := "SELECT * FROM " + rel.AssocTable + " WHERE " + strcase.ToSnake(rel.AssocField) + " IN (" + strings.Join(primaryIDs, ",") + ")"

			err := sqlscan.Select(ctx, db, resultRel.Interface(), queryRel)
			if err != nil {
				return nil, fmt.Errorf("query failed: %w", err)
			}

			for i := 0; i < resultRel.Elem().Len(); i++ {
				elem := resultRel.Elem().Index(i)

				key := rel.GetAssocField(elem)

				for j := 0; j < len(mapMainField[key]); j++ {
					rel.Attach(reflect.ValueOf(mapMainField[key][j]).Elem(), elem)
				}
			}
		}
	}

	return results, nil
}

func Get[T Model](ctx context.Context, db *sql.DB, id int) (T, error) {
	var m T
	query := "SELECT * FROM " + m.TableName() + " WHERE id = ?"

	fmt.Println(query)

	var result T
	err := sqlscan.Get(ctx, db, &result, query, id)
	if err != nil {
		return m, fmt.Errorf("query failed: %w", err)
	}

	return result, nil
}
