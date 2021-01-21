package serializer

import (
	"encoding/json"
	"reflect"
)

// 不同结构体间赋值
func StructValue(a, b interface{}) bool {
	//data := StructToMap(a)
	//
	//if MapToStruct(data, b) != nil {
	//	return false
	//}
	//return true

	bytes, err := json.Marshal(a)
	if err != nil {
		return false
	}
	err = json.Unmarshal(bytes, &b)
	if err != nil {
		return false
	}
	return true
}

func StructToMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}
	return data
}
