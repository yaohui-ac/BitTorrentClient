package ben_code

import (
	"bytes"
	"fmt"
	"testing"
)

type testMarshal struct {
	val  []int64
	name string `bencode:"name"`
	Age  int32  `bencode:"age"`
}

func newTestMarshal() *testMarshal {
	return &testMarshal{
		val:  []int64{114, 514},
		name: "ikun",
		Age:  114514,
	}
}
func TestMarshal(t *testing.T) {
	str := "123"
	buf := new(bytes.Buffer)
	_ = Marshal(buf, str)
	fmt.Println(buf.String())
	buf.Reset()

	v := 123
	_ = Marshal(buf, v)
	fmt.Println(buf.String())
	buf.Reset()

	list := []string{"a", "AA", "BBB"}
	_ = Marshal(buf, list)
	fmt.Println(buf.String())
	buf.Reset()

	_ = Marshal(buf, newTestMarshal())
	fmt.Println(buf.String())
	buf.Reset()
}
