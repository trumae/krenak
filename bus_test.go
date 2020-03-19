package bus

import (
	"fmt"
	"log"
	"testing"
)

func TestBase(t *testing.T) {
	bus := NewDefault()
	bus.Start()

	sub1 := bus.Subscribe("test1")
	sub2 := bus.Subscribe("test2")
	sub3 := bus.Subscribe("test3")

	sub1.Publish("type1", 1234)

	msg1 := sub1.Receive()
	fmt.Println(bus)
	msg2 := sub2.Receive()
	fmt.Println(bus)
	msg3 := sub3.Receive()
	fmt.Println(bus)

	if msg1.Content.(int) != 1234 {
		t.Fatal("Wrong message content")
	}
	if msg2.Content.(int) != 1234 {
		t.Fatal("Wrong message content")
	}
	if msg3.Content.(int) != 1234 {
		t.Fatal("Wrong message content")
	}
}

func TestBase2(t *testing.T) {
	bus := NewDefault()
	bus.Start()

	sub1 := bus.Subscribe("test1")
	sub2 := bus.Subscribe("test2")
	sub3 := bus.Subscribe("test3")

	sub1.Publish("type1", 1234)
	sub1.Publish("type1", 5678)

	msg1 := sub1.Receive()
	msg2 := sub2.Receive()
	msg3 := sub3.Receive()

	log.Println(msg1, msg2, msg3)
}
