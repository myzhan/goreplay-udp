package output

import (
	"time"
	"github.com/myzhan/goreplay-udp/stats"
	"github.com/myzhan/goreplay-udp/client"
	"sync/atomic"
	"github.com/myzhan/goreplay-udp/proto"
)

const initialDynamicWorkers = 10

type UDPOutputConfig struct {
	Workers        int
	Timeout        time.Duration
	Stats          bool
	IgnoreResponse bool
}

type UDPOutPut struct {
	// Keep this as first element of struct because it guarantees 64bit
	// alignment. atomic.* functions crash on 32bit machines if operand is not
	// aligned at 64bit. See https://github.com/golang/go/issues/599
	activeWorkers int64

	needWorker chan int

	address string
	queue   chan []byte

	config     *UDPOutputConfig
	queueStats *stats.GorStat
}

func NewUDPOutput(address string, config *UDPOutputConfig) (o *UDPOutPut) {
	o = new(UDPOutPut)
	o.address = address
	o.config = config

	if o.config.Stats {
		o.queueStats = stats.NewGorStat("output_udp")
	}

	o.queue = make(chan []byte, 10000)
	o.needWorker = make(chan int, 1)

	// Initial workers count
	if o.config.Workers == 0 {
		o.needWorker <- initialDynamicWorkers
	} else {
		o.needWorker <- o.config.Workers
	}

	go o.workerMaster()
	return o
}

func (o *UDPOutPut) workerMaster() {
	for {
		newWorkers := <-o.needWorker
		for i := 0; i < newWorkers; i++ {
			go o.startWorker()
		}

		// Disable dynamic scaling if workers poll fixed size
		if o.config.Workers != 0 {
			return
		}
	}
}

func (o *UDPOutPut) startWorker() {
	c := client.NewUDPClient(o.address, o.config.Timeout, o.config.IgnoreResponse)
	deathCount := 0
	atomic.AddInt64(&o.activeWorkers, 1)
	for {
		select {
		case data := <-o.queue:
			o.sendRequest(c, data)
			deathCount = 0
		case <-time.After(time.Millisecond * 100):
			// When dynamic scaling enabled workers die after 2s of inactivity
			if o.config.Workers == 0 {
				deathCount++
			} else {
				continue
			}

			if deathCount > 20 {
				workersCount := atomic.LoadInt64(&o.activeWorkers)

				// At least 1 startWorker should be alive
				if workersCount != 1 {
					atomic.AddInt64(&o.activeWorkers, -1)
					return
				}
			}
		}
	}
}

func (o *UDPOutPut) Write(data []byte) (n int, err error) {
	if !proto.IsRequestPayload(data) {
		return len(data), nil
	}

	buf := make([]byte, len(data))
	copy(buf, data)

	o.queue <- buf

	if o.config.Stats {
		o.queueStats.Write(len(o.queue))
	}

	if o.config.Workers == 0 {
		workersCount := atomic.LoadInt64(&o.activeWorkers)

		if len(o.queue) > int(workersCount) {
			o.needWorker <- len(o.queue)
		}
	}

	return len(data), nil
}

func (o *UDPOutPut) sendRequest(client *client.UDPClient, request []byte) {
	body := proto.PayloadBody(request)
	client.Send(body)
}

func (o *UDPOutPut) String() string {
	return "UDP output: " + o.address
}
