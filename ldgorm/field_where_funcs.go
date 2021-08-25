/*
 * Copyright (C) distroy
 */

package ldgorm

import "reflect"

func Equal(value interface{}) FieldWherer    { return newFieldWhereWithCheck(" = ?", value) }
func NotEqual(value interface{}) FieldWherer { return newFieldWhereWithCheck(" <> ?", value) }

func IsNull() FieldWherer    { return newFieldWhere(" IS NULL") }
func IsNotNull() FieldWherer { return newFieldWhere(" IS NOT NULL") }

func Between(min, max interface{}) FieldWherer {
	minVal, isMinSet := getWhereValue(min)
	maxVal, isMaxSet := getWhereValue(max)
	if !isMinSet && !isMaxSet {
		return fieldWhereEmpty{}
	}
	if !isMinSet {
		return newFieldWhere(" <= ?", maxVal)
	}
	if !isMaxSet {
		return newFieldWhere(" >= ?", minVal)
	}
	if minVal == maxVal {
		return newFieldWhere(" = ?", minVal)
	}

	return newFieldWhere(" BETWEEN ? AND ?", minVal, maxVal)
}

func Gt(value interface{}) FieldWherer {
	return newFieldWhereWithCheck(" > ?", value)
}
func Lt(value interface{}) FieldWherer {
	return newFieldWhereWithCheck(" < ?", value)
}
func Gte(value interface{}) FieldWherer {
	return newFieldWhereWithCheck(" >= ?", value)
}
func Lte(value interface{}) FieldWherer {
	return newFieldWhereWithCheck(" <= ?", value)
}

func In(value interface{}) FieldWherer {
	if value == nil {
		return fieldWhereEmpty{}
	}

	val := reflect.ValueOf(value)
	switch val.Kind() {
	default:
		panic("the condition value for in must be slice or array type")

	case reflect.Slice, reflect.Array:
		break
	}

	if val.Len() == 0 {
		return fieldWhereEmpty{}
	}

	return newFieldWhere(" IN (?)", value)
}

func NotIn(value interface{}) FieldWherer {
	if value == nil {
		return fieldWhereEmpty{}
	}

	val := reflect.ValueOf(value)
	switch val.Kind() {
	default:
		panic("the where value for in must be slice or array type")

	case reflect.Slice, reflect.Array:
		break
	}

	if val.Len() == 0 {
		return fieldWhereEmpty{}
	}

	return newFieldWhere(" NOT IN (?)", value)
}

func Like(value string) FieldWherer {
	return newFieldWhereWithCheck(" LIKE ?", value)
}

func LikePrefix(value string) FieldWherer {
	value = value + "%"
	return newFieldWhereWithCheck(" LIKE ?", value)
}

func LikeSuffix(value string) FieldWherer {
	value = "%" + value
	return newFieldWhereWithCheck(" LIKE ?", value)
}

func LikeContain(value string) FieldWherer {
	value = "%" + value + "%"
	return newFieldWhereWithCheck(" LIKE ?", value)
}

func NotLike(value string) FieldWherer {
	return newFieldWhereWithCheck(" NOT LIKE ?", value)
}

func NotLikePrefix(value string) FieldWherer {
	value = value + "%"
	return newFieldWhereWithCheck(" NOT LIKE ?", value)
}

func NotLikeSuffix(value string) FieldWherer {
	value = "%" + value
	return newFieldWhereWithCheck(" NOT LIKE ?", value)
}

func NotLikeContain(value string) FieldWherer {
	value = "%" + value + "%"
	return newFieldWhereWithCheck(" NOT LIKE ?", value)
}