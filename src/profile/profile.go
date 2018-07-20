package profile

import (
	"runtime"
	"sync"
	"encoding/json"
	"time"
	"os"
	"fmt"
	"bufio"
	"log"
)

var (
	profile  *Profile
	once     sync.Once
	memStats *runtime.MemStats = &runtime.MemStats{}
)

type ProfileData struct {
	// 内部profile信息
	goVersion  string
	numCpu     int
	Alloc      uint64 // 已申请且仍在使用的字节数
	TotalAlloc uint64 // 已申请的总字节数（已释放的部分也算在内）
	Sys        uint64 // 从系统中获取的字节数（下面XxxSys之和）
	Lookups    uint64 // 指针查找的次数
	Mallocs    uint64 // 申请内存的次数
	Frees      uint64 // 释放内存的次数
	// 主分配堆统计
	HeapAlloc    uint64 // 已申请且仍在使用的字节数
	HeapSys      uint64 // 从系统中获取的字节数
	HeapIdle     uint64 // 闲置span中的字节数
	HeapInuse    uint64 // 非闲置span中的字节数
	HeapReleased uint64 // 释放到系统的字节数
	HeapObjects  uint64 // 已分配对象的总个数
	// L低层次、大小固定的结构体分配器统计，Inuse为正在使用的字节数，Sys为从系统获取的字节数
	StackInuse  uint64 // 引导程序的堆栈
	StackSys    uint64
	MSpanInuse  uint64 // mspan结构体
	MSpanSys    uint64
	MCacheInuse uint64 // mcache结构体
	MCacheSys   uint64
	BuckHashSys uint64 // profile桶散列表
	GCSys       uint64 // GC元数据
	OtherSys    uint64 // 其他系统申请
	// 垃圾收集器统计
	NextGC       uint64 // 会在HeapAlloc字段到达该值（字节数）时运行下次GC
	LastGC       uint64 // 上次运行的绝对时间（纳秒）
	PauseTotalNs uint64
	NumGC        uint32

	InUseBytes   int64 // 正在使用的字节数
	InUseObjects int64 // 正在使用的对象数
	NumGoroutine int   // 当前存在的Go程数
}

type Profile struct {
	cycle    int
	dumpFile string
	shutdown chan int
}

func GetProfile(dump string, cycle int) *Profile {
	if cycle < 10 {
		cycle = 10
	}

	if dump == "" {
		dump = "./profile.dump"
	}

	once.Do(func() {
		profile = new(Profile)
	})

	profile.cycle = cycle
	profile.dumpFile = dump

	return profile
}

func (p *Profile) Load() ([]byte, error){
	runtime.ReadMemStats(memStats)

	data := &ProfileData{}

	data.goVersion = runtime.Version()
	data.numCpu = runtime.NumCPU()
	data.Alloc = memStats.Alloc
	data.TotalAlloc = memStats.TotalAlloc
	data.Sys = memStats.Sys
	data.Lookups = memStats.Lookups
	data.Mallocs = memStats.Mallocs
	data.Frees = memStats.Frees
	data.HeapAlloc = memStats.HeapAlloc
	data.HeapSys = memStats.HeapSys
	data.HeapIdle = memStats.HeapIdle
	data.HeapInuse = memStats.HeapInuse
	data.HeapReleased = memStats.HeapReleased
	data.HeapObjects = memStats.HeapObjects
	data.StackInuse = memStats.StackInuse
	data.StackSys = memStats.StackSys
	data.MSpanInuse = memStats.MSpanInuse
	data.MSpanSys = memStats.MSpanSys
	data.MCacheInuse = memStats.MCacheInuse
	data.MCacheSys = memStats.MCacheSys
	data.BuckHashSys = memStats.BuckHashSys
	data.GCSys = memStats.GCSys
	data.OtherSys = memStats.OtherSys
	data.NextGC = memStats.NextGC
	data.LastGC = memStats.LastGC
	data.PauseTotalNs = memStats.PauseTotalNs
	data.NumGC = memStats.NumGC

	js, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	return js, nil
}

func (p *Profile) save(data []byte) {
	f, err := os.OpenFile(p.dumpFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatalf("Error to open file:%v \n", err)
		return
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	writer.Write(data)
	writer.Flush()
}

func (p *Profile) dumpProfile() {
	for {
		select {
			case <- p.shutdown:
				log.Println("Exit profile")
				break
		}
		data, err := p.Load()
		if err != nil {
			fmt.Printf("Load profile json error:%v", err)
		} else {
			p.save(data)
		}
		time.Sleep(time.Millisecond * time.Duration(p.cycle))
	}
}

func (p *Profile) Start() {
	go p.dumpProfile()
}

func (p *Profile) Stop() {
	p.shutdown <- 1
}
