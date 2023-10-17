package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"os"
)

func main() {
	group, err := kafka.NewConsumerGroup(kafka.ConsumerGroupConfig{
		ID:      "my-group",
		Brokers: []string{"localhost:9092", "localhost:9093", "localhost:9094"},
		Topics:  []string{"my-topic2"},
		//StartOffset: kafka.FirstOffset,
	})
	if err != nil {
		fmt.Println("new consume group fail: ", err)
		os.Exit(0)
	}
	defer group.Close()

	ctx := context.Background()
	for {
		next, err := group.Next(ctx)
		if err != nil {
			fmt.Println("for - break ", err)
			break
		}
		assignments := next.Assignments["my-topic2"]
		for _, ass := range assignments {
			partitions, offset := ass.ID, ass.Offset
			fmt.Println(partitions, offset)
			//if partitions == 0 || partitions == 1 {
			//	continue
			//}
			fmt.Println("分区 偏移量：", partitions, offset)
			next.Start(func(ctx context.Context) {
				reader := kafka.NewReader(kafka.ReaderConfig{
					Brokers:   []string{"localhost:9092", "localhost:9093", "localhost:9094"},
					Topic:     "my-topic2",
					Partition: partitions,
				})
				defer reader.Close()
				reader.SetOffset(offset)

				for {
					msg, err := reader.ReadMessage(ctx)
					if err != nil {
						if errors.Is(err, kafka.ErrGenerationEnded) {
							// generation has ended.  commit offsets.  in a real app,
							// offsets would be committed periodically.
							next.CommitOffsets(map[string]map[int]int64{"my-topic2": {partitions: offset + 1}})
							return
						}

						fmt.Printf("error reading message: %+v\n", err)
						return
					}

					fmt.Printf("received message %s/%d/%d : %s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
					offset = msg.Offset
					//提交了偏移量 下次消费的时候 就从之前消费解释后的第一个偏移量开始心的 这个会记录到一个__consume_state的主题
					//记录了 主题 分组 分区  偏移量
					next.CommitOffsets(map[string]map[int]int64{"my-topic2": {partitions: offset + 1}})
				}
			})

		}
	}

}
