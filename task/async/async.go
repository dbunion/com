package async

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/tasks"
	loginter "github.com/dbunion/com/log"
	"github.com/dbunion/com/task"
	"reflect"
	"time"
)

var (
	sleepDuration = time.Second * 5
)

// item - task item info
type item struct {
	param     *task.Param
	callbacks []task.CallbackFunc
	onSuccess []*task.Param
	onError   []*task.Param
}

// Task is task adaptor for github.com/RichardKnop/machinery
type Task struct {
	server *machinery.Server
	items  []*item
	logger loginter.Logger
	wraps  map[string]task.FuncWrap
}

// NewTask create new Task with default collection name.
func NewTask() task.Task {
	return &Task{
		wraps: map[string]task.FuncWrap{},
	}
}

// AddTask - add new task
func (t *Task) AddTask(param *task.Param, onSuccess []*task.Param, onError []*task.Param, callbacks ...task.CallbackFunc) error {
	if param == nil {
		return fmt.Errorf("INVALID task param")
	}

	item := &item{param: param, onSuccess: onSuccess, onError: onError}
	if callbacks != nil {
		item.callbacks = make([]task.CallbackFunc, 0)
		item.callbacks = append(item.callbacks, callbacks...)
	}

	t.items = append(t.items, item)

	return nil
}

func (t *Task) callback(item *item, r task.Result) {
	var results []reflect.Value
	var err error

	if item.param.WaitTimeOut > 0 {
		results, err = r.GetWithTimeout(item.param.WaitTimeOut, sleepDuration)
	} else {
		results, err = r.Get(sleepDuration)
	}

	var message string
	if err == nil {
		message = tasks.HumanReadableResults(results)
	}

	// call all callback function
	for i := 0; i < len(item.callbacks); i++ {
		if err := item.callbacks[i](item.param, err, message); err != nil {
			fmt.Printf("task[%v]call back function call err:%v index:%v\n ", item.param.UUID, err, i)
		}
	}
}

// Run - run all task
func (t *Task) Run(chain bool) error {
	defer func() {
		t.items = make([]*item, 0)
	}()

	makeArgsFunc := func(list []task.Arg) []tasks.Arg {
		args := make([]tasks.Arg, 0)
		for _, arg := range list {
			args = append(args, tasks.Arg{
				Name:  arg.Name,
				Type:  arg.Type,
				Value: arg.Value,
			})
		}
		return args
	}

	makeSignaturesFunc := func(list []*task.Param) []*tasks.Signature {
		var signs []*tasks.Signature
		if len(list) > 0 {
			signs = make([]*tasks.Signature, len(list))
		}
		for i := 0; i < len(list); i++ {
			signs[i] = &tasks.Signature{
				UUID:         list[i].UUID,
				Name:         list[i].Fun,
				ETA:          list[i].Option.ETA,
				Args:         makeArgsFunc(list[i].Args),
				Priority:     list[i].Option.Priority,
				Immutable:    list[i].Option.Immutable,
				RetryCount:   list[i].Option.RetryCount,
				RetryTimeout: list[i].Option.RetryTimeout,
			}
		}
		return signs
	}

	length := len(t.items)

	if !chain {
		for i := 0; i < length; i++ {
			item := t.items[i]

			args := makeArgsFunc(item.param.Args)

			// construct onSuccess call task
			onSuccess := makeSignaturesFunc(item.onSuccess)

			// construct onSuccess call task
			onError := makeSignaturesFunc(item.onError)

			taskSign := &tasks.Signature{
				UUID:         item.param.UUID,
				Name:         item.param.Fun,
				ETA:          item.param.Option.ETA,
				Args:         args,
				Priority:     item.param.Option.Priority,
				Immutable:    item.param.Option.Immutable,
				RetryCount:   item.param.Option.RetryCount,
				RetryTimeout: item.param.Option.RetryTimeout,
				OnSuccess:    onSuccess,
				OnError:      onError,
			}

			asyncResult, err := t.server.SendTaskWithContext(context.Background(), taskSign)
			if err != nil {
				return err
			}

			// if callbacks define, run call back function
			if item.callbacks != nil {
				go t.callback(item, asyncResult)
			}
		}
	} else {
		list := make([]*tasks.Signature, 0)

		// use to notice last item
		var item *item
		for i := 0; i < length; i++ {
			item = t.items[i]

			args := makeArgsFunc(item.param.Args)

			// construct onSuccess call task
			onError := makeSignaturesFunc(item.onError)

			taskSign := &tasks.Signature{
				UUID:         item.param.UUID,
				Name:         item.param.Fun,
				ETA:          item.param.Option.ETA,
				Args:         args,
				Priority:     item.param.Option.Priority,
				Immutable:    item.param.Option.Immutable,
				RetryCount:   item.param.Option.RetryCount,
				RetryTimeout: item.param.Option.RetryTimeout,
				OnError:      onError,
			}

			list = append(list, taskSign)
		}

		chainAsynResult, err := t.server.SendChain(&tasks.Chain{
			Tasks: list,
		})

		if err != nil {
			return err
		}

		// if callbacks define, run call back function
		// call last item callback
		if item != nil && item.callbacks != nil {
			go t.callback(item, chainAsynResult)
		}
	}

	return nil
}

// Stop - stop all task
func (t *Task) Stop() error {
	return nil
}

// registerFuncWrap - register new external task FuncWrap implementations
func (t *Task) registerFuncWrap(name string, maintainer task.FuncWrap) error {
	if maintainer == nil {
		return fmt.Errorf("invalid mainter")
	}

	if _, found := t.wraps[name]; found {
		return fmt.Errorf("wrap already register")
	}

	t.wraps[name] = maintainer
	return t.server.RegisterTasks(maintainer.GetTasks())
}

// StartAndGC start task adapter.
func (t *Task) StartAndGC(cfg task.Config) error {
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

	if cfg.Logger != nil {
		t.logger = cfg.Logger
		log.Set(cfg.Logger)
	}

	server, err := machinery.NewServer(&cnf)
	if err != nil {
		panic(err)
	}

	t.server = server

	// register func wraps and tasks
	for key, value := range cfg.FuncWraps {
		if err := t.registerFuncWrap(key, value); err != nil {
			return err
		}
	}

	return nil
}

func paramFromContext(ctx context.Context) *task.Param {
	sign := tasks.SignatureFromContext(ctx)
	if sign == nil {
		return nil
	}

	args := make([]task.Arg, len(sign.Args))
	for i := 0; i < len(sign.Args); i++ {
		item := sign.Args[i]
		args[i] = task.Arg{
			Name:  item.Name,
			Type:  item.Type,
			Value: item.Value,
		}
	}

	return &task.Param{
		UUID: sign.UUID,
		Name: sign.Name,
		Fun:  sign.Name,
		Option: task.Option{
			ETA:          sign.ETA,
			Priority:     sign.Priority,
			Immutable:    sign.Immutable,
			RetryCount:   sign.RetryCount,
			RetryTimeout: sign.RetryTimeout,
		},
		Args: args,
	}

}

func init() {
	task.Register(task.TypeAsync, NewTask)
	task.ParamFunc = paramFromContext
}
