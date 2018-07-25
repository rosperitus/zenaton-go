package zenaton

var taskManagerInstance *TaskManager

type TaskManager struct {
	tasks map[string]*Task
}

func NewTaskManager() *TaskManager {
	if taskManagerInstance == nil {
		taskManagerInstance = &TaskManager{
			tasks: make(map[string]*Task),
		}
	}
	return taskManagerInstance
}

func (tm *TaskManager) setClass(name string, task *Task) {
	// check that this task does not exist yet
	if tm.getClass(name) != nil {
		panic(`"` + name + `" task can not be defined twice`)
	}

	tm.tasks[name] = task
}

func (tm *TaskManager) getClass(name string) *Task {
	return tm.tasks[name]
}

func (tm *TaskManager) GetTask(name, encodedData string) *Task {
	// get task class
	task := tm.getClass(name)
	// unserialize data
	err := Serializer{}.Decode(encodedData, &task.Data)
	if err != nil {
		panic(err)
	}

	//todo: what is this:?
	//// do not use construct function to set data
	//taskClass._useInit = false
	//// get new task instance
	//const task = new taskClass(data)
	//// avoid side effect
	//taskClass._useInit = true
	//// return task

	return task
}
