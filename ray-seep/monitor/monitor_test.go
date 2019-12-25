package monitor

import (
	"github.com/rcrowley/go-metrics"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestBaseMonitor_Inc(t *testing.T) {
	m := NewBaseMonitor("test", "counter")
	m.Inc(1)
	m.Inc(1)
	m.Inc(1)
	m.Inc(1)
	go metrics.Log(m.reg, 5*time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
	time.Sleep(time.Second * 20)
}
func TestBaseMonitor_Gauge(t *testing.T) {
	m := NewBaseMonitor("test", "gauge")
	go func() {
		cnt := 0
		for {
			m.Gauge(int64(cnt))
			time.Sleep(time.Second * 1)
			if cnt >= 19 {
				return
			}
			cnt++
		}
	}()

	go metrics.Log(m.reg, 5*time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
	time.Sleep(time.Second * 20)

}

func TestBaseMonitor_Meter(t *testing.T) {
	m := NewBaseMonitor("test", "meter")
	rand.Seed(time.Now().Unix())
	go func() {
		for {
			m.Meter(rand.Int63n(1000))
			time.Sleep(time.Millisecond * 10)
		}
	}()
	m.StartPrint(DefautlMetricePrint, 2*time.Second)
	time.Sleep(time.Second * 20)
}
