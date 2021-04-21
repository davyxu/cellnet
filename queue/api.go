package cellqueue

// 有队列时队列调用，无队列时直接调用
func QueuedCall(queue *Queue, callback func()) {
	if queue == nil {
		callback()
	} else {
		queue.Post(callback)
	}
}
