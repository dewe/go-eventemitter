package output

type EventEmitter struct {
	listeners map[string][]chan interface{}
}

func NewEventEmitter() *EventEmitter {
	return &EventEmitter{listeners: make(map[string][]chan interface{})}
}

func (e *EventEmitter) On(event string) chan interface{} {
	ch := make(chan interface{})
	e.listeners[event] = append(e.listeners[event], ch)
	return ch
}

func (e *EventEmitter) Emit(event string, data interface{}) {
	for _, ch := range e.listeners[event] {
		go func(c chan interface{}) {
			c <- data
		}(ch)
	}
}

func (e *EventEmitter) Off(event string, ch chan interface{}) {
	if listeners, ok := e.listeners[event]; ok {
		for i, listener := range listeners {
			if listener == ch {
				e.listeners[event] = append(e.listeners[event][:i], e.listeners[event][i+1:]...)
				close(listener)
				break
			}
		}
	}
}
