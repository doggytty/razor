package razor

import (
	"testing"
	"fmt"
	"github.com/doggytty/goutils/stringutil"
)

func TestNewInstance(t *testing.T) {
	factory, err := NewInstance(1, 2)
	if err != nil {
		fmt.Println("error!")
		return
	}
	fmt.Println(factory == nil)
	factory.CreateKnife()

	// 创建协程添加hair


	fmt.Println("OK")

}

func TestStart(t *testing.T) {
	fmt.Println(fmt.Sprintf("RAZOR-Knife-%s", stringutil.GenerateStringsSize(8)))
}
