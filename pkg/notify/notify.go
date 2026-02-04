package notify

import "sync"

var (
	mu          sync.RWMutex
	subscribers = make(map[int]map[chan struct{}]struct{})
)

// Subscribe returns a channel that receives a signal each time the given match is notified.
func Subscribe(matchId int) chan struct{} {
	ch := make(chan struct{}, 1)
	mu.Lock()
	if subscribers[matchId] == nil {
		subscribers[matchId] = make(map[chan struct{}]struct{})
	}
	subscribers[matchId][ch] = struct{}{}
	mu.Unlock()
	return ch
}

// Unsubscribe removes the channel from the match's subscriber set.
func Unsubscribe(matchId int, ch chan struct{}) {
	mu.Lock()
	if subs, ok := subscribers[matchId]; ok {
		delete(subs, ch)
		if len(subs) == 0 {
			delete(subscribers, matchId)
		}
	}
	mu.Unlock()
}

// Notify sends a non-blocking signal to every subscriber of the given match.
func Notify(matchId int) {
	mu.RLock()
	subs := subscribers[matchId]
	mu.RUnlock()
	for ch := range subs {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}
