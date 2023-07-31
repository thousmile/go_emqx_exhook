package channelx

import (
	"log"
	"runtime"
	"sync"
	"time"
)

// Aggregator Represents the aggregator
type Aggregator[T any] struct {
	option         AggregatorOption[T]
	wg             *sync.WaitGroup
	quit           chan struct{}
	eventQueue     chan T
	batchProcessor BatchProcessFunc[T]
}

// AggregatorOption Represents the aggregator option
type AggregatorOption[T any] struct {
	BatchSize         int
	Workers           int
	ChannelBufferSize int
	LingerTime        time.Duration
	ErrorHandler      ErrorHandlerFunc[T]
	Logger            *log.Logger
}

// BatchProcessFunc the func to batch process items
type BatchProcessFunc[T any] func([]T) error

// SetAggregatorOptionFunc the func to set option for aggregator
type SetAggregatorOptionFunc[T any] func(option AggregatorOption[T]) AggregatorOption[T]

// ErrorHandlerFunc the func to handle error
type ErrorHandlerFunc[T any] func(err error, items []T, batchProcessFunc BatchProcessFunc[T], aggregator *Aggregator[T])

// NewAggregator Creates a new aggregator
func NewAggregator[T any](batchProcessor BatchProcessFunc[T], optionFuncs ...SetAggregatorOptionFunc[T]) *Aggregator[T] {
	option := AggregatorOption[T]{
		BatchSize:  100,
		Workers:    2,
		LingerTime: 1 * time.Second,
	}

	for _, optionFunc := range optionFuncs {
		option = optionFunc(option)
	}

	if option.ChannelBufferSize <= option.Workers {
		option.ChannelBufferSize = option.Workers
	}

	return &Aggregator[T]{
		eventQueue:     make(chan T, option.ChannelBufferSize),
		option:         option,
		quit:           make(chan struct{}),
		wg:             new(sync.WaitGroup),
		batchProcessor: batchProcessor,
	}
}

// TryEnqueue Try enqueue an item, and it is non-blocked
func (agt *Aggregator[T]) TryEnqueue(item T) bool {
	select {
	case agt.eventQueue <- item:
		return true
	default:
		if agt.option.Logger != nil {
			agt.option.Logger.Printf("Aggregator: Event queue is full and try reschedule")
		}
		runtime.Gosched()
		select {
		case agt.eventQueue <- item:
			return true
		default:
			if agt.option.Logger != nil {
				agt.option.Logger.Printf("Aggregator: Event queue is still full and %+v is skipped.", item)
			}
			return false
		}
	}
}

// Enqueue an item, will be blocked if the queue is full
func (agt *Aggregator[T]) Enqueue(item T) {
	agt.eventQueue <- item
}

// Start the aggregator
func (agt *Aggregator[T]) Start() {
	for i := 0; i < agt.option.Workers; i++ {
		index := i
		go agt.work(index)
	}
}

// Stop the aggregator
func (agt *Aggregator[T]) Stop() {
	close(agt.quit)
	agt.wg.Wait()
}

// SafeStop Stop the aggregator safely, the difference with Stop is it guarantees no item is missed during stop
func (agt *Aggregator[T]) SafeStop() {
	if len(agt.eventQueue) == 0 {
		close(agt.quit)
	} else {
		ticker := time.NewTicker(50 * time.Millisecond)
		for range ticker.C {
			if len(agt.eventQueue) == 0 {
				close(agt.quit)
				break
			}
		}
		ticker.Stop()
	}
	agt.wg.Wait()
}

func (agt *Aggregator[T]) work(index int) {
	defer func() {
		if r := recover(); r != nil {
			if agt.option.Logger != nil {
				agt.option.Logger.Printf("Aggregator: recover worker as bad thing happens %+v", r)
			}
			agt.work(index)
		}
	}()

	agt.wg.Add(1)
	defer agt.wg.Done()

	batch := make([]T, 0, agt.option.BatchSize)
	lingerTimer := time.NewTimer(0)
	if !lingerTimer.Stop() {
		<-lingerTimer.C
	}
	defer lingerTimer.Stop()

loop:
	for {
		select {
		case req := <-agt.eventQueue:
			batch = append(batch, req)

			batchSize := len(batch)
			if batchSize < agt.option.BatchSize {
				if batchSize == 1 {
					lingerTimer.Reset(agt.option.LingerTime)
				}
				break
			}

			agt.batchProcess(batch)

			if !lingerTimer.Stop() {
				<-lingerTimer.C
			}
			batch = make([]T, 0, agt.option.BatchSize)
		case <-lingerTimer.C:
			if len(batch) == 0 {
				break
			}

			agt.batchProcess(batch)
			batch = make([]T, 0, agt.option.BatchSize)
		case <-agt.quit:
			if len(batch) != 0 {
				agt.batchProcess(batch)
			}

			break loop
		}
	}
}

func (agt *Aggregator[T]) batchProcess(items []T) {
	agt.wg.Add(1)
	defer agt.wg.Done()
	if err := agt.batchProcessor(items); err != nil {
		if agt.option.Logger != nil {
			agt.option.Logger.Printf("Aggregator: error happens")
		}
		if agt.option.ErrorHandler != nil {
			go agt.option.ErrorHandler(err, items, agt.batchProcessor, agt)
		} else if agt.option.Logger != nil {
			agt.option.Logger.Printf("Aggregator: error happens in batchProcess and is skipped")
		}
	} else if agt.option.Logger != nil {
		agt.option.Logger.Printf("Aggregator: %d items have been sent.", len(items))
	}
}
