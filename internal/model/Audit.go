package model

import (
	"container/list"
	"context"
	"time"
)

type Audit struct {
	ID     string    `json:"id"`
	Time   time.Time `json:"time"`
	Entity string    `json:"entity"`
	Action string    `json:"action"`
	Who    string    `json:"who"`
	What   string    `json:"what"`
}

type Observable struct {
	Subs *list.List
}

type Observer interface {
	Notify(ctx context.Context, audit *Audit)
}

func (o *Observable) Subscribe(x Observer) {
	o.Subs.PushBack(x)
}

func (o *Observable) Unsubscribe(x Observer) {
	for z := o.Subs.Front(); z != nil; z = z.Next() {
		if z.Value.(Observer) == x {
			o.Subs.Remove(z)
		}
	}
}

func (o *Observable) Fire(ctx context.Context, a *Audit) {
	for z := o.Subs.Front(); z != nil; z = z.Next() {
		z.Value.(Observer).Notify(ctx, a)
	}
}
