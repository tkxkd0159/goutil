package ds

type IterList[T any] struct {
	size  int
	items []T
}

func NewIterList[T any]() *IterList[T] {
	return new(IterList[T])
}

func (c *IterList[T]) Size() int {
	return c.size
}

func (c *IterList[T]) Get(idx int) T {
	return c.items[idx]
}

func (c *IterList[T]) Delete(idx int) {
	if (idx + 1) > c.size {
		panic("index out of range on deletion")
	}
	c.items = append(c.items[:idx], c.items[idx+1:]...)
	c.size -= 1
}

func (c *IterList[T]) Append(item T) {
	c.items = append(c.items, item)
	c.size += 1
}

func (c *IterList[T]) Iter() <-chan T {
	ch := make(chan T)
	go func() {
		for i := 0; i < c.size; i++ {
			ch <- c.items[i]
		}
		close(ch)
	}()
	return ch
}
