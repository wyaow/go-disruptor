package main

import (
	"fmt"

	"github.com/smartystreets-prototypes/go-disruptor"
)

func main() {
	sequencer, listener := wireup()

	go func() {
		publish(sequencer)
		_ = listener.Close()
	}()

	listener.Listen()
}
func wireup() (disruptor.Sequencer, disruptor.ListenCloser) {
	wireup, err := disruptor.New(
		disruptor.WithCapacity(BufferSize),
		disruptor.WithConsumerGroup(MyConsumer{}, MyConsumer{}, MyConsumer{}, MyConsumer{}, MyConsumer{}),
		disruptor.WithConsumerGroup(MyConsumer{}, MyConsumer{}, MyConsumer{}, MyConsumer{}, MyConsumer{}),
		disruptor.WithConsumerGroup(MyConsumer{}, MyConsumer{}, MyConsumer{}, MyConsumer{}, MyConsumer{}),
		disruptor.WithConsumerGroup(MyConsumer{}, MyConsumer{}, MyConsumer{}, MyConsumer{}, MyConsumer{}),
		disruptor.WithConsumerGroup(MyConsumer{}, MyConsumer{}, MyConsumer{}, MyConsumer{}, MyConsumer{}),
	)
	if err != nil {
		panic(err)
	}

	return wireup.Build()
}

func publish(sequencer disruptor.Sequencer) {
	for sequence := int64(0); sequence <= Iterations; {
		sequence = sequencer.Reserve(Reservations)

		for lower := sequence - Reservations + 1; lower <= sequence; lower++ {
			ringBuffer[lower&BufferMask] = lower
		}

		sequencer.Commit(sequence-Reservations+1, sequence)
	}
}

// ////////////////////

type MyConsumer struct{}

func (this MyConsumer) Consume(lower, upper int64) {
	for ; lower <= upper; lower++ {
		message := ringBuffer[lower&BufferMask]
		if message != lower {
			panic(fmt.Errorf("race condition: %d %d", message, lower))
		}
	}
}

const (
	BufferSize   = 1024 * 64
	BufferMask   = BufferSize - 1
	Iterations   = 128 * 1024 * 32
	Reservations = 1
)

var ringBuffer = [BufferSize]int64{}
