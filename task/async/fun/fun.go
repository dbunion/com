package fun

import (
	"context"
	"fmt"
	"github.com/dbunion/com/task"
)

// defaultFuncWrap - provide default FuncWrap
type defaultFuncWrap struct {
	worker        task.Worker
	cancelFuncMap map[string]context.CancelFunc
	count         int64
}

// NewFuncWrap - create new func maintainer
func NewFuncWrap(worker task.Worker) task.FuncWrap {
	return &defaultFuncWrap{
		worker:        worker,
		cancelFuncMap: map[string]context.CancelFunc{},
	}
}

// GetTasks - return internal tasks list
func (m *defaultFuncWrap) GetTasks() map[string]interface{} {
	return map[string]interface{}{
		"httpGet":    m.HTTPGet,
		"sshCommand": m.ExecSSHCommand,
	}
}

// StopTask - stop task by uuid
func (m *defaultFuncWrap) StopTask(uuid string) error {
	cancelFunc, found := m.cancelFuncMap[uuid]
	if !found {
		return fmt.Errorf("task not found")
	}

	cancelFunc()
	return nil
}

func (m *defaultFuncWrap) wrangler(ctx context.Context, callback func(ctx context.Context) ([]string, error)) ([]string, error) {
	fmt.Printf("m:%v count:%v\n", m, m.count)
	m.count++
	param := task.ParamFunc(ctx)
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	// mapping cancel func and remove when func exit
	m.cancelFuncMap[param.UUID] = cancelFunc
	defer delete(m.cancelFuncMap, param.UUID)

	m.Add("running", 1)
	results, err := callback(cancelCtx)
	if err != nil {
		m.Add("failure", 1)
		m.Add("running", -1)
		return nil, err
	}
	m.Add("success", 1)
	m.Add("running", -1)
	return results, nil
}

func (m *defaultFuncWrap) Add(name string, delta int64) {
	if m.worker == nil {
		return
	}
}
