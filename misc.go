package main

type Index[TKey comparable, TVal comparable] struct {
	kv map[TKey]TVal
	vk map[TVal]TKey
}

func NewIndex[TKey comparable, TVal comparable]() *Index[TKey, TVal] {
	return &Index[TKey, TVal]{
		kv: map[TKey]TVal{},
		vk: map[TVal]TKey{},
	}
}

func (i *Index[TKey, TVal]) Add(key TKey, value TVal) *Index[TKey, TVal] {
	i.kv[key] = value
	i.vk[value] = key
	return i
}

func (i *Index[TKey, TVal]) GetValue(key TKey) (TVal, bool) {
	res, ok := i.kv[key]
	return res, ok
}

func (i *Index[TKey, TVal]) GetKey(value TVal) (TKey, bool) {
	res, ok := i.vk[value]
	return res, ok
}
