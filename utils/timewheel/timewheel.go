package timewheel

import (
	"container/list"
	"sync"
	"time"
)

type Task struct {
	key     string
	delay   time.Duration
	circle  int
	data    interface{}
	element *list.Element
}

type TimeWheel struct {
	tick     time.Duration
	slots    []list.List
	slotNum  int
	current  int
	mutex    sync.Mutex
	ticker   *time.Ticker
	taskMap  map[string]*Task
	callback func(data interface{})
	stopChan chan struct{}
}

func New(tick time.Duration, slotNum int, callback func(data interface{})) *TimeWheel {
	tw := &TimeWheel{
		tick:     tick,
		slots:    make([]list.List, slotNum),
		slotNum:  slotNum,
		taskMap:  make(map[string]*Task),
		callback: callback,
		stopChan: make(chan struct{}),
	}
	tw.Start()
	return tw
}

func (tw *TimeWheel) Start() {
	tw.ticker = time.NewTicker(tw.tick)
	go func() {
		for {
			select {
			case <-tw.ticker.C:
				tw.tickHandler()
			case <-tw.stopChan:
				tw.ticker.Stop()
				return
			}
		}
	}()
}

func (tw *TimeWheel) Stop() {
	close(tw.stopChan)
}

func (tw *TimeWheel) AddTimer(delay time.Duration, key string, data interface{}) {
	if delay < 0 {
		return
	}

	tw.mutex.Lock()
	defer tw.mutex.Unlock()

	steps := int(delay / tw.tick)
	circle := steps / tw.slotNum
	slot := (tw.current + steps) % tw.slotNum

	task := &Task{
		key:    key,
		delay:  delay,
		circle: circle,
		data:   data,
	}

	e := tw.slots[slot].PushBack(task)
	task.element = e
	tw.taskMap[key] = task
}

func (tw *TimeWheel) RemoveTimer(key string) {
	tw.mutex.Lock()
	defer tw.mutex.Unlock()
	if task, ok := tw.taskMap[key]; ok {
		tw.slots[(tw.current+int(task.delay/tw.tick))%tw.slotNum].Remove(task.element)
		delete(tw.taskMap, key)
	}
}

func (tw *TimeWheel) tickHandler() {
	tw.mutex.Lock()
	defer tw.mutex.Unlock()

	slotList := &tw.slots[tw.current]
	var next *list.Element

	for e := slotList.Front(); e != nil; e = next {
		next = e.Next()
		task := e.Value.(*Task)
		if task.circle > 0 {
			task.circle--
			continue
		}
		go tw.callback(task.data)
		slotList.Remove(e)
		delete(tw.taskMap, task.key)
	}

	tw.current = (tw.current + 1) % tw.slotNum
}
