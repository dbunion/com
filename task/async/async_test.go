package async

import (
	"github.com/dbunion/com/log"
	_ "github.com/dbunion/com/log/zssky"
	"github.com/dbunion/com/task"
	"github.com/dbunion/com/task/async/fun"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestTask(t *testing.T) {
	logger, err := log.NewLogger(log.TypeZsskyLog, log.Config{
		Level:         log.LevelInfo,
		HighLighting:  true,
		JSONFormatter: false,
		AlsoToStdOut:  true,
		CallerSkip:    5,
	})

	if err != nil {
		t.Fatalf("create new logger error, err:%v", err)
	}

	cfg := task.Config{
		BrokerType:    "redis",
		Broker:        "redis://192.168.64.5:6379",
		DefaultQueue:  "machinery_tasks",
		BrokerConfig:  `{"MaxIdle":10, "MaxActive":100, "IdleTimeout": 300, "Wait": true, "ReadTimeout": 15, "WriteTimeout": 15, "ConnectTimeout": 15, "NormalTasksPollPeriod": 1000, "DelayedTasksPollPeriod": 20, "DelayedTasksKey": "REDIS_DELAYED_TASKS_KEY", "Password": "test"}`,
		BackendType:   "redis",
		ResultBackend: "redis://192.168.64.5:6379",
		BackendConfig: `{"MaxIdle":10, "MaxActive":100, "IdleTimeout": 300, "Wait": true, "ReadTimeout": 15, "WriteTimeout": 15, "ConnectTimeout": 15, "NormalTasksPollPeriod": 1000, "DelayedTasksPollPeriod": 20, "DelayedTasksKey": "REDIS_DELAYED_TASKS_KEY", "Password": "test"}`,

		ResultsExpireIn: 3600000,
		FuncWraps:       map[string]task.FuncWrap{"xx": fun.NewFuncWrap(nil)},
		Logger:          logger,
		ErrorHandler: func(err error) {
			logger.Errorf("task handler error, err:%v", err)
		},
		PreTaskHandler: func(param *task.Param) {
			logger.Infof("task pre post, name:%v func:%v", param.Name, param.Fun)
		},
		PostTaskHandler: func(param *task.Param) {
			logger.Infof("task post, name:%v func:%v", param.Name, param.Fun)
		},
	}

	tsk, err := task.NewTask(task.TypeAsync, cfg)
	if err != nil {
		t.Fatalf("create new task error, err:%v", err)
	}

	worker, err := task.NewWorker(task.TypeAsyncWorker, cfg)
	if err != nil {
		t.Fatalf("create new worker error, err:%v", err)
	}

	// run worker
	go func() {
		if err := worker.Run(); err != nil {
			t.Logf("worker run error:%v", err)
		}
	}()

	time.Sleep(time.Second * 10)

	// add task
	if err := tsk.AddTask(&task.Param{
		UUID:   uuid.New().String(),
		Name:   "add",
		Fun:    "add",
		Option: task.Option{},
		Args: []task.Arg{
			{Type: "int64", Value: 1},
			{Type: "int64", Value: 1},
		},
		WaitTimeOut: 0,
	}, nil, nil); err != nil {
		t.Fatalf("add task error, err:%v", err)
	}

	// add http task
	/*
		if err := tsk.AddTask(&task.Param{
			UUID:   uuid.New().String(),
			Name:   "httpGet",
			Fun:    "httpGet",
			Option: task.Option{},
			Args: []task.Arg{
				{Type: "string", Value: "https://www.baidu.com"},
				{Type: "int64", Value: time.Second * 10},
				{Type: "int64", Value: time.Second * 10},
			},
		}, nil, nil); err != nil {
			t.Fatalf("add task error, err:%v", err)
		}

	*/

	/*
		// add ssh task
		if err := tsk.AddTask(&task.Param{
			UUID:   uuid.New().String(),
			Name:   "sshCommand1",
			Fun:    "sshCommand",
			Option: task.Option{},
			Args: []task.Arg{
				{Type: "string", Value: "127.0.0.1"},
				{Type: "int", Value: 22},
				{Type: "string", Value: "root"},
				{Type: "string", Value: "root"},
				{Type: "string", Value: "whoami"},
			},
		}, nil, nil); err != nil {
			t.Fatalf("add task error, err:%v", err)
		}

		// add ssh task
		if err := tsk.AddTask(&task.Param{
			UUID:   uuid.New().String(),
			Name:   "sshCommand2",
			Fun:    "sshCommand",
			Option: task.Option{},
			Args: []task.Arg{
				{Type: "string", Value: "127.0.0.1"},
				{Type: "int", Value: 22},
				{Type: "string", Value: "root"},
				{Type: "string", Value: "root"},
				{Type: "string", Value: "date"},
			},
		}, nil, nil); err != nil {
			t.Fatalf("add task error, err:%v", err)
		}

	*/

	if err := tsk.Run(false); err != nil {
		t.Fatalf("task run err:%v", err)
	}

	time.Sleep(time.Second * 10)
}
