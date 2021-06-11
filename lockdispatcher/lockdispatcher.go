package lockdispatcher

import (
	"context"
)

type ResourceCallback func() error

type dispatchRequest struct {
	resource interface{}
	waitChan chan interface{}
	returns  chan interface{}
}

func newDispatchRequest(resource interface{}, returns chan interface{}) (result dispatchRequest) {
	return dispatchRequest{
		resource: resource,
		waitChan: make(chan interface{}),
		returns:  returns,
	}
}

func (r *dispatchRequest) acquire() {
	<-r.waitChan
}

func (r *dispatchRequest) release() {
	r.returns <- r.resource
}

type LockDispatcher struct {
	context   context.Context
	resources map[interface{}][]*dispatchRequest
	requests  chan dispatchRequest
	returns   chan interface{}
}

func NewLockDispatcher(ctx context.Context) (result *LockDispatcher) {
	result = &LockDispatcher{
		context:   ctx,
		resources: make(map[interface{}][]*dispatchRequest),
		requests:  make(chan dispatchRequest),
		returns:   make(chan interface{}),
	}
	go result.run()
	return
}

func (d *LockDispatcher) run() {
	for {
		select {
		case requested := <-d.requests:
			if pendingList, ok := d.resources[requested.resource]; ok {
				d.resources[requested.resource] = append(pendingList, &requested)
			} else {
				d.resources[requested.resource] = []*dispatchRequest{}
				requested.waitChan <- 0
			}
		case returned := <-d.returns:
			if pendingList, ok := d.resources[returned]; ok {
				if len(pendingList) == 0 {
					delete(d.resources, returned)
				} else {
					first := pendingList[0]
					d.resources[returned] = pendingList[1:]
					first.waitChan <- 0
				}
			} else {
				panic("never acquired resource returned")
			}
		case <-d.context.Done():
			return
		}
	}
}

func (d *LockDispatcher) Aax(resource interface{}, callback ResourceCallback) (err error) {
	request := newDispatchRequest(resource, d.returns)
	d.requests <- request

	request.acquire()
	defer request.release()

	err = callback()
	return
}
