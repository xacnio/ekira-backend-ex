package models

import (
	"fmt"
	"reflect"
	"strings"
)

type Filter map[string]interface{}

func (f *Filter) ToPreloadQuery() []interface{} {
	if reflect.ValueOf(f).IsNil() {
		return []interface{}{""}
	}
	if len(*f) == 0 {
		return []interface{}{""}
	}
	var query []string
	var args = []interface{}{"query"}
	for k, v := range *f {
		typeOfV := reflect.TypeOf(v)
		if typeOfV.Kind() == reflect.Bool {
			boolStr := "false"
			if v.(bool) {
				boolStr = "true"
			}
			query = append(query, fmt.Sprintf("%s is %s", k, boolStr))
		} else {
			query = append(query, fmt.Sprintf("%s = ?", k))
			args = append(args, v)
		}
	}
	args[0] = strings.Join(query, " AND ")
	return args
}

func (f *Filter) ToWhereQuery() (string, []interface{}) {
	if reflect.ValueOf(f).IsNil() {
		return "", []interface{}{}
	}
	if len(*f) == 0 {
		return "", []interface{}{}
	}
	var strs []string
	var args []interface{}
	for k, v := range *f {
		typeOfV := reflect.TypeOf(v)
		if typeOfV.Kind() == reflect.Bool {
			boolStr := "false"
			if v.(bool) {
				boolStr = "true"
			}
			strs = append(strs, fmt.Sprintf("%s is %s", k, boolStr))
		} else {
			strs = append(strs, fmt.Sprintf("%s = ?", k))
			args = append(args, v)
		}
	}
	return strings.Join(strs, " AND "), args
}

type Pagination struct {
	Limit   int    `json:"limit"`
	Page    int    `json:"page"`
	Sort    string `json:"sort"`
	Search  string `json:"search"`
	Filters Filter `json:"filters"`
}
