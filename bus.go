package bus

import (
	"fmt"
	"log"
	"sync"
	"time"
)

//DefaultCapacity is the normal capacity to subscrition stack
const DefaultCapacity = 10

//Message is the msg format to bus
type Message struct {
	Type      string
	Timestamp int64
	From      string
	To        string
	Content   interface{}
}

//Subscription listen and can publish in the bus
type Subscription struct {
	Name      string
	Chan      chan Message
	bus       *Bus
	sent      int64
	delivered int64
	published int64
	lock      sync.Mutex
}

//Bus delivers messages to subscritions
type Bus struct {
	input    chan Message
	subs     []*Subscription
	capacity int64

	lock sync.Mutex
}

//NewWithCapacity create a new Bus using stack with capacity msgs
func NewWithCapacity(capacity int64) *Bus {
	return &Bus{
		input:    make(chan Message, capacity),
		subs:     []*Subscription{},
		capacity: capacity,
		lock:     sync.Mutex{},
	}
}

//NewDefault create a new Bus with default configuration
func NewDefault() *Bus {
	return NewWithCapacity(DefaultCapacity)
}

//Start run with Bus
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
						func() {
							sub.lock.Lock()
							defer sub.lock.Unlock()

							sub.sent++
						}()
						sub.Chan <- msg
					} else {
						log.Println("Capacity exhausted - Subscription:", sub.Name)
					}
				}
			}
		}
	}()
}

//Subscribe return a subscription to bus
func (bus *Bus) Subscribe(name string) *Subscription {
	sub := &Subscription{
		Name: name,
		Chan: make(chan Message),
		bus:  bus,
		lock: sync.Mutex{},
	}
	msg := Message{
		Type:    "subscribe",
		Content: sub,
	}
	bus.input <- msg
	return sub
}

//Publish can inject a message into the bus
func (sub *Subscription) Publish(typ string, cont interface{}) {
	msg := Message{
		Type:      typ,
		From:      sub.Name,
		Timestamp: time.Now().UnixNano(),
		Content:   cont,
	}
	sub.bus.input <- msg
	func() {
		sub.lock.Lock()
		defer sub.lock.Unlock()

		sub.published++
	}()

}

//Receive get a message from bs
func (sub *Subscription) Receive() Message {
	msg := <-sub.Chan
	func() {
		sub.lock.Lock()
		defer sub.lock.Unlock()

		sub.delivered++
	}()
	return msg
}

//String show info about a sub
func (sub *Subscription) String() string {
	sub.lock.Lock()
	defer sub.lock.Unlock()

	ret := fmt.Sprintf("%s - sent: %d delivered: %d published: %d intoqueue: %d",
		sub.Name, sub.sent, sub.delivered, sub.published, sub.sent-sub.delivered)
	return ret
}

//String show basic info about the bus
func (bus *Bus) String() string {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	ret := fmt.Sprintf("Capacity: %d\n", bus.capacity)
	for _, sub := range bus.subs {
		ret += sub.String() + "\n"
	}
	return ret
}
