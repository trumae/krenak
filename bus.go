package bus

import (
	"fmt"
	"log"
	"time"
)

const DefaultCapacity = 10

type Message struct {
	Type      string
	Timestamp int64
	From      string
	To        string
	Content   interface{}
}

type Subscription struct {
	Name      string
	Chan      chan Message
	bus       *Bus
	sent      int64
	delivered int64
	published int64
}

type Bus struct {
	input    chan Message
	subs     []*Subscription
	capacity int64
}

func NewWithCapacity(capacity int64) *Bus {
	return &Bus{
		input:    make(chan Message, capacity),
		subs:     []*Subscription{},
		capacity: capacity,
	}
}

func NewDefault() *Bus {
	return NewWithCapacity(DefaultCapacity)
}

func (bus *Bus) Start() {
	go func() {
		for {
			msg := <-bus.input
			switch {
			case msg.Type == "subscribe":
				bus.subs = append(bus.subs, msg.Content.(*Subscription))
			default:
				for _, sub := range bus.subs {
					if sub.sent-sub.delivered < bus.capacity-1 {
						sub.Chan <- msg
						sub.sent++
					} else {
						log.Println("Capacity exhausted - Subscription:", sub.Name)
					}
				}
			}
		}
	}()
}

func (bus *Bus) Subscribe(name string) *Subscription {
	sub := &Subscription{
		Name: name,
		Chan: make(chan Message),
		bus:  bus,
	}
	msg := Message{
		Type:    "subscribe",
		Content: sub,
	}
	bus.input <- msg
	return sub
}

func (sub *Subscription) Publish(typ string, cont interface{}) {
	msg := Message{
		Type:      typ,
		From:      sub.Name,
		Timestamp: time.Now().UnixNano(),
		Content:   cont,
	}
	sub.bus.input <- msg
	sub.published++
}

func (sub *Subscription) Receive() Message {
	msg := <-sub.Chan
	sub.delivered++
	return msg
}

func (sub *Subscription) String() string {
	return fmt.Sprintf("%s - sent: %d delivered: %d published: %d intoqueue: %d",
		sub.Name, sub.sent, sub.delivered, sub.published, sub.sent-sub.delivered)
}

func (bus *Bus) String() string {
	ret := fmt.Sprintf("Capacity: %d\n", bus.capacity)
	for _, sub := range bus.subs {
		ret += sub.String() + "\n"
	}
	return ret
}
