package inits

import (
	"fmt"
	"github.com/go-zookeeper/zk"
	"sync"
	"time"
)

var ZkConn *zk.Conn
var once sync.Once

func InitZk() {
	once.Do(func() {

		connect, _, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second*3)
		if err != nil {
			fmt.Println("zookeeper init err:", err)
		}
		ZkConn = connect
	})
}
