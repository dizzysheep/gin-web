package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
)

var kafkaBrokers []string
var kafkaTopic string

func init() {
	//todo 暂时测试到时候移到配置文件
	kafkaBrokers = []string{"127.0.0.1:9092"}
	kafkaTopic = "topic001"
}

func Produce(sendMsg string) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Version = sarama.V0_10_0_1

	fmt.Println("start make producer")

	// 连接kafka
	producer, err := sarama.NewAsyncProducer(kafkaBrokers, config)
	if err != nil {
		fmt.Println("create producer error, err:", err)
		return
	}
	defer producer.AsyncClose()
	fmt.Println("start goroutine")

	go func(p sarama.AsyncProducer) {
		for {
			select {
			case <-p.Successes():
				//fmt.Println("offset:", suc.Offset, "timestamp:", suc.Timestamp.String(), "partitions:", suc.Partition)
			case fail := <-p.Errors():
				fmt.Println("err : ", fail.Error())
			}
		}
	}(producer)
	// 构造一个消息
	msg := &sarama.ProducerMessage{
		Topic: kafkaTopic,
		Value: sarama.ByteEncoder(sendMsg),
	}
	producer.Input() <- msg
}
