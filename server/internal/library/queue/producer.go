// Package queue
package queue

import (
	"github.com/gogf/gf/v2/util/gconv"
)

// Push 推送队列
func Push(topic string, data any) (err error) {
	q, err := InstanceProducer()
	if err != nil {
		return
	}
	mqMsg, err := q.SendMsg(topic, gconv.String(data))
	ProducerLog(ctx, topic, mqMsg, err)
	return
}

// DelayPush 推送延迟队列
// redis delay 传入 秒。如：10代表延迟10秒
// rocketmq delay 传入 延迟级别。如：2代表延迟5秒
// rocketmq reference delay level definition: 1s 5s 10s 30s 1m 2m 3m 4m 5m 6m 7m 8m 9m 10m 20m 30m 1h 2h
// rocketmq delay level starts from 1. for example, if we set param level=1, then the delay time is 1s.
func DelayPush(topic string, data any, delay int64) (err error) {
	q, err := InstanceProducer()
	if err != nil {
		return
	}
	mqMsg, err := q.SendDelayMsg(topic, gconv.String(data), delay)
	ProducerLog(ctx, topic, mqMsg, err)
	return
}
