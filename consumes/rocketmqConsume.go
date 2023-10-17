package main

import (
	"context"
	"fmt"
	rocketmq "github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"

	"os"
)

func main() {
	sig := make(chan os.Signal)
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName("testGroup"),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
		//consumer.WithAutoCommit(false),
		//	consumer.WithConsumerOrder(true),
		//	consumer.WithConsumerModel(consumer.Clustering),
		//	consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset),
		//consumer.WithMaxReconsumeTimes(2),
	)

	err := c.Subscribe("test", consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range msgs {
			fmt.Println("--------------------------")
			fmt.Printf("subscribe callback: %d---- %s \n", msgs[i].Queue.QueueId, msgs[i].Message.String())

		}

		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	// Note: start after subscribe
	err = c.Start()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	<-sig
	err = c.Shutdown()
	if err != nil {
		fmt.Printf("shutdown Consumer error: %s", err.Error())
	}
}
