package proto_util
import (
	"testing"
	pb "mustard/base/proto_util/testdata"
)
///////init///////
func genCard() pb.Card {
	var card pb.Card
	card.Hobbies = make([]string, 2)
	card.Hobbies[0] = "swimming"
	card.Hobbies[1] = "football"
	card.Keyword = make(map[string]int32)
	card.Keyword["c"] = 22
	card.Keyword["java"] = 33
	card.Person = &pb.Person{"Alice",21}
	return card
}
func TestSerialize(t *testing.T) {
	c := genCard()
	b,_ := Serialize(&c)
	var nc pb.Card
	Deserialize(b,&nc)
	if nc.Hobbies[0] != "swimming" {
		t.Error("Serialize error")
	}
}
func TestSerialize2(t *testing.T) {
	c := genCard()
	s := FromProtoToString(&c)
	var nc pb.Card
	FromStringToProto(s,&nc)
	if nc.Hobbies[0] != "swimming" {
		t.Error("Serialize error")
	}
}

