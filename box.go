package set

type builtin interface {
	~string | ~int | ~uint | ~int64 | ~uint64 | ~int32 | ~uint32 | ~int16 | ~uint16 | ~int8 | ~uint8
}

type Builtin[T builtin] struct {
	item T
}

func (b *Builtin[T]) Less(o *Builtin[T]) bool {
	return b.item < o.item
}

func (b *Builtin[T]) Equal(o *Builtin[T]) bool {
	return b.item == o.item
}

func Box[T builtin](item T) *Builtin[T] {
	return &Builtin[T]{item: item}
}
