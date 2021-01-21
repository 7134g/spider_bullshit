package connect

import (
	"fmt"
	"testing"
)

func TestParam_URLString(t *testing.T) {
	m := map[string]interface{}{
		"aaa": 123,
		"bbb": "456",
		"ccc": 7.89,
	}
	p := ParseParam(m)
	u := p.URLString()
	fmt.Println(u)
}

func TestParam_HeaderMap(t *testing.T) {
	m := map[string]interface{}{
		"aaa": 123,
		"bbb": "456",
		"ccc": 7.89,
	}
	p := ParseParam(m)
	h := p.HeaderMap()
	fmt.Println(h)
}
