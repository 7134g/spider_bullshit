package connect

import "strconv"

// 只支持 int, int64, float, float64 和 string
type Param map[string]interface{}

func (u Param) URLString() string {
	params := ""
	for k, v := range u {
		if len(params) != 0 {
			params += "&"
		}
		params += k
		params += "=" + u.ConvertString(v)
	}
	return params
}

func (u Param) HeaderMap() map[string]string {
	header := make(map[string]string, 0)
	for k, v := range u {
		header[k] = u.ConvertString(v)
	}
	return header
}

func (u Param) Set(k string, v interface{}) {
	u[k] = u.ConvertString(v)
}

func (u Param) GetString(k string) string {
	return u.ConvertString(u[k])
}

func (u Param) ConvertString(v interface{}) string {
	value := ""
	switch v.(type) {
	case string:
		value = v.(string)
	case int:
		value = strconv.Itoa(v.(int))
	case bool:
		value = strconv.FormatBool(v.(bool))
	case int64:
		value = strconv.FormatInt(v.(int64), 10)
	case float64:
		value = strconv.FormatFloat(v.(float64), 'f', -1, 64)
	case uint64:
		value = strconv.FormatUint(v.(uint64), 10)
	default:
		panic("Param 无法解析该类型")
	}
	return value
}

// 将 map[string]interface{} 转 Param
func ParseParam(m map[string]interface{}) Param {
	return Param(m)
}
