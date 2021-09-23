package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len       int
	firstItem *ListItem
	lastItem  *ListItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.firstItem
}

func (l *list) Back() *ListItem {
	return l.lastItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	i := &ListItem{}
	i.Value = v
	frontItem := l.firstItem
	l.len++
	if frontItem == nil /*Empty List*/ {
		i.Prev = nil
		i.Next = nil
		l.firstItem = i
		l.lastItem = i
		return i
	}
	i.Prev = nil
	i.Next = frontItem
	frontItem.Prev = i
	l.firstItem = i
	return i
}

func (l *list) PushBack(v interface{}) *ListItem {
	i := &ListItem{}
	i.Value = v
	lastItem := l.lastItem
	l.len++
	if lastItem == nil /*Empty List*/ {
		i.Prev = nil
		i.Next = nil
		l.firstItem = i
		l.lastItem = i
		return i
	}
	lastItem.Next = i
	i.Prev = lastItem
	i.Next = nil
	l.lastItem = i
	return i
}

func (l *list) Remove(i *ListItem) {
	if i.Prev == nil {
		l.firstItem = i.Next
	} else {
		i.Prev.Next = i.Next
	}
	if i.Next == nil {
		l.lastItem = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}
