package controllers

import (
	"context"
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/gin-gonic/gin"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"strconv"

	kafka "github.com/segmentio/kafka-go"

	rocketmq "github.com/apache/rocketmq-client-go/v2"
)

type TestController struct {
}

func NewTestController() TestController {
	return TestController{}
}
func (t TestController) Index() gin.HandlerFunc {
	return func(context *gin.Context) {
		sentil, blockError := sentinel.Entry("some-test", sentinel.WithTrafficType(base.Inbound))

		if blockError != nil {
			context.JSON(400, gin.H{
				"message": "访问太频繁,已被限流",
			})
			context.Abort()
			return
		}
		//正常请求后需要关闭
		defer func() {
			sentil.Exit()
		}()

		context.JSON(200, gin.H{
			"message": "index",
			"data":    "haha",
		})

	}
}

func (t TestController) Test() gin.HandlerFunc {
	return func(context *gin.Context) {

		context.JSON(200, gin.H{
			"message": "test",
			"data":    "haha",
		})
	}
}

func (t TestController) Rocketmq() gin.HandlerFunc {
	return func(c *gin.Context) {
		p, _ := rocketmq.NewProducer(
			producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
			producer.WithRetry(2),
		)
		err := p.Start()
		if err != nil {
			fmt.Printf("start producer error: %s", err.Error())
			//os.Exit(1)
			c.JSON(200, gin.H{
				"message": err.Error(),
			})
			return
		}
		topic := "test"

		ctx := context.Background()
		start := time.Now()
		for i := 0; i < 9; i++ {
			msg := &primitive.Message{
				Topic: topic,
				Body:  []byte("Hello RocketMQ Go Client! " + strconv.Itoa(i)),
			}
			//msg.WithDelayTimeLevel(2)
			_, err := p.SendSync(ctx, msg)

			if err != nil {
				fmt.Printf("send message error: %s\n", err)
			} else {
				//fmt.Printf("send message success: result=%s\n", res.String())
			}
		}

		err = p.Shutdown()
		if err != nil {
			c.JSON(200, gin.H{
				"message": err.Error(),
			})
			return
			fmt.Printf("s"+
				"hutdown producer error: %s", err.Error())
		}
		c.JSON(200, gin.H{
			"message": "ok",
			"time":    time.Since(start).Seconds(),
		})
	}

}

type MyTranListener struct {
	localTrans       *sync.Map
	transactionIndex int32
}

func NewMyTranListener() *MyTranListener {
	return &MyTranListener{
		localTrans: new(sync.Map),
	}
}

func (m *MyTranListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	nextIndex := atomic.AddInt32(&m.transactionIndex, 1)
	fmt.Printf("nextIndex: %v for transactionID: %v\n", nextIndex, msg.TransactionId)
	status := nextIndex % 3
	m.localTrans.Store(msg.TransactionId, primitive.LocalTransactionState(status+1))

	fmt.Printf("dl")
	return primitive.CommitMessageState
}

// When no response to prepare(half) message. broker will send check message to check the transaction status, and this
// method will be invoked to get local transaction status.
func (m *MyTranListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	fmt.Printf("%v msg transactionID : %v\n", time.Now(), msg.TransactionId)
	v, existed := m.localTrans.Load(msg.TransactionId)
	if !existed {
		fmt.Printf("unknow msg: %v, return Commit", msg)
		return primitive.CommitMessageState
	}
	state := v.(primitive.LocalTransactionState)
	switch state {
	case 1:
		fmt.Printf("checkLocalTransaction COMMIT_MESSAGE: %v\n", msg)
		return primitive.CommitMessageState
	case 2:
		fmt.Printf("checkLocalTransaction ROLLBACK_MESSAGE: %v\n", msg)
		return primitive.RollbackMessageState
	case 3:
		fmt.Printf("checkLocalTransaction unknow: %v\n", msg)
		return primitive.UnknowState
	default:
		fmt.Printf("checkLocalTransaction default COMMIT_MESSAGE: %v\n", msg)
		return primitive.CommitMessageState
	}
}

func (t TestController) RocketmqTran() gin.HandlerFunc {
	return func(c *gin.Context) {
		p, _ := rocketmq.NewTransactionProducer(
			NewMyTranListener(),
			producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
			producer.WithRetry(2),
		)
		err := p.Start()
		if err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}

		for i := 0; i < 1; i++ {
			res, err := p.SendMessageInTransaction(context.Background(), primitive.NewMessage("test", []byte("this is test tran message"+strconv.Itoa(i))))

			if err != nil {
				fmt.Printf("send message error: %s\n", err)
			} else {
				fmt.Printf("send message success: result=%s\n", res.String())
			}
		}

		time.Sleep(30 * time.Second)
		err = p.Shutdown()
		if err != nil {
			fmt.Printf("shutdown producer error: %s", err.Error())
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "send tran message ok",
		})
	}
}

func (t TestController) Kafka() gin.HandlerFunc {
	return func(context *gin.Context) {
		conn, err := kafka.DialLeader(context, "tcp", "localhost:9092", "my-topic2", 2)

		if err != nil {
			log.Fatal("faid to dial leader:", err)

		}
		defer conn.Close()
		//conn.SetRequiredAcks(1)

		start := time.Now()
		mess := make([]kafka.Message, 0)
		for i := 0; i < 120; i++ {
			//messa := kafka.Message{
			//	Key:   []byte("order"+strconv.Itoa(i)),
			//	Value: []byte(strconv.Itoa(i)),
			//}
			mess = append(mess, kafka.Message{
				Key:   []byte("order" + strconv.Itoa(i)),
				Value: []byte(strconv.Itoa(i)),
			})

		}

		for j := 0; j < 60; j++ {
			_, err = conn.WriteMessages(mess...)
			if err != nil {
				context.JSON(400, gin.H{
					"message": "fail to send 1k afka message: " + err.Error(),
				})
				return
			}
		}

		cost := time.Since(start).Seconds()
		if err != nil {
			context.JSON(400, gin.H{
				"message": "fail to send kafka message: " + err.Error(),
				"cost":    cost,
			})
			return
		}
		context.JSON(200, gin.H{
			"message": "success to send kafka message",
			"cost":    cost,
		})

	}
}
