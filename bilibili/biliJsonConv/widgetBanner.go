package biliJsonConv

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
)

/*!
这个关键字的JSON结构的数量多且复杂，用 map[string]interface{} 可能是更好的选择
*/

func WidgetBanner(b []byte) (map[string]interface{}, error) {
	var err error
	var wbs map[string]interface{}

	decoder := jsoniter.NewDecoder(bytes.NewReader(b))
	decoder.UseNumber()
	err = decoder.Decode(&wbs)

	return wbs, err
}
