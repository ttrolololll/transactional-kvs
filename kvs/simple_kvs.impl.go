package kvs

type simpleKvsImpl struct {
	kvMap     map[string]string
	vCountMap map[string]uint64
}

// NewSimpleKvs creates a new instance of simple KVS
func NewSimpleKvs() IKvs {
	return &simpleKvsImpl{
		kvMap:     map[string]string{},
		vCountMap: map[string]uint64{},
	}
}

// Set creates a transaction and commit the value
func (kvs *simpleKvsImpl) Set(k, v string) {
	kvs.kvMap[k] = v
	if _, exists := kvs.vCountMap[v]; exists {
		kvs.vCountMap[v] += 1
	} else {
		kvs.vCountMap[v] = 1
	}
}

// Get retrieves a value by key
func (kvs *simpleKvsImpl) Get(k string) string {
	v, exists := kvs.kvMap[k]
	if !exists {
		return ""
	}
	return v
}

// Delete removes a value by key
func (kvs *simpleKvsImpl) Delete(k string) {
	v := kvs.Get(k)
	if v == "" {
		return
	}

	delete(kvs.kvMap, k)
	kvs.vCountMap[v] -= 1
}

// Count returns the number of value occurrences
func (kvs *simpleKvsImpl) Count(v string) uint64 {
	if _, exists := kvs.vCountMap[v]; exists {
		return kvs.vCountMap[v]
	}
	return 0
}
