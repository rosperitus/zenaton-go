package zenaton

import (
	"fmt"
	"reflect"
)

type Task struct {
	name string
	//todo: would be nice if the handle func could take many arguments, instead of just one. would have to think how that would be done (maybe pass in argments into execute?)
	handleFunc interface{}
	data       interface{}
	id         func() string
}

type TaskParams struct {
	Name       string
	HandleFunc interface{}
	Data       interface{}
	ID         func() string
}

func NewTask(params TaskParams) *Task {
	validateTaskParams(params)

	task := &Task{
		name:       params.Name,
		handleFunc: params.HandleFunc,
		data:       params.Data,
		id:         params.ID,
	}

	NewTaskManager().setClass(params.Name, task)
	return task
}

//todo: should I panic here?
//todo: should I really require that you return an error?
func validateTaskParams(params TaskParams) {
	if params.Name == "" {
		panic("must set a Name for the task")
	}
	if params.HandleFunc == nil {
		panic("must set a HandleFunc for the task")
	}

	fnType := reflect.TypeOf(params.HandleFunc)
	// check that the HandleFunc is in fact a function
	if fnType.Kind() != reflect.Func {
		panic(fmt.Sprintf("HandlerFunc must be a function, not a : %s", fnType.Kind().String()))
	}
	// check that the number inputs to the HandleFunc is not greater than 1
	if fnType.NumIn() > 1 {
		panic("HandlerFunc must take a maximum of 1 argument")
	}

	// check that the number of outputs is either 0, 1 or 2 (for either (result, error), or just error, or no return)
	if fnType.NumOut() > 2 {
		panic(fmt.Sprintf("HandlerFunc must return (result, error) or just error, but found %d return values", fnType.NumOut()))
	}

	// check that the return type is valid (channels, functions, and unsafe pointers cannot be serialized)
	if fnType.NumOut() > 1 && !isValidResultType(fnType.Out(0)) {
		panic(fmt.Sprintf("HandlerFunc's first return value cannot be a channel, function, or unsafe pointer; found: %v", fnType.Out(0).Kind()))
	}

	// check that the last return value is an error
	if fnType.NumOut() > 0 && !isError(fnType.Out(fnType.NumOut()-1)) {
		panic(fmt.Sprintf("expected function second return value to return error but found %v", fnType.Out(fnType.NumOut()-1).Kind()))
	}

	// if Data is defined, the type of Data must be the same as the type of the receiver for HandleFunc
	if params.Data != nil {
		t := reflect.TypeOf(params.HandleFunc)
		if t.NumIn() != 1 {
			panic("if you specify a data field for a task, your handler function must have a receiver to accept that data" + t.Kind().String())
		}
		if t.In(0) != reflect.TypeOf(params.Data) {
			panic("type of data must be the same as the parameter type of the HandlerFunc. handlerFunction type: " +
				t.String() + " Data type: " + reflect.TypeOf(params.Data).String())
		}
	}

}

func isError(inType reflect.Type) bool {
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()
	return inType.Implements(errorInterface)
}

//todo: I'm not sure if I actually need this. I don't think the output's are actually serialized.
func isValidResultType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Func, reflect.Chan, reflect.UnsafePointer:
		return false
	}

	return true
}

func (t *Task) Handle() (interface{}, error) {

	handlFuncValue := reflect.ValueOf(t.handleFunc)
	handlFuncType := reflect.TypeOf(t.handleFunc)

	var in []reflect.Value
	if handlFuncType.NumIn() > 0 {
		in = []reflect.Value{reflect.ValueOf(t.data)}
	}

	values := handlFuncValue.Call(in)

	var err error

	if len(values) == 0 {
		return nil, nil
	}

	if !values[len(values)-1].IsNil() {
		err = values[len(values)-1].Interface().(error)
	}

	if len(values) == 1 {
		return nil, err
	}

	return values[0].Interface(), err
}

//todo: would be great if we could take a pointer to execute and modify that like json.unmarshal does, but it's hard to figure out how they do it
func (t *Task) Execute() (interface{}, error) {
	e := NewEngine()
	output, err := e.Execute([]Job{t})
	//todo: make sure this is impossible to get index out of bounds
	if output == nil {
		return nil, err
	}
	return output[0], err
}

func (t *Task) Dispatch() error {
	e := NewEngine()
	err := e.Dispatch([]Job{t})
	return err
}

func (ts Tasks) Dispatch() error {
	e := NewEngine()
	var jobs []Job
	for _, task := range ts {
		jobs = append(jobs, task)
	}
	return e.Dispatch(jobs)
}

type Tasks []*Task

func (ts Tasks) Execute() ([]interface{}, error) {
	e := NewEngine()
	var jobs []Job
	for _, task := range ts {
		jobs = append(jobs, task)
	}
	return e.Execute(jobs)
}

func (t *Task) GetName() string {
	return t.name
}

func (t *Task) GetData() interface{} {
	return t.data
}
