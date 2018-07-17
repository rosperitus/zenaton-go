package engine

import (
	"reflect"

	"github.com/davecgh/go-spew/spew"
	"github.com/zenaton/zenaton-go/v1/zenaton/client"
	"github.com/zenaton/zenaton-go/v1/zenaton/job"
)

var instance *Engine

type Engine struct {
	client    *client.Client
	processor Processor //todo: not sure what type a processor should be
}

func New() *Engine {
	if instance == nil {
		instance = &Engine{
			client: client.New(false),
		}
	}
	return instance
}

type Processor interface {
	Process([]job.Job, bool) []interface{}
}

type chanResult struct {
	result interface{}
	index  int
}

func wrapper(index int, outcome chan chanResult, handle func() interface{}) {
	outcome <- chanResult{
		index:  index,
		result: handle(),
	}
}

//todo: error handling
func (e *Engine) Execute(jobs []job.Job) []interface{} {

	// local execution
	if e.processor == nil || len(jobs) == 0 {
		outputChan := make(chan chanResult)
		for i, job := range jobs {
			go wrapper(i, outputChan, job.Handle)
		}

		outputs := make([]interface{}, len(jobs))
		for range jobs {
			o, ok := <-outputChan
			if !ok {
				break
			}
			outputs[o.index] = o.result
		}

		return outputs
	}

	//todo: figure out the processor stuff
	return e.processor.Process(jobs, true)
}

func (e *Engine) Dispatch(jobs []job.Job) []chan interface{} {
	// local execution
	var chans []chan interface{}
	for range jobs {
		chans = append(chans, make(chan interface{}))
	}
	spew.Dump(!reflect.ValueOf(e.processor).IsValid())
	if !reflect.ValueOf(e.processor).IsValid() || len(jobs) == 0 {
		for i, job := range jobs {
			go job.AsyncHandle(chans[i])
		}
	}
	return chans
	//todo: (figure out what to do with processor here) return e.processor.Process(jobs, false)
}

func (e *Engine) SetProcessor(processor Processor) {
	e.processor = processor
}
