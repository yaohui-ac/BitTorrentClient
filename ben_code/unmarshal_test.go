package ben_code

import (
	"bytes"
	"fmt"
	"testing"
)

type User struct {
	Name string `bencode:"name"`
	Age  int    `bencode:"age"`
}
type Team struct {
	Name   string `bencode:"name"`
	Size   int    `bencode:"size"`
	Member []User `bencode:"member"`
}

func TestUnmarshal(t *testing.T) {

	str := "d4:name6:archer3:agei29ee"
	uer := &User{}
	_ = Unmarshal(bytes.NewBufferString(str), uer)
	fmt.Println(uer.Age, uer.Name)

	str = "d4:name3:ace4:sizei2e6:memberld4:name6:archer3:agei29eed4:name5:nancy3:agei31eeee"
	team := &Team{}
	err := Unmarshal(bytes.NewBufferString(str), team)
	fmt.Println(err, team.Name, team.Size, len(team.Member))

	buf := new(bytes.Buffer)
	Marshal(buf, team)
	fmt.Println(buf.String())

}
