package serializer

import (
	"bullshit"
	"fmt"
	"testing"
)

func TestStructValue(t *testing.T) {
	p := bullshit.NewCtrlParams()
	p.ProxyValue = "127.0.0.1:8888"
	ctrl := bullshit.NewController(p)
	s1 := ctrl.NewSpider()
	s1.ProxyValue = "127.0.0.1:10808"
	StructValue(p, &s1)
	fmt.Println(p.ProxyValue, s1.ProxyValue)
}
