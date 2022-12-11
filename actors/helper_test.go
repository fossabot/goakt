package actors

import (
	"context"
	"sync"
	"time"

	actorsv1 "github.com/tochemey/goakt/actors/testdata/actors/v1"
	"github.com/tochemey/goakt/log"
)

type BenchActor struct {
	Wg sync.WaitGroup
}

func (p *BenchActor) ID() string {
	return "BenchActor"
}

func (p *BenchActor) PreStart(ctx context.Context) error {
	return nil
}

func (p *BenchActor) Receive(message Message) error {
	switch message.Payload().(type) {
	case *actorsv1.TestSend:
		p.Wg.Done()
	case *actorsv1.TestReply:
		message.SetResponse(&actorsv1.Reply{Content: "received message"})
		p.Wg.Done()
	}
	return nil
}

func (p *BenchActor) PostStop(ctx context.Context) error {
	return nil
}

type TestActor struct {
	id string
}

var _ Actor = (*TestActor)(nil)

// NewTestActor creates a TestActor
func NewTestActor(id string) *TestActor {
	return &TestActor{
		id: id,
	}
}

func (p *TestActor) ID() string {
	return p.id
}

// Init initialize the actor. This function can be used to set up some database connections
// or some sort of initialization before the actor init processing messages
func (p *TestActor) PreStart(ctx context.Context) error {
	return nil
}

// Stop gracefully shuts down the given actor
func (p *TestActor) PostStop(ctx context.Context) error {
	return nil
}

// Receive processes any message dropped into the actor mailbox without a reply
func (p *TestActor) Receive(message Message) error {
	switch message.Payload().(type) {
	case *actorsv1.TestSend:
		// pass
	case *actorsv1.TestPanic:
		log.Panic("Boom")
	case *actorsv1.TestReply:
		message.SetResponse(&actorsv1.Reply{Content: "received message"})
	case *actorsv1.TestTimeout:
		// delay for a while before sending the reply
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			time.Sleep(recvDelay)
			wg.Done()
		}()
		// block until timer is up
		wg.Wait()
	default:
		return ErrUnhandled
	}
	return nil
}

type ParentActor struct {
	id string
}

var _ Actor = (*ParentActor)(nil)

func NewParentActor(id string) *ParentActor {
	return &ParentActor{id: id}
}

func (p *ParentActor) ID() string {
	return p.id
}

func (p *ParentActor) PreStart(ctx context.Context) error {
	return nil
}

func (p *ParentActor) Receive(message Message) error {
	switch message.Payload().(type) {
	case *actorsv1.TestSend:
		return nil
	default:
		return ErrUnhandled
	}
}

func (p *ParentActor) PostStop(ctx context.Context) error {
	return nil
}

type ChildActor struct {
	id string
}

var _ Actor = (*ChildActor)(nil)

func NewChildActor(id string) *ChildActor {
	return &ChildActor{id: id}
}

func (c *ChildActor) ID() string {
	return c.id
}

func (c *ChildActor) PreStart(ctx context.Context) error {
	return nil
}

func (c *ChildActor) Receive(message Message) error {
	switch message.Payload().(type) {
	case *actorsv1.TestSend:
		return nil
	default:
		return ErrUnhandled
	}
}

func (c *ChildActor) PostStop(ctx context.Context) error {
	return nil
}
