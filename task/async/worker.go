package async

import (
	"encoding/json"
	"fmt"
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/tasks"
	loginter "github.com/dbunion/com/log"
	"github.com/dbunion/com/task"
)

// Worker - worker type define
type Worker struct {
	server *machinery.Server
	worker *machinery.Worker
	logger loginter.Logger
	wraps  map[string]task.FuncWrap
}

// NewWorker create new Task worker with default collection name.
func NewWorker() task.Worker {
	return &Worker{
		wraps: map[string]task.FuncWrap{},
	}
}

// Run - run worker
func (w *Worker) Run() error {
	return w.worker.Launch()
}

// registerFuncWrap - register new external task FuncWrap implementations
func (w *Worker) registerFuncWrap(name string, wrap task.FuncWrap) error {
	if wrap == nil {
		return fmt.Errorf("invalid mainter")
	}

	if _, found := w.wraps[name]; found {
		return fmt.Errorf("wrap already register")
	}

	w.wraps[name] = wrap
	return w.server.RegisterTasks(wrap.GetTasks())
}

// Close - close worker
func (w *Worker) Close() error {
	w.worker.Quit()
	return nil
}

// StartAndGC start task worker adapter.
func (w *Worker) StartAndGC(cfg task.Config) error {
	cnf := config.Config{
		Broker:          cfg.Broker,
		DefaultQueue:    cfg.DefaultQueue,
		ResultBackend:   cfg.ResultBackend,
		ResultsExpireIn: cfg.ResultsExpireIn,
		TLSConfig:       nil,
	}

	if cfg.BrokerType == "amqp" {
		var broker config.AMQPConfig
		if err := json.Unmarshal([]byte(cfg.BrokerConfig), &broker); err != nil {
			return err
		}
		cnf.AMQP = &broker
	} else if cfg.BrokerType == "sqs" {
		var broker config.SQSConfig
		if err := json.Unmarshal([]byte(cfg.BrokerConfig), &broker); err != nil {
			return err
		}
		cnf.SQS = &broker
	} else if cfg.BrokerType == "redis" {
		var broker config.RedisConfig
		if err := json.Unmarshal([]byte(cfg.BrokerConfig), &broker); err != nil {
			return err
		}
		cnf.Redis = &broker
	} else if cfg.BrokerType == "dynamodb" {
		var broker config.DynamoDBConfig
		if err := json.Unmarshal([]byte(cfg.BrokerConfig), &broker); err != nil {
			return err
		}
		cnf.DynamoDB = &broker
	} else if cfg.BrokerType == "kafka" {
		var broker config.KafkaConfig
		if err := json.Unmarshal([]byte(cfg.BrokerConfig), &broker); err != nil {
			return err
		}
		cnf.Kafka = &broker
	}

	if cfg.BackendType == "redis" {
		var backend config.RedisConfig
		if err := json.Unmarshal([]byte(cfg.BackendConfig), &backend); err != nil {
			return err
		}
		cnf.Redis = &backend
	} else if cfg.BackendType == "mongodb" {
		var backend config.MongoDBConfig
		if err := json.Unmarshal([]byte(cfg.BackendConfig), &backend); err != nil {
			return err
		}
		cnf.MongoDB = &backend
	}

	server, err := machinery.NewServer(&cnf)
	if err != nil {
		panic(err)
	}

	consumerTag := "machinery_worker"
	concurrency := cfg.Concurrency
	if concurrency == 0 {
		concurrency = 10
	}

	w.server = server
	w.worker = server.NewWorker(consumerTag, concurrency)
	if cfg.Logger != nil {
		w.logger = cfg.Logger
		log.Set(cfg.Logger)
	}

	// register func wraps and tasks
	for key, value := range cfg.FuncWraps {
		if err := w.registerFuncWrap(key, value); err != nil {
			return err
		}
	}

	if cfg.ErrorHandler != nil {
		w.worker.SetErrorHandler(cfg.ErrorHandler)
	}

	if cfg.PostTaskHandler != nil {
		w.worker.SetPostTaskHandler(func(signature *tasks.Signature) {
			cfg.PostTaskHandler(toParam(signature))
		})
	}

	if cfg.PreTaskHandler != nil {
		w.worker.SetPreTaskHandler(func(signature *tasks.Signature) {
			cfg.PreTaskHandler(toParam(signature))
		})
	}

	return nil
}

func toParam(signature *tasks.Signature) *task.Param {
	return &task.Param{
		UUID: signature.UUID,
		Name: signature.Name,
		Fun:  signature.Name,
		Option: task.Option{
			ETA:          signature.ETA,
			Priority:     signature.Priority,
			Immutable:    signature.Immutable,
			RetryCount:   signature.RetryCount,
			RetryTimeout: signature.RetryTimeout,
		},
		Args:        toArgs(signature.Args),
		WaitTimeOut: 0,
	}
}

func toArgs(args []tasks.Arg) []task.Arg {
	as := make([]task.Arg, len(args))
	for i := 0; i < len(args); i++ {
		as[i] = toArg(args[i])
	}
	return as
}

func toArg(arg tasks.Arg) task.Arg {
	return task.Arg{
		Name:  arg.Name,
		Type:  arg.Type,
		Value: arg.Value,
	}
}

func init() {
	task.RegisterWorker(task.TypeAsyncWorker, NewWorker)
}
