package main

type MultiMap[K, V comparable] struct {
	containedMap map[K]*[]V
}

func NewMultiMap[K comparable, V comparable]() *MultiMap[K, V] {
	m := new(MultiMap[K, V])
	m.containedMap = make(map[K]*[]V)
	return m
}

func (m *MultiMap[K, V]) GetFirstValueIfKeyExists(key K) (V, bool) {
	containedSlice := m.getOrCreateSliceForKey(key)
	if containedSlice != nil && len(*containedSlice) > 0 {
		return (*containedSlice)[0], true
	} else {
		return *new(V), false
	}
}

func (m *MultiMap[K, V]) GetFirstValueForKey(key K) V {
	containedSlice := m.getOrCreateSliceForKey(key)
	if containedSlice != nil && len(*containedSlice) > 0 {
		return (*containedSlice)[0]
	} else {
		return *new(V)
	}
}

func (m *MultiMap[K, V]) ContainsKey(key K) bool {
	_, exists := m.containedMap[key]
	return exists
}

func (m *MultiMap[K, V]) AddKeyValue(key K, value V) {
	containedSlice := m.getOrCreateSliceForKey(key)
	*containedSlice = append(*containedSlice, value)
}

func (m *MultiMap[K, V]) RemoveKeyValue(key K, value V) {
	containedSlice := m.getOrCreateSliceForKey(key)

	if containedSlice != nil {
		for index, v := range *containedSlice {
			if v == value {
				*containedSlice = append((*containedSlice)[:index], (*containedSlice)[index+1:]...)
				break
			}
		}
		if len(*containedSlice) == 0 {
			m.removeListForKey(key)
		}
	}
}

func (m *MultiMap[K, V]) getOrCreateSliceForKey(key K) *[]V {
	var containedSlice *[]V
	containedSlice = m.containedMap[key]

	if containedSlice == nil {
		var slice = make([]V, 0)
		containedSlice = &slice
		m.containedMap[key] = containedSlice
	}
	return containedSlice
}

func (m *MultiMap[K, V]) removeListForKey(key K) {
	delete(m.containedMap, key)
}
