package razor

import (
	"time"
	"sync"
)

const (
	KNIFE_ORIGINAL = iota // 初始态
	KNIFE_RUNNING  = iota // 运行中
	KNIFE_SUSPEND  = iota // 暂停状态
	KNIFE_ISOLATE  = iota // 异常态
	KNIFE_OVER     = iota // 结束态
)

type Knife struct {
	Id         int       //当前编号
	Name       string    // 名称
	CreateTime time.Time //创建时间
	CurrStatus int       //当前状态
	StartTime  time.Time //开始运行时间
}

type KnifeLog struct {
	LogId      int
	KnifeId    int
	JobId      string	// 工作标识
	StartTime  time.Time
	EndTime    time.Time
	CurrStatus int //当前状态
}

// 协程的工厂，一个razor对应一个KnifeFactory实例
type KnifeFactory struct {
	KnifeRef    []*Knife       // 每把刀的地址
	KnifeStatus map[int]string // 每把刀的状态,刀的编号:状态
	MaxKnife    int            // 最大
	MinKnife    int            // 最小
	CurrKnife   int            // 当前多少
	Lock        sync.Mutex     // 共享锁
}

// 启动razor系统
func init() {
	// 启动
	//factory := new(KnifeFactory)
	//factory.KnifeRef = make([]*Knife, 0)
}

// 启动razor
func Start() {
	// 获取当前处理的值，与预设的系统最大值做比较
	// 如果
}

// 停止razor
func Close() {

}
