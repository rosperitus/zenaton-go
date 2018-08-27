package zenaton

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const (
	MODE_AT        = "AT"
	MODE_WEEK_DAY  = "WEEK_DAY"
	MODE_MONTH_DAY = "MONTH_DAY"
	MODE_TIMESTAMP = "TIMESTAMP"
)

var waitTask = &Wait{
	task: NewTask(TaskParams{
		Name:       "_Wait",
		HandleFunc: func() {},
	}),
}

type Wait struct {
	task     *Task
	event    string
	buffer   []duration
	mode     string
	timezone *time.Location
	nowFunc  func() time.Time
}

type duration struct {
	method string
	value  interface{}
}

func NewWait() *Wait {
	return waitTask
}

func (w *Wait) FakeNow(nowFunc func() time.Time) *Wait {
	w.nowFunc = nowFunc
	return w
}

func (w *Wait) WithEvent(event string) *Wait {
	w.event = event
	return w
}

func (w *Wait) Event() string {
	return w.event
}

func (w *Wait) Seconds(value int64) *Wait {
	w.push(duration{
		method: "seconds",
		value:  value,
	})

	return w
}

func (w *Wait) Minutes(value int64) *Wait {
	w.push(duration{
		method: "minutes",
		value:  value,
	})

	return w
}

func (w *Wait) Hours(value int64) *Wait {
	w.push(duration{
		method: "hours",
		value:  value,
	})

	return w
}

func (w *Wait) Days(value int64) *Wait {
	w.push(duration{
		method: "days",
		value:  value,
	})

	return w
}

func (w *Wait) Weeks(value int64) *Wait {
	w.push(duration{
		method: "weeks",
		value:  value,
	})

	return w
}

func (w *Wait) Months(value int64) *Wait {
	w.push(duration{
		method: "months",
		value:  value,
	})

	return w
}

func (w *Wait) Years(value int64) *Wait {
	w.push(duration{
		method: "years",
		value:  value,
	})

	return w
}

func (w *Wait) Timezone(timezone string) error {
	tz, err := time.LoadLocation(timezone)
	if err != nil {
		return err
	}
	w.timezone = tz
	return nil
}

func (w *Wait) Timestamp(value int64) *Wait {
	w.push(duration{"timestamp", value})
	return w
}

func (w *Wait) At(value string) *Wait {
	w.push(duration{"at", value})
	return w
}

func (w *Wait) DayOfMonth(value string) *Wait {
	w.push(duration{"dayOfMonth", value})
	return w
}

func (w *Wait) Monday(value int) *Wait {
	w.push(duration{"Monday", value})
	return w
}

func (w *Wait) Tuesday(value int) *Wait {
	w.push(duration{"Tuesday", value})
	return w
}

func (w *Wait) Wednesday(value int) *Wait {
	w.push(duration{"Wednesday", value})
	return w
}

func (w *Wait) Thursday(value int) *Wait {
	w.push(duration{"Thursday", value})
	return w
}

func (w *Wait) Friday(value int) *Wait {
	w.push(duration{"Friday", value})
	return w
}

func (w *Wait) Saturday(value int) *Wait {
	w.push(duration{"Saturday", value})
	return w
}

func (w *Wait) Sunday(value int) *Wait {
	w.push(duration{"Sunday", value})
	return w
}

func (w *Wait) push(data duration) {
	w.buffer = append(w.buffer, data)
}

func (w *Wait) initNowThen() (time.Time, time.Time) {
	// get set or current time zone

	var tz *time.Location
	if w.timezone == nil {
		tz = time.Local
	}
	n := time.Now()
	var now = time.Date(n.Year(), n.Month(), n.Day(), n.Hour(), n.Minute(), n.Second(), n.Nanosecond(), tz)
	var then = now
	return now, then
}

//todo: would be nice to make this unexported
func (w *Wait) GetTimestampOrDuration() (int64, float64, error) {

	now, then := w.initNowThen()

	w.mode = ""

	var err error
	for _, duration := range w.buffer {
		then, err = w.apply(duration.method, duration.value, now, then)
		if err != nil {
			return 0, 0, nil
		}
	}

	isTimestamp := w.mode == ""

	if isTimestamp {
		//todo: these shouldn't be 0, right? what if the time until then is actually 0?
		return then.Unix(), 0, nil
	} else {
		return 0, time.Until(then).Seconds(), nil
	}
}

func (w *Wait) apply(method string, value interface{}, now, then time.Time) (time.Time, error) {
	switch method {
	case "timestamp":
		return w._timestamp(value.(int64)), nil
	case "at":
		return w._at(value.(string), now, then)
	case "dayOfMonth":
		return w._dayOfMonth(value.(int), now, then), nil
	case "monday":
		return w._weekDay(value.(int), 1, then), nil
	case "tuesday":
		return w._weekDay(value.(int), 2, then), nil
	case "wednesday":
		return w._weekDay(value.(int), 3, then), nil
	case "thursday":
		return w._weekDay(value.(int), 4, then), nil
	case "friday":
		return w._weekDay(value.(int), 5, then), nil
	case "saturday":
		return w._weekDay(value.(int), 6, then), nil
	case "sunday":
		return w._weekDay(value.(int), 7, then), nil
	default:
		return w._applyDuration(method, value.(int64), then)
	}
}

func (w *Wait) _timestamp(timestamp int64) time.Time {
	w._setMode(MODE_TIMESTAMP)

	return time.Unix(timestamp, 0)
}

func (w *Wait) _at(t string, now, then time.Time) (time.Time, error) {
	w._setMode(MODE_AT)

	segments := strings.Split(t, ":")
	h, err := strconv.Atoi(segments[0])
	if err != nil {
		return time.Time{}, errors.New("time formatted incorrectly")
	}
	var m int
	if len(segments) > 1 {
		m, err = strconv.Atoi(segments[1])
		if err != nil {
			return time.Time{}, errors.New("time formatted incorrectly")
		}
	}
	var s int
	if len(segments) > 2 {
		s, err = strconv.Atoi(segments[2])
		if err != nil {
			return time.Time{}, errors.New("time formatted incorrectly")
		}
	}

	then = time.Date(now.Year(), now.Month(), now.Day(), h, m, s, 0, w.timezone)

	if now.After(then) {
		switch w.mode {
		case MODE_AT:
			then = then.AddDate(0, 0, 1)
			break
		case MODE_WEEK_DAY:
			then = then.AddDate(0, 0, 7)
			break
		case MODE_MONTH_DAY:
			then = then.AddDate(0, 1, 0)
			break
		default:
			return time.Time{}, errors.New("Unknown mode: " + w.mode)
		}
	}

	return then, nil
}

func (w *Wait) _dayOfMonth(day int, now, then time.Time) time.Time {
	w._setMode(MODE_MONTH_DAY)

	then = time.Date(now.Year(), now.Month(), day, now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), w.timezone)

	if now.After(then) {
		then = then.AddDate(0, 1, 0)
	}

	return then
}

func (w *Wait) _weekDay(n int, day int, then time.Time) time.Time {
	w._setMode(MODE_WEEK_DAY)

	d := int(then.Weekday())
	then = then.AddDate(0, 0, day-d)

	if d > day {
		then.AddDate(0, 0, n*7)
	} else {
		then.AddDate(0, 0, (n-1)*7)
	}

	return then
}

func (w *Wait) _setMode(mode string) error {
	// can not apply twice the same method
	if mode == w.mode {
		return errors.New("incompatible definition in Wait methods")
	}
	// timestamp can only be used alone
	if (w.mode != "" && mode == MODE_TIMESTAMP) || w.mode == MODE_TIMESTAMP {
		return errors.New("incompatible definition in Wait methods")
	}

	// other mode takes precedence to MODE_AT
	if w.mode == "" || MODE_AT == w.mode {
		w.mode = mode
	}

	return nil
}

func (w *Wait) _applyDuration(method string, value int64, then time.Time) (time.Time, error) {
	switch method {
	case "seconds":
		return then.Add(time.Duration(value) * time.Second), nil
	case "minutes":
		return then.Add(time.Duration(value) * time.Minute), nil
	case "hours":
		return then.Add(time.Duration(value) * time.Hour), nil
	case "days":
		return then.AddDate(0, 0, int(value)), nil
	case "weeks":
		return then.AddDate(0, 0, int(value)*7), nil
	case "months":
		return then.AddDate(0, int(value), 0), nil
	case "years":
		return then.AddDate(int(value), 0, 0), nil
	default:
		return time.Time{}, errors.New("Unknown method " + method)
	}
}

func (w *Wait) Handle() (interface{}, error) {
	return w.task.Handle()
}

func (w *Wait) GetName() string {
	return w.task.GetName()
}

func (w *Wait) GetData() interface{} {
	return w.task.GetData()
}

func (w *Wait) Execute() (interface{}, error) {
	return w.task.Execute()
}
