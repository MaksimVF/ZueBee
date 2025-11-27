





package priority

import (
    "container/heap"
    "sync"
    "time"
)

type Job struct {
    Priority int
    Data     any
    AddedAt  time.Time
    Index    int
}

type JobQueue []*Job

func (q JobQueue) Len() int { return len(q) }
func (q JobQueue) Less(i, j int) bool { return q[i].Priority > q[j].Priority }
func (q JobQueue) Swap(i, j int) { q[i], q[j] = q[j], q[i]; q[i].Index = i; q[j].Index = j }

func (q *JobQueue) Push(x interface{}) {
    item := x.(*Job)
    item.Index = len(*q)
    *q = append(*q, item)
}

func (q *JobQueue) Pop() interface{} {
    old := *q
    n := len(old)
    item := old[n-1]
    *q = old[0 : n-1]
    return item
}

type PriorityQueue struct {
    mu  sync.Mutex
    q   JobQueue
}

func New() *PriorityQueue {
    return &PriorityQueue{}
}

func (p *PriorityQueue) Push(priority int, data any) {
    p.mu.Lock()
    defer p.mu.Unlock()

    heap.Push(&p.q, &Job{
        Priority: priority,
        Data:     data,
        AddedAt:  time.Now(),
    })
}

func (p *PriorityQueue) Pop() any {
    p.mu.Lock()
    defer p.mu.Unlock()

    if len(p.q) == 0 {
        return nil
    }

    item := heap.Pop(&p.q).(*Job)
    return item.Data
}





