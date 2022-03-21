package utils

import (
	"encoding/json"
	"sync"
)

type Map struct {
	sync.RWMutex
	m map[interface{}]interface{}
}

func (m *Map) init() {
	if m.m == nil {
		m.m = make(map[interface{}]interface{})
	}
}

func (m *Map) UnsafeGet(key interface{}) interface{} {
	if m.m == nil {
		return nil
	} else {
		return m.m[key]
	}
}

func (m *Map) Get(key interface{}) interface{} {
	m.RLock()
	defer m.RUnlock()
	return m.UnsafeGet(key)
}

func (m *Map) UnsafeSet(key interface{}, value interface{}) {
	m.init()
	m.m[key] = value
}

func (m *Map) Set(key interface{}, value interface{}) {
	m.Lock()
	defer m.Unlock()
	m.UnsafeSet(key, value)
}

func (m *Map) TestAndSet(key interface{}, value interface{}) interface{} {
	m.Lock()
	defer m.Unlock()

	m.init()

	if v, ok := m.m[key]; ok {
		return v
	} else {
		m.m[key] = value
		return nil
	}
}

func (m *Map) UnsafeDel(key interface{}) {
	m.init()
	delete(m.m, key)
}

func (m *Map) Del(key interface{}) {
	m.Lock()
	defer m.Unlock()
	m.UnsafeDel(key)
}

func (m *Map) UnsafeLen() int {
	if m.m == nil {
		return 0
	} else {
		return len(m.m)
	}
}

func (m *Map) Len() int {
	m.RLock()
	defer m.RUnlock()
	return m.UnsafeLen()
}

func (m *Map) UnsafeRange(f func(interface{}, interface{})) {
	if m.m == nil {
		return
	}
	for k, v := range m.m {
		f(k, v)
	}
}

func (m *Map) RLockRange(f func(interface{}, interface{})) {
	m.RLock()
	defer m.RUnlock()
	m.UnsafeRange(f)
}

func (m *Map) LockRange(f func(interface{}, interface{})) {
	m.Lock()
	defer m.Unlock()
	m.UnsafeRange(f)
}

func MapToJsonString(param map[string]interface{}) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}
func MapIntToJsonString(param map[string]int32) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}
func JsonStringToMap(str string) (tempMap map[string]int32) {
	_ = json.Unmarshal([]byte(str), &tempMap)
	return tempMap
}
func GetSwitchFromOptions(Options map[string]bool, key string) (result bool) {
	if flag, ok := Options[key]; !ok || flag {
		return true
	}
	return false
}
func SetSwitchFromOptions(options map[string]bool, key string, value bool) {
	if options == nil {
		options = make(map[string]bool, 5)
	}
	options[key] = value
}
