package filter

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

/**
 * Created by Muhammad Muflih Kholidin
 * https://github.com/mmuflih
 * muflic.24@gmail.com
 **/

type W map[string]interface{}

type Where map[string]W

func (w Where) generateFilter(field, op string, val interface{}) string {
	if op == "raw" {
		return field + " " + val.(string)
	}
	if val == nil {
		return field + " is null"
	} else {
		return field + " " + op + " " + w.getValue(val)
	}
}

func (w Where) getValue(val interface{}) string {
	v := reflect.ValueOf(val)
	fmt.Println(v.Type().Name())
	switch v.Type().Name() {
	case "int":
		return strconv.Itoa(val.(int))
	case "uint8":
		return strconv.Itoa(int(val.(uint8)))
	case "uint16":
		return strconv.Itoa(int(val.(uint16)))
	case "uint32":
		return strconv.Itoa(int(val.(uint32)))
	case "uint64":
		return strconv.Itoa(int(val.(uint64)))
	case "int8":
		return strconv.Itoa(int(val.(int8)))
	case "int16":
		return strconv.Itoa(int(val.(int16)))
	case "int32":
		return strconv.Itoa(int(val.(int32)))
	case "int64":
		return strconv.Itoa(int(val.(int64)))
	case "float64":
		return fmt.Sprintf("%f", val)
	case "float32":
		return fmt.Sprintf("%f", val)
	case "bool":
		return strconv.FormatBool(val.(bool))
	case "string":
		return "'" + val.(string) + "'"
	default:
		return val.(string)
	}
}

func (w Where) GenerateConditionRaw() string {
	var where string
	var i int
	for field, val := range w {
		if strings.ToLower(field) == "like" {
			var wheres []string
			for op, v := range val {
				wheres = append(wheres, op+" like "+w.getValue(v))
			}
			if i == 0 {
				where += " where (" + strings.Join(wheres, " or ") + ")"
				continue
			}
			where += " and (" + strings.Join(wheres, " or ") + ")"
			i++
			continue
		}
		if strings.ToLower(field) == "ilike" {
			var wheres []string
			for op, v := range val {
				wheres = append(wheres, op+" ilike "+w.getValue(v))
			}
			if i == 0 {
				where += " where (" + strings.Join(wheres, " or ") + ")"
				continue
			}
			where += " and (" + strings.Join(wheres, " or ") + ")"
			i++
			continue
		}
		for op, v := range val {
			where += w.genRawWhere(i, field, op, v)
		}
		i++
	}
	return where
}

func (w Where) genRawWhere(i int, field, op string, v interface{}) string {
	if i == 0 {
		return " where " + w.generateFilter(field, op, v)
	}
	return "	and " + w.generateFilter(field, op, v)
}

func (w Where) GenerateCondition(db *gorm.DB) *gorm.DB {
	for field, val := range w {
		if strings.ToLower(field) == "like" {
			var wheres []string
			for op, v := range val {
				wheres = append(wheres, op+" like "+w.getValue(v))
			}
			db.Where(strings.Join(wheres, " or "))
			continue
		}
		if strings.ToLower(field) == "ilike" {
			var wheres []string
			for op, v := range val {
				wheres = append(wheres, op+" ilike "+w.getValue(v))
			}
			db.Where(strings.Join(wheres, " or "))
			continue
		}
		if val == nil {
			db.Where(field + " is null")
			continue
		}
		for op, v := range val {
			if op == "raw" {
				db.Where(field + " " + w.getValue(v))
				continue
			}
			if v == nil {
				db.Where(field + " is null")
				continue
			} else {
				db.Where(field+" "+op+" ?", v)
				continue
			}

		}
	}
	return db
}
