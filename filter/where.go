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

type Between map[interface{}]interface{}
type Or map[interface{}]interface{}
type Where map[string]W

func (oO Or) Extract(field string) string {
	for v1, v2 := range oO {
		var vv1 string = " is null "
		var vv2 string = " is null "
		if v1 != nil {
			vv1 = " = " + _getValue(v1)
		}
		if v2 != nil {
			vv2 = " = " + _getValue(v2)
		}
		return "(" + field + vv1 + " or " + field + vv2 + ")"
	}
	return ""
}

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
	return _getValue(val)
}

func _getValue(val interface{}) string {
	if val == nil {
		return "null"
	}
	v := reflect.ValueOf(val)
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
	case "Between":
		for k, vv := range val.(Between) {
			return _getValue(k) + " and " + _getValue(vv)
		}
		return val.(string)
	default:
		return val.(string)
	}
}

func (w Where) GenerateConditionRaw() string {
	var i int
	var wheres []string
	for field, val := range w {
		if val == nil {
			wheres = append(wheres, w.generateFilter(field, " is null ", nil))
			continue
		}
		if strings.ToLower(field) == "like" {
			var whs []string
			for op, v := range val {
				whs = append(wheres, op+" like "+w.getValue(v))
			}
			wheres = append(wheres, " ("+strings.Join(whs, " or ")+")")
			continue
		}
		if strings.ToLower(field) == "ilike" {
			var whs []string
			for op, v := range val {
				whs = append(whs, op+" ilike "+w.getValue(v))
			}
			wheres = append(wheres, " ("+strings.Join(whs, " or ")+")")
			continue
		}
		for op, v := range val {
			if strings.ToLower(op) == "or" {
				oOr, ok := v.(Or)
				if ok {
					wheres = append(wheres, oOr.Extract(field))
					continue
				}
			} else {
				wheres = append(wheres, w.generateFilter(field, op, v))
			}
		}
		i++
	}
	var where string
	for k, wh := range wheres {
		if k == 0 {
			where += " where " + wh
			continue
		}
		where += " and " + wh
	}
	return where
}

func (w Where) genRawWhere(i int, field, op string, v interface{}) string {
	return w.generateFilter(field, op, v)
}

func (w Where) GenerateCondition(db *gorm.DB) *gorm.DB {
	for field, val := range w {
		if val == nil {
			db.Where(field + " is null")
			continue
		}
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
		for op, v := range val {
			if op == "raw" {
				db.Where(field + " " + w.getValue(v))
				continue
			}
			if strings.ToLower(op) == "or" {
				oOr, ok := v.(Or)
				if ok {
					db.Where(oOr.Extract(field))
					continue
				}
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
