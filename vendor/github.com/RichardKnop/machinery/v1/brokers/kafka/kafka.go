package kafka

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/RichardKnop/machinery/v1/brokers/errs"
	"github.com/RichardKnop/machinery/v1/brokers/iface"
	"github.com/RichardKnop/machinery/v1/common"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/tasks"
	kafka "github.com/Shopify/sarama"
)

// MessageReader - The kafka message reader interface
type MessageReader interface {
	Consume(ctx context.Context) (kafka.Message, error)

	CommitMessages(ctx context.Context, msgs ...kafka.Message) error

	// Close shuts down the consumer. It must be called after all child
	// PartitionConsumers have already been closed.
	Close() error
}

// MessageWriter - The kafka message writer interface
type MessageWriter interface {
	SendMessage(msg kafka.Message) error

	// Close shuts down the producer; you must call this function before a producer
	// object passes out of scope, as it may otherwise leak memory.
	// You must call this before calling Close on the underlying client.
	Close() error
}

type kafkaWriter struct {
	kafka.SyncProducer
	topics []string
}

func newKafkaWriter(brokers []string, config *kafka.Config, topics []string) *kafkaWriter {
	producer, err := kafka.NewSyncProducer(brokers, config)
	if err != nil {
		panic(err)
	}
	return &kafkaWriter{
		SyncProducer: producer,
		topics:       topics,
	}
}

// SendMessage - send message
func (k *kafkaWriter) SendMessage(msg kafka.Message) error {
	for i := 0; i < len(k.topics); i++ {
		_, _, err := k.SyncProducer.SendMessage(&kafka.ProducerMessage{
			Topic:     k.topics[i],
			Key:       kafka.ByteEncoder(msg.Key),
			Value:     kafka.ByteEncoder(msg.Value),
			Offset:    0,
			Partition: 0,
			Timestamp: msg.Timestamp,
		})

		if err != nil {
			return err
		}
	}
	return nil
}

type kafkaReader struct {
	kafka.ConsumerGroup
	topics []string
}

func newKafkaReader(brokers []string, groupID string, config *kafka.Config, topics []string) *kafkaReader {
	consumer, err := kafka.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		panic(err)
	}
	return &kafkaReader{
		ConsumerGroup: consumer,
		topics:        topics,
	}
}

// CommitMessages - commit message
func (k *kafkaReader) CommitMessages(ctx context.Context, msgs ...kafka.Message) error {
	return nil
}

// Consume - consumer
func (k *kafkaReader) Consume(ctx context.Context) (kafka.Message, error) {
	consumer := Consumer{
		message: make(chan kafka.Message),
	}

	go func() {
		if err := k.ConsumerGroup.Consume(ctx, k.topics, &consumer); err != nil {
			fmt.Printf("consumer failure %v\n", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return kafka.Message{}, ctx.Err()
		case msg := <-consumer.message:
			return msg, nil
		}
	}
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	message chan kafka.Message
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(kafka.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(kafka.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session kafka.ConsumerGroupSession, claim kafka.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		session.MarkMessage(message, "")
		consumer.message <- kafka.Message{
			Key:       message.Key,
			Value:     message.Value,
			Timestamp: message.Timestamp,
		}

		session.Commit()
		break
	}

	return nil
}

// KafkaBroker - kafka broker
type KafkaBroker struct {
	common.Broker
	reader        MessageReader
	writer        MessageWriter
	delayedReader MessageReader
	delayedWriter MessageWriter

	consumePeriod  time.Duration
	consumeTimeout time.Duration

	consumingWG       sync.WaitGroup // wait group to make sure whole consumption completes
	processingWG      sync.WaitGroup // use wait group to make sure task processing completes
	delayedWG         sync.WaitGroup
	enableDelayedTask bool
}

type messageInfo struct {
	message kafka.Message
	reader  MessageReader
}

func New(cnf *config.Config) *KafkaBroker {
	brokers := strings.Split(strings.TrimLeft(cnf.Broker, "kafka://"), ",")

	readerCfg := func() *kafka.Config {
		conf := kafka.NewConfig()
		conf.ClientID = cnf.Kafka.ClientID
		return conf
	}

	writerCfg := func() *kafka.Config {
		conf := kafka.NewConfig()
		conf.Producer.RequiredAcks = kafka.WaitForAll
		conf.Producer.Return.Successes = true
		conf.Producer.Return.Errors = true
		conf.ClientID = cnf.Kafka.ClientID
		return conf
	}

	consumePeriod := 500 // default poll period for delayed tasks
	if cnf.Kafka != nil {
		configuredConsumePeriod := cnf.Kafka.DelayedTasksConsumePeriod
		if configuredConsumePeriod > 0 {
			consumePeriod = configuredConsumePeriod
		}
	}

	var delayedReader MessageReader
	var delayedWriter MessageWriter
	if cnf.Kafka != nil && cnf.Kafka.EnableDelayedTasks {
		delayedReader = newKafkaReader(brokers, cnf.Kafka.ConsumerGroupId, readerCfg(), []string{cnf.Kafka.DelayedTasksTopic})
		delayedWriter = newKafkaWriter(brokers, writerCfg(), []string{cnf.Kafka.DelayedTasksTopic})
	}

	// topic, delayedTasksTopic := cnf.Kafka.Topic, cnf.Kafka.DelayedTasksTopic
	return &KafkaBroker{
		Broker:            common.NewBroker(cnf),
		reader:            newKafkaReader(brokers, cnf.Kafka.ConsumerGroupId, readerCfg(), []string{cnf.Kafka.Topic}),
		writer:            newKafkaWriter(brokers, writerCfg(), []string{cnf.Kafka.Topic}),
		delayedReader:     delayedReader,
		delayedWriter:     delayedWriter,
		consumePeriod:     time.Duration(consumePeriod) * time.Millisecond,
		consumeTimeout:    time.Second * 30,
		enableDelayedTask: cnf.Kafka.EnableDelayedTasks,
	}
}

func (b *KafkaBroker) StartConsuming(consumerTag string, concurrency int, taskProcessor iface.TaskProcessor) (bool, error) {
	b.consumingWG.Add(1)
	defer b.consumingWG.Done()

	if concurrency < 1 {
		concurrency = runtime.NumCPU() * 2
	}

	b.Broker.StartConsuming(consumerTag, concurrency, taskProcessor)

	// Channel to which we will push tasks ready for processing by worker
	deliveries := make(chan messageInfo, concurrency)

	errorsChan := make(chan error, 1)

	go func() {
		log.INFO.Print("[*] Waiting for messages. To exit press CTRL+C")

		for {
			select {
			// A way to stop this goroutine from b.StopConsuming
			case <-b.GetStopChan():
				close(deliveries)
				return
			default:
				ctx, cancelFunc := context.WithTimeout(context.Background(), b.consumeTimeout)
				defer cancelFunc()
				m, err := b.reader.Consume(ctx)

				// timeout error, then retry
				if errors.Is(err, context.DeadlineExceeded) {
					continue
				}
				if err != nil {
					errorsChan <- err
					return
				}
				deliveries <- messageInfo{message: m, reader: b.reader}
			}
		}
	}()

	if b.enableDelayedTask {
		b.delayedWG.Add(1)
		go func() {
			defer b.delayedWG.Done()
			for {
				select {
				// A way to stop this goroutine from b.StopConsuming
				case <-b.GetStopChan():
					return
				default:
					err := b.processDelayedTask()
					if err != nil {
						errorsChan <- err
						return
					}
				}
			}
		}()
	}

	if err := b.consume(deliveries, concurrency, taskProcessor, errorsChan); err != nil {
		return b.GetRetry(), err
	}

	b.processingWG.Wait()

	return b.GetRetry(), nil
}

func (b *KafkaBroker) consume(deliveries <-chan messageInfo, concurrency int, taskProcessor iface.TaskProcessor, errorsChan chan error) error {
	pool := make(chan struct{}, concurrency)

	// init pool for Worker tasks execution, as many slots as Worker concurrency param
	for i := 0; i < concurrency; i++ {
		pool <- struct{}{}
	}

	for {
		select {
		case err := <-errorsChan:
			return err
		case d, open := <-deliveries:
			if !open {
				return nil
			}
			if concurrency > 0 {
				// get execution slot from pool (blocks until one is available)
				<-pool
			}

			b.processingWG.Add(1)

			// Consume the task inside a goroutine so multiple tasks
			// can be processed concurrently
			go func() {
				if err := b.consumeOne(d.reader, d.message, taskProcessor); err != nil {
					errorsChan <- err
				}

				b.processingWG.Done()

				if concurrency > 0 {
					// give slot back to pool
					pool <- struct{}{}
				}
			}()
		}
	}
}

func (b *KafkaBroker) consumeOne(reader MessageReader, message kafka.Message, taskProcessor iface.TaskProcessor) error {
	defer reader.CommitMessages(context.Background(), message)

	// Unmarshal message body into signature struct
	signature := new(tasks.Signature)
	decoder := json.NewDecoder(bytes.NewReader(message.Value))
	decoder.UseNumber()
	if err := decoder.Decode(signature); err != nil {
		return errs.NewErrCouldNotUnmarshalTaskSignature(message.Value, err)
	}

	// If the task is not registered, we nack it and requeue,
	// there might be different workers for processing specific tasks
	if !b.IsTaskRegistered(signature.Name) {
		log.INFO.Printf("Task not registered with this worker. Requeing message: %s", message.Value)
		return nil
	}
	log.DEBUG.Printf("Received new message: %s", message.Value)
	err := taskProcessor.Process(signature)
	return err
}

func (b *KafkaBroker) processDelayedTask() error {

	time.Sleep(b.consumePeriod)

	ctx, cancelFunc := context.WithTimeout(context.Background(), b.consumeTimeout)
	defer cancelFunc()
	m, err := b.delayedReader.Consume(ctx)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil
		}
		return err
	}

	defer b.delayedReader.CommitMessages(context.Background(), m)

	signature := new(tasks.Signature)
	decoder := json.NewDecoder(bytes.NewReader(m.Value))
	decoder.UseNumber()
	if err := decoder.Decode(signature); err != nil {
		return errs.NewErrCouldNotUnmarshalTaskSignature(m.Value, err)
	}

	if err := b.Publish(context.Background(), signature); err != nil {
		return err
	}
	return nil
}

func (b *KafkaBroker) Publish(ctx context.Context, signature *tasks.Signature) error {
	msg, err := json.Marshal(signature)
	if err != nil {
		return fmt.Errorf("JSON marshal error: %s", err)
	}

	// Check the ETA signature field, if it is set and it is in the future,
	// delay the task
	if b.enableDelayedTask && signature.ETA != nil {
		now := time.Now().UTC()

		if signature.ETA.After(now) {
			err = b.delayedWriter.SendMessage(kafka.Message{Value: msg})
			return err
		}
	}

	err = b.writer.SendMessage(kafka.Message{Value: msg})
	return err
}

// StopConsuming quits the loop
func (b *KafkaBroker) StopConsuming() {
	b.Broker.StopConsuming()

	if b.enableDelayedTask {
		// Waiting for the delayed tasks goroutine to have stopped
		b.delayedWG.Wait()
	}
	// Waiting for consumption to finish
	b.consumingWG.Wait()

	b.reader.Close()
	b.writer.Close()

	if b.enableDelayedTask {
		b.delayedReader.Close()
		b.delayedWriter.Close()
	}
}
