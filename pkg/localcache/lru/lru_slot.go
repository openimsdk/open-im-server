package lru

func NewSlotLRU[K comparable, V any](slotNum int, hash func(K) uint64, create func() LRU[K, V]) LRU[K, V] {
	x := &slotLRU[K, V]{
		n:     uint64(slotNum),
		slots: make([]LRU[K, V], slotNum),
		hash:  hash,
	}
	for i := 0; i < slotNum; i++ {
		x.slots[i] = create()
	}
	return x
}

type slotLRU[K comparable, V any] struct {
	n     uint64
	slots []LRU[K, V]
	hash  func(k K) uint64
}

func (x *slotLRU[K, V]) GetBatch(keys []K, fetch func(keys []K) (map[K]V, error)) (map[K]V, error) {
	var (
		slotKeys = make(map[uint64][]K)
		kVs       = make(map[K]V)
	)

	for _, k := range keys {
		index := x.getIndex(k)
		slotKeys[index] = append(slotKeys[index], k)
	}

	for k, v := range slotKeys {
		batches, err := x.slots[k].GetBatch(v, fetch)
		if err != nil {
			return nil, err
		}
		for key, value := range batches {
			kVs[key] = value
		}
	}
	return kVs, nil
}

func (x *slotLRU[K, V]) getIndex(k K) uint64 {
	return x.hash(k) % x.n
}

func (x *slotLRU[K, V]) Get(key K, fetch func() (V, error)) (V, error) {
	return x.slots[x.getIndex(key)].Get(key, fetch)
}

func (x *slotLRU[K, V]) Set(key K, value V) {
	x.slots[x.getIndex(key)].Set(key, value)
}

func (x *slotLRU[K, V]) SetHas(key K, value V) bool {
	return x.slots[x.getIndex(key)].SetHas(key, value)
}

func (x *slotLRU[K, V]) Del(key K) bool {
	return x.slots[x.getIndex(key)].Del(key)
}

func (x *slotLRU[K, V]) Stop() {
	for _, slot := range x.slots {
		slot.Stop()
	}
}
