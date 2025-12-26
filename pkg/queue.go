// Package pkg 提供 xingfinger 的核心功能
// 本文件实现线程安全的队列，用于管理扫描任务
package pkg

import (
	"container/list"
	"sync"
)

// Queue 线程安全的队列结构体
// 基于双向链表实现，使用互斥锁保证并发安全
type Queue struct {
	l    sync.Mutex // 互斥锁，保证并发安全
	data *list.List // 底层双向链表
}

// NewQueue 创建新的队列实例
//
// 返回：
//   - 初始化完成的队列指针
func NewQueue() *Queue {
	q := new(Queue)
	q.data = list.New()
	return q
}

// Push 将元素添加到队列头部
// 线程安全操作
//
// 参数：
//   - v: 要添加的元素（任意类型）
//
// 返回：
//   - 新添加元素的链表节点
func (q *Queue) Push(v interface{}) *list.Element {
	q.l.Lock()
	defer q.l.Unlock()
	return q.data.PushFront(v)
}

// Pop 从队列尾部取出并移除一个元素
// 实现 FIFO（先进先出）行为
// 线程安全操作
//
// 返回：
//   - 取出的元素，队列为空时返回 nil
func (q *Queue) Pop() interface{} {
	q.l.Lock()
	defer q.l.Unlock()

	// 获取尾部元素
	iter := q.data.Back()
	if nil == iter {
		return nil
	}

	// 移除并返回元素值
	v := iter.Value
	q.data.Remove(iter)
	return v
}

// Len 获取队列当前长度
// 注意：此方法不加锁，返回值可能在多线程环境下不准确
//
// 返回：
//   - 队列中的元素数量
func (q *Queue) Len() int {
	return q.data.Len()
}
