package extracter

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/mylxsw/coll"
)

// Column is a sql column info
type Column struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Rows sql rows object
type Rows struct {
	Columns  []Column        `json:"columns"`
	DataSets [][]interface{} `json:"data_sets"`
}

// Extract export sql rows to Rows object
func Extract(rows *sql.Rows) (*Rows, error) {
	types, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	var columns []Column
	if err := coll.MustNew(types).Map(func(t *sql.ColumnType) Column {
		return Column{
			Name: t.Name(),
			Type: t.DatabaseTypeName(),
		}
	}).All(&columns); err != nil {
		return nil, err
	}

	dataSets := make([][]interface{}, 0)

	for rows.Next() {
		var data = coll.MustNew(types).
			Map(func(t *sql.ColumnType) interface{} {
				var tt interface{}
				return &tt
			}).Items().([]interface{})

		if err := rows.Scan(data...); err != nil {
			return nil, err
		}

		dataSets = append(dataSets, coll.MustNew(data).Map(func(k *interface{}, index int) interface{} {
			if k == nil || *k == nil {
				return nil
			}

			res := fmt.Sprintf("%s", *k)
			switch types[index].DatabaseTypeName() {
			case "INT", "TINYINT", "BIGINT", "MEDIUMINT", "SMALLINT", "DECIMAL":
				intRes, _ := strconv.Atoi(res)
				return intRes
			}

			return res
		}).Items().([]interface{}))
	}

	res := Rows{
		Columns:  columns,
		DataSets: dataSets,
	}

	return &res, nil
}
