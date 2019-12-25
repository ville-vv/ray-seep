package monitor

import (
	"github.com/rcrowley/go-metrics"
	"time"
	"vilgo/vlog"
)

type Monitor interface {
	StartPrint(l metrics.Logger, freq time.Duration)
	Inc(int64)
	Dec(int64)
	Gauge(int64)
	Meter(int64)
	Histograms(n int64)
}

func NewMonitor(name string, ut ...string) Monitor {
	return NewBaseMonitor(name, ut...)
}

type BaseMonitor struct {
	name  string
	reg   metrics.Registry
	cnt   metrics.Counter
	gus   metrics.Gauge
	mt    metrics.Meter
	hst   metrics.Histogram
	print metrics.Logger
}

func NewBaseMonitor(name string, ut ...string) *BaseMonitor {
	bm := &BaseMonitor{
		name: name,
		reg:  metrics.NewRegistry(),
		cnt:  metrics.NewCounter(),
		gus:  metrics.NewGauge(),
		mt:   metrics.NewMeter(),
		hst:  metrics.NewHistogram(metrics.NewUniformSample(1028)),
	}
	for _, v := range ut {
		switch v {
		case "counter":
			bm.reg.GetOrRegister(name+"_counter", bm.cnt)
		case "gauge":
			bm.reg.GetOrRegister(name+"_gauge", bm.gus)
		case "meter":
			bm.reg.GetOrRegister(name+"_meter", bm.mt)
		case "hist":
			bm.reg.GetOrRegister(name+"_histogram", bm.hst)
		}
	}
	return bm
}

// 计数类统计，可以进行加或减，也可以进行归零操作，所有的操作都是在旧值的基础上进行的．
func (m *BaseMonitor) Inc(n int64) {
	m.cnt.Inc(n)
}

// 计数类统计进行减也可以进行归零操作，所有的操作都是在旧值的基础上进行的．
func (m *BaseMonitor) Dec(n int64) {
	m.cnt.Dec(n)
}

// 用于对瞬时值的测量，如我们可以过一段时间就对内存的使用量进行统计，并上报，
// 那么所有的数据点集就是对应时间点的内存值，Gauges只有value指标．也就是上报的是什么就是什么
func (m *BaseMonitor) Gauge(n int64) {
	m.gus.Update(n)
}

// 用于计算一段时间内的计量，通常用于计算接口调用频率，如QPS(每秒的次数)，
// 主要分为rateMean,Rate1/Rate5/Rate15等指标．RateMean
// 单位时间内发生的次数，如一分钟发送100次，则该值为100/60.
// Rate1/Rate5/Rate15
func (m *BaseMonitor) Meter(n int64) {
	m.mt.Mark(n)
}

// 主要用于对数据集中的值分布情况进行统计，典型的应用场景为接口耗时，
// 接口每次调用都会产生耗时，记录每次调用耗时来对接口耗时情况进行分析显然不现实．
// 因此将接口一段时间内的耗时看做数据集，
// 并采集Count，Min, Max, Mean, Median, 75%, 95%, 99%等指标．以相对较小的资源消耗，来尽可能反应数据集的真实情况．
func (m *BaseMonitor) Histograms(n int64) {
	m.hst.Update(n)
}

func (m *BaseMonitor) StartPrint(l metrics.Logger, freq time.Duration) {
	go metrics.Log(m.reg, freq, l)
}

type MetricsPrint struct {
	log vlog.ILogger
}

var DefautlMetricePrint = &MetricsPrint{log: vlog.NewGoLogger(&vlog.LogCnf{
	ProgramName:   "ray-seep-metrics",
	OutPutFile:    []string{"./metrics.log"},
	OutPutErrFile: nil,
	Level:         vlog.LogLevelInfo,
})}

func (sel *MetricsPrint) Printf(format string, v ...interface{}) {
	sel.log.LogI(format, v...)
	return
}
