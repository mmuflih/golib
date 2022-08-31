package filter

import (
	"fmt"
	"reflect"
	"strconv"

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
	switch v.Type().Name() {
	case "int":
		return strconv.Itoa(val.(int))
	case "string":
		return val.(string)
	}
	return ""
}

func (w Where) GenerateConditionRaw() string {
	var where string
	var id int
	for field, val := range w {
		for op, v := range val {
			if id == 0 {
				where += " where " + w.generateFilter(field, op, v)
				break
			}
			where += "	and " + w.generateFilter(field, op, v)
		}
		id++
	}
	return where
}

func (w Where) GenerateCondition(db *gorm.DB) *gorm.DB {
	var id int

	for field, val := range w {
		for op, v := range val {
			fmt.Println(field, op, v)
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
		id++
	}
	return db
}
