package snowflake

import (
	"sync"
	"time"
)

const (
	//起始时间 2020-05-23 00:00:00
	START_STMP int64 = 1590163200000

	//序列号占用的位数
	SEQUENCE_BIT int64 = 12

	//数据中心占用的位数
	DATACENTER_BIT int64 = 5

	//机器标识占用的位数
	MACHINE_BIT int64 = 5

	//每一部分的最大值
	MAX_SEQUENCE       int64 = -1 ^ (-1 << SEQUENCE_BIT)
	MAX_DATACENTER_NUM int64 = -1 ^ (-1 << DATACENTER_BIT)
	MAX_MACHINE_NUM    int64 = -1 ^ (-1 << MACHINE_BIT)

	//每一部分向左的位移
	MACHINE_LEFT    int64 = SEQUENCE_BIT
	DATACENTER_LEFT int64 = MACHINE_LEFT + MACHINE_BIT
	TIMESTMP_LEFT   int64 = DATACENTER_LEFT + DATACENTER_BIT
)

type SnowFlake struct {
	mu           sync.Mutex
	dataCenterId int64
	machineId    int64
	sequence     int64
	lastStmp     int64
}

func NewSnowFlake(dataCenterId int64, machineId int64) *SnowFlake {
	if dataCenterId > MAX_DATACENTER_NUM || dataCenterId < 0 {
		panic("dataCenterId 非法")
	}
	if machineId > MAX_MACHINE_NUM || machineId < 0 {
		panic("machineId 非法")
	}
	return &SnowFlake{dataCenterId: dataCenterId, machineId: machineId, sequence: 0, lastStmp: -1}
}

func (this *SnowFlake) NextId() int64 {
	this.mu.Lock()
	defer this.mu.Unlock()

	currStmp := this.getNewstmp()
	if currStmp < this.lastStmp { //时钟回拨
		panic("时钟回拨")
	}

	if currStmp == this.lastStmp {
		//相同毫秒内，序列号自增
		this.sequence = (this.sequence + 1) & MAX_SEQUENCE

		if this.sequence == int64(0) { //同一毫秒的序列数已经达到最大
			currStmp = this.getNextMill()
		}
	} else {
		//不同毫秒内，序列号置为0
		this.sequence = 0
	}

	this.lastStmp = currStmp
	return (currStmp-START_STMP)<<TIMESTMP_LEFT | //时间戳部分
		this.dataCenterId<<DATACENTER_LEFT | //数据中心部分
		this.machineId<<MACHINE_LEFT | //机器标识部分
		this.sequence //序列号部分
}

func (this *SnowFlake) getNextMill() int64 {
	mill := this.getNewstmp()
	for mill <= this.lastStmp {
		mill = this.getNewstmp()
	}
	return mill
}

func (this *SnowFlake) getNewstmp() int64 {
	return time.Now().UnixNano() / 1e6
}
