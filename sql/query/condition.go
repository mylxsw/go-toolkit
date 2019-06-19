package query

import (
	"fmt"
	"strings"
)

type Condition interface {
	WhereColumn(field, operator string, value string) Condition
	OrWhereColumn(field, operator string, value string) Condition
	OrWhereNotExist(subQuery SubQuery) Condition
	OrWhereExist(subQuery SubQuery) Condition
	WhereNotExist(subQuery SubQuery) Condition
	WhereExist(subQuery SubQuery) Condition
	OrWhereNotNull(field string) Condition
	OrWhereNull(field string) Condition
	WhereNotNull(field string) Condition
	WhereNull(field string) Condition
	OrWhereRaw(raw string, items ...interface{}) Condition
	WhereRaw(raw string, items ...interface{}) Condition
	OrWhereNotIn(field string, items ...interface{}) Condition
	OrWhereIn(field string, items ...interface{}) Condition
	WhereNotIn(field string, items ...interface{}) Condition
	WhereIn(field string, items ...interface{}) Condition
	WhereGroup(wc ConditionGroup) Condition
	OrWhereGroup(wc ConditionGroup) Condition
	Where(field, operator string, value interface{}) Condition
	OrWhere(field, operator string, value interface{}) Condition
	WhereCondition(cond sqlCondition) Condition

	When(when When, cg ConditionGroup) Condition
	OrWhen(when When, cg ConditionGroup) Condition

	Clone() Condition
	Empty() bool
	Resolve(tableAlias string) (string, []interface{})
}

type When func() bool

type connectType string
type conditionType int

const (
	connectTypeAnd connectType = "AND"
	connectTypeOr              = "OR"
)

const (
	condTypeSimple conditionType = iota
	condTypeColumn
	condTypeRaw
	condTypeIn
	condTypeNotIn
	condTypeNull
	condTypeNotNull
	condTypeExists
	condTypeNotExists
	condTypeGroup
)

type SubQuery interface {
	ResolveQuery() (string, []interface{})
}

type sqlCondition struct {
	Connector connectType
	Type      conditionType
	Field     string
	Operate   string
	Values    []interface{}
	Nested    ConditionGroup
	SubQuery  SubQuery
	When      When
}

type ConditionGroup func(builder Condition)

type conditionBuilder struct {
	conditions []sqlCondition
}

func ConditionBuilder() Condition {
	return &conditionBuilder{
		conditions: make([]sqlCondition, 0),
	}
}

func (builder *conditionBuilder) When(when When, wc ConditionGroup) Condition {
	return builder.WhereCondition(sqlCondition{
		Connector: connectTypeAnd,
		Type:      condTypeGroup,
		Nested:    wc,
		When:      when,
	})
}

func (builder *conditionBuilder) OrWhen(when When, wc ConditionGroup) Condition {
	return builder.WhereCondition(sqlCondition{
		Connector: connectTypeOr,
		Type:      condTypeGroup,
		Nested:    wc,
		When:      when,
	})
}

func (builder *conditionBuilder) Clone() Condition {
	newBuilder := &conditionBuilder{}
	newBuilder.conditions = make([]sqlCondition, len(builder.conditions))

	for i, c := range builder.conditions {
		newBuilder.conditions[i] = sqlCondition{
			Connector: c.Connector,
			Type:      c.Type,
			Field:     c.Field,
			Operate:   c.Operate,
			Values:    c.Values,
			Nested:    c.Nested,
			SubQuery:  c.SubQuery,
			When:      c.When,
		}
	}

	return newBuilder
}

func (builder *conditionBuilder) Empty() bool {
	return len(builder.conditions) == 0
}

func (builder *conditionBuilder) WhereColumn(field, operator string, value string) Condition {
	return builder.WhereCondition(sqlCondition{
		Connector: connectTypeAnd,
		Type:      condTypeColumn,
		Field:     field,
		Operate:   operator,
		Values:    []interface{}{value},
	})
}

func (builder *conditionBuilder) OrWhereColumn(field, operator string, value string) Condition {
	return builder.WhereCondition(sqlCondition{
		Connector: connectTypeOr,
		Type:      condTypeColumn,
		Field:     field,
		Operate:   operator,
		Values:    []interface{}{value},
	})
}

// --------------

func (builder *conditionBuilder) OrWhereNotExist(subQuery SubQuery) Condition {
	return builder.WhereCondition(sqlCondition{
		Connector: connectTypeOr,
		Type:      condTypeNotExists,
		SubQuery:  subQuery,
	})
}

func (builder *conditionBuilder) OrWhereExist(subQuery SubQuery) Condition {
	return builder.WhereCondition(sqlCondition{
		Connector: connectTypeOr,
		Type:      condTypeExists,
		SubQuery:  subQuery,
	})
}

func (builder *conditionBuilder) WhereNotExist(subQuery SubQuery) Condition {
	return builder.WhereCondition(sqlCondition{
		Connector: connectTypeAnd,
		Type:      condTypeNotExists,
		SubQuery:  subQuery,
	})
}

func (builder *conditionBuilder) WhereExist(subQuery SubQuery) Condition {
	return builder.WhereCondition(sqlCondition{
		Connector: connectTypeAnd,
		Type:      condTypeExists,
		SubQuery:  subQuery,
	})
}

// --------------

func (builder *conditionBuilder) OrWhereNotNull(field string) Condition {
	return builder.WhereCondition(sqlCondition{
		Connector: connectTypeOr,
		Type:      condTypeNotNull,
		Field:     field,
	})
}

func (builder *conditionBuilder) OrWhereNull(field string) Condition {
	return builder.WhereCondition(sqlCondition{
		Connector: connectTypeOr,
		Type:      condTypeNull,
		Field:     field,
	})
}

func (builder *conditionBuilder) WhereNotNull(field string) Condition {
	return builder.WhereCondition(sqlCondition{
		Connector: connectTypeAnd,
		Type:      condTypeNotNull,
		Field:     field,
	})
}

func (builder *conditionBuilder) WhereNull(field string) Condition {
	return builder.WhereCondition(sqlCondition{
		Connector: connectTypeAnd,
		Type:      condTypeNull,
		Field:     field,
	})
}

func (builder *conditionBuilder) OrWhereRaw(raw string, items ...interface{}) Condition {
	return builder.WhereCondition(sqlCondition{
		Connector: connectTypeOr,
		Type:      condTypeRaw,
		Field:     raw,
		Values:    items,
	})
}

func (builder *conditionBuilder) WhereRaw(raw string, items ...interface{}) Condition {
	return builder.WhereCondition(sqlCondition{
		Connector: connectTypeAnd,
		Type:      condTypeRaw,
		Field:     raw,
		Values:    items,
	})
}

func (builder *conditionBuilder) OrWhereNotIn(field string, items ...interface{}) Condition {
	return builder.WhereCondition(sqlCondition{
		Connector: connectTypeOr,
		Type:      condTypeNotIn,
		Field:     field,
		Values:    items,
	})
}

func (builder *conditionBuilder) OrWhereIn(field string, items ...interface{}) Condition {
	return builder.WhereCondition(sqlCondition{
		Connector: connectTypeOr,
		Type:      condTypeIn,
		Field:     field,
		Values:    items,
	})
}

func (builder *conditionBuilder) WhereNotIn(field string, items ...interface{}) Condition {
	return builder.WhereCondition(sqlCondition{
		Connector: connectTypeAnd,
		Type:      condTypeNotIn,
		Field:     field,
		Values:    items,
	})
}

func (builder *conditionBuilder) WhereIn(field string, items ...interface{}) Condition {
	return builder.WhereCondition(sqlCondition{
		Connector: connectTypeAnd,
		Type:      condTypeIn,
		Field:     field,
		Values:    items,
	})
}

func (builder *conditionBuilder) WhereGroup(wc ConditionGroup) Condition {
	return builder.WhereCondition(sqlCondition{
		Connector: connectTypeAnd,
		Type:      condTypeGroup,
		Nested:    wc,
	})
}

func (builder *conditionBuilder) OrWhereGroup(wc ConditionGroup) Condition {
	return builder.WhereCondition(sqlCondition{
		Connector: connectTypeOr,
		Type:      condTypeGroup,
		Nested:    wc,
	})
}

func (builder *conditionBuilder) Where(field, operator string, value interface{}) Condition {
	return builder.WhereCondition(sqlCondition{
		Connector: connectTypeAnd,
		Type:      condTypeSimple,
		Field:     field,
		Operate:   operator,
		Values:    []interface{}{value},
	})
}

func (builder *conditionBuilder) OrWhere(field, operator string, value interface{}) Condition {
	return builder.WhereCondition(sqlCondition{
		Connector: connectTypeOr,
		Type:      condTypeSimple,
		Field:     field,
		Operate:   operator,
		Values:    []interface{}{value},
	})
}

func (builder *conditionBuilder) WhereCondition(cond sqlCondition) Condition {
	if cond.When == nil {
		cond.When = func() bool {
			return true
		}
	}

	builder.conditions = append(builder.conditions, cond)
	return builder
}

func (builder *conditionBuilder) Resolve(tableAlias string) (string, []interface{}) {
	var result = ""
	var params = make([]interface{}, 0)
	for i, c := range builder.conditions {
		if !c.When() {
			continue
		}

		connector := c.Connector
		if i == 0 {
			connector = ""
		}

		r, p := builder.resolveCondition(tableAlias, connector, c)

		result += r
		params = append(params, p...)
	}

	return result, params
}

func (builder *conditionBuilder) resolveCondition(tableAlias string, connector connectType, c sqlCondition) (string, []interface{}) {
	result := ""
	params := make([]interface{}, 0)

	switch c.Type {
	case condTypeSimple:
		if isSubQuery(c.Values) {
			s, p := c.Values[0].(SubQuery).ResolveQuery()
			params = append(params, p...)
			result += fmt.Sprintf(" %s %s %s (%s)", connector, replaceTableField(tableAlias, c.Field), c.Operate, s)
		} else {
			result += fmt.Sprintf(" %s %s %s ?", connector, replaceTableField(tableAlias, c.Field), c.Operate)
			params = append(params, c.Values...)
		}
	case condTypeColumn:
		result += fmt.Sprintf(" %s %s %s %s", connector, replaceTableField(tableAlias, c.Field), c.Operate, replaceTableField(tableAlias, c.Values[0].(string)))
	case condTypeRaw:
		result += fmt.Sprintf(" %s %s", connector, c.Field)
		params = append(params, c.Values...)
	case condTypeIn, condTypeNotIn:
		operator := "IN"
		if c.Type == condTypeNotIn {
			operator = "NOT IN"
		}

		if isSubQuery(c.Values) {
			s, p := c.Values[0].(SubQuery).ResolveQuery()
			result += fmt.Sprintf(" %s %s %s (%s)", connector, replaceTableField(tableAlias, c.Field), operator, s)
			params = append(params, p...)
		} else {
			result += fmt.Sprintf(
				" %s %s %s (%s)",
				connector,
				replaceTableField(tableAlias, c.Field),
				operator,
				strings.Trim(strings.Repeat(", ?", len(c.Values)), ","),
			)

			params = append(params, c.Values...)
		}

	case condTypeNull, condTypeNotNull:
		if c.Type == condTypeNull {
			result += fmt.Sprintf(" %s %s IS NULL", connector, replaceTableField(tableAlias, c.Field))
		} else {
			result += fmt.Sprintf(" %s %s IS NOT NULL", connector, replaceTableField(tableAlias, c.Field))
		}
	case condTypeExists, condTypeNotExists:
		s, p := c.SubQuery.ResolveQuery()
		params = append(params, p...)

		if c.Type == condTypeExists {
			result += fmt.Sprintf(" %s EXISTS (%s)", connector, s)
		} else {
			result += fmt.Sprintf(" %s NOT EXISTS (%s)", connector, s)
		}
	case condTypeGroup:
		newBuilder := ConditionBuilder()
		c.Nested(newBuilder)

		newCond, newParams := newBuilder.Resolve(tableAlias)
		params = append(params, newParams...)
		result += fmt.Sprintf(" %s (%s)", connector, newCond)
	}
	return result, params
}