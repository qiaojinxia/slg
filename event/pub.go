package event

import (
	"errors"
	"sync"
	"sync/atomic"
)

//订阅的消息
type Subject struct {
}

//游戏的物体
type GameObject interface {
	GetID() int
	Notify(msg string)
}

type Channel struct {
	Name        string
	subscribers map[int]GameObject
	//  exitChan   chan int
	sync.RWMutex
	waitGroup    WaitGroupWrapper
	messageCount uint64
	exitFlag     int32
}

//发布消息
func (srv *Topic) PublishMessage(channelName, message string) (bool, error) {
	srv.RLock()
	ch, found := srv.Dict[channelName]
	if !found {
		srv.RUnlock()
		return false, errors.New("channelName不存在!")
	}
	srv.RUnlock()

	ch.Notify(message)
	//ch.Wait()
	return true, nil
}

//取消订阅
func (srv *Topic) Unsubscribe(client GameObject, channelName string) {
	srv.RLock()
	ch, found := srv.Dict[channelName]
	srv.RUnlock()
	if found {
		if ch.DeleteSubscribe(client) == 0 {
			//ch.Exit()
			srv.Lock()
			delete(srv.Dict, channelName)
			srv.Unlock()
		}
	}
}

type Topic struct {
	Dict map[string]*Channel //map[Channel.Name]*Channel
	sync.RWMutex
}

func NewChannel(channelName string) *Channel {
	return &Channel{
		Name: channelName,
		//  exitChan:       make(chan int),
		subscribers: make(map[int]GameObject),
	}
}

func (ch *Channel) AddSubscriber(subscriber GameObject) bool {
	ch.RLock()
	_, found := ch.subscribers[subscriber.GetID()]
	ch.RUnlock()

	ch.Lock()
	if !found {
		ch.subscribers[subscriber.GetID()] = subscriber
	}
	ch.Unlock()
	return found
}

func (ch *Channel) DeleteSubscribe(subscriber GameObject) int {
	var ret int
	//ch.ReplyMsg(
	//	fmt.Sprintf("从channel:%s 中删除client:%d ", ch.Name, client.Id))

	ch.Lock()
	delete(ch.subscribers, subscriber.GetID())
	ch.Unlock()

	ch.RLock()
	ret = len(ch.subscribers)
	ch.RUnlock()

	return ret
}

func (ch *Channel) Notify(message string) bool {
	ch.RLock()
	defer ch.RUnlock()
	for _, subscriber := range ch.subscribers {
		//	ch.ReplyMsg(
		//		fmt.Sprintf("channel:%s client:%d message:%s", ch.Name, cid, message))
		//}
		ch.waitGroup.Wrap(func() {
			subscriber.Notify(message)
		})

	}
	return true
}

func (ch *Channel) Exiting() bool {
	return atomic.LoadInt32(&ch.exitFlag) == 1
}

func (ch *Channel) Wait() {
	ch.waitGroup.Wait()
}

func (ch *Channel) Exit() {
	if !atomic.CompareAndSwapInt32(&ch.exitFlag, 0, 1) {
		return
	}
	//close(ch.exitChan)
	ch.Wait()
}
