package main

import (
	"encoding/json"
	"fmt"

	pf "github.com/dakaraj/go-api/chapter06/protofiles"
	"github.com/golang/protobuf/proto"
)

func main() {
	p := &pf.Person{
		Id:    1234,
		Name:  "Anton K",
		Email: "ankra@test.com",
		Phones: []*pf.Person_PhoneNumber{
			&pf.Person_PhoneNumber{Number: "555-4321", Type: pf.Person_HOME},
		},
	}

	p1 := &pf.Person{}
	body, _ := proto.Marshal(p)
	jsonBody, _ := json.Marshal(p)
	_ = proto.Unmarshal(body, p1)
	fmt.Printf("Original struct loaded from proto file: %v\n", p)
	fmt.Printf("Marshalled JSON: %s\n", jsonBody)
	fmt.Printf("Marshalled proto data: %v\n", body)
	fmt.Printf("Unmarshalled struct: %v\n", p1)
}
