package razor

import (
	"time"
	"sync"
	"github.com/astaxie/beego/logs"
	"errors"
	"github.com/doggytty/goutils/stringutil"
	"fmt"
)

const (
	KNIFE_ORIGINAL = iota // 初始态
	KNIFE_RUNNING  = iota // 运行中
	KNIFE_SUSPEND  = iota // 暂停状态
	KNIFE_ISOLATE  = iota // 异常态
	KNIFE_OVER     = iota // 结束态

	CREATE = "create"
	DELETE = "delete"
	ADD = "add"
	REMOVE = "remove"
	ERROR = "error"

	KILL = 9
)

var Lock        sync.Mutex     // 共享锁
var createChanel chan int
var deleteChanel chan string
var addChannel chan Hair
var rmChannel chan string
var errChannel chan error

var Instance *KnifeFactory = new(KnifeFactory)
var once sync.Once


// 刀片信息
type Knife struct {
	Id         string       // 当前编号
	Name       string    // 名称
	CreateTime time.Time // 创建时间
	CurrStatus int       // 当前状态
	StartTime  time.Time // 开始运行时间
	EndTime	time.Time	// 结束运行时间
}

// 刀片日志
type KnifeLog struct {
	LogId      int
	KnifeId    string
	JobId      string	// 工作标识
	StartTime  time.Time
	EndTime    time.Time
	Status int 			// 状态
}

// 协程的工厂，一个razor对应一个KnifeFactory实例
type KnifeFactory struct {
	KnifeRef    map[string]*Knife // 每把刀的地址,刀的编号:刀的地址
	KnifeStatus map[string]int // 每把刀的状态,刀的编号:状态
	MaxKnife    int            // 最大
	MinKnife    int            // 最小
	CurrKnife   int            // 当前多少
	Signal 	int	// 信号
}

// 所有外部任务实现该接口
type Hair interface {
	Cutting() bool
}

// 启动razor系统
func init() {
	// 初始化日志
	logs.GetBeeLogger().SetLogger("file", `{"filename":"razor.log"}`)
	// 初始化通信管道
	createChanel = make(chan int, 1)	// 添加刀片通道
	deleteChanel = make(chan string, 1) // 删除刀片通道
	addChannel = make(chan Hair, 1)	// 添加hair通道
	rmChannel = make(chan string, 1)	// 移出hair通道
}

func instanceInit(min, max int) {
	Instance.MinKnife = min
	Instance.MaxKnife = max
	Instance.CurrKnife = 0
	Instance.KnifeRef = make(map[string]*Knife, 0)
	Instance.KnifeStatus = make(map[string]int, 0)
	// 记录状态
	logs.GetBeeLogger().Debug("init razor. min=%d, max=%d, curr=0", min, max)
	// 发送消息给协程,创建新的刀片
	createChanel <- min
	// 启动协程
	go Instance.daemon()
}

func NewInstance(min, max int) (*KnifeFactory, error) {
	if max < min {
		return nil, errors.New("参数错误")
	}
	// 初始化函数
	once.Do(func() {
		instanceInit(min, max)
	})
	return Instance, nil
}

func (kf *KnifeFactory) daemon() {
	// knifefactory 后台
	for {
		if kf.Signal == KILL {
			// 退出关闭factory
			logs.GetBeeLogger().Debug("now close razor!")
			break
		}
		// 获取通道信号并处理
		select {
		case newCount := <- createChanel:
			// 创建刀片
			logs.GetBeeLogger().Debug("now create knife %d!", newCount)
			for i := 0; i < newCount; i++ {
				_, err := kf.CreateKnife()
				if err != nil {
					// 创建失败,重新创建
					createChanel <- 1
				}
			}
		case delKnife := <- deleteChanel:
			// 删除刀片
			logs.GetBeeLogger().Debug("now create knife %s!", delKnife)
			delete(kf.KnifeStatus, delKnife)
			delete(kf.KnifeRef, delKnife)
		case <- addChannel:
			// 新增任务
		case <- rmChannel:
			// 完成任务
		case <- errChannel:
			// 任务错误
		}


		time.Sleep(10)
	}
}

type OutOfMaxError struct {}

func (oe *OutOfMaxError) Error() string {
	return "已经达到最大值,请等待"
}


func (kf *KnifeFactory) CreateKnife() (*Knife, error) {
	// 判断是否已经到max
	if kf.CurrKnife >= kf.MaxKnife {
		logs.GetBeeLogger().Debug("已经达到最大值,请等待!")
		return nil, new(OutOfMaxError)
	}
	// 创建新的knife
	knife := new(Knife)
	tmpString := stringutil.GenerateStringsSize(8)
	knife.Name = fmt.Sprintf("RAZOR-Knife-%s", tmpString)
	knife.CreateTime = time.Now()
	knife.CurrStatus = KNIFE_ORIGINAL
	tmpString = stringutil.GenerateStringsSize(16)
	knife.Id = tmpString

	kf.CurrKnife ++
	kf.KnifeRef[knife.Id] = knife
	kf.KnifeStatus[knife.Id] = KNIFE_ORIGINAL
	return knife, nil
}


