package safelist

import (
	"fmt"
	"os"
	"sync"
)

// SafeList is a thread-safe map that can be used to store single values (for avoiding duplication) or mapped values
type SafeList struct {
	sync.RWMutex
	s map[string]struct{}
	m map[string]map[string]struct{}
}

// NewSafeList creates a thread-safe map
func New() *SafeList {
	return &SafeList{
		s: make(map[string]struct{}),
		m: make(map[string]map[string]struct{}),
	}
}

func (sm *SafeList) String() string {
	if len(sm.s) > 0 {
		return fmt.Sprintln(sm.GetValues())
	} else if len(sm.m) > 0 {
		str := ""
		for s, v := range sm.GetListValues() {
			str += fmt.Sprintf("%s: %s\n", s, v)
		}
		return str
	}
	return ""
}

// Set adds a value to the SafeList
func (sm *SafeList) Set(s string) {
	sm.Lock()
	defer sm.Unlock()

	sm.set(s)
}

func (sm *SafeList) set(s string) {
	sm.s[s] = struct{}{}
}

// SetMap adds a mapped value to the SafeList map
func (sm *SafeList) SetList(s string, v string) {
	sm.Lock()
	defer sm.Unlock()

	sm.setList(s, v)
}

func (sm *SafeList) setList(s string, v string) {
	if _, ok := sm.m[s]; !ok {
		sm.m[s] = map[string]struct{}{}
	}
	sm.m[s][v] = struct{}{}
}

// SetMultiple sets multiple values in the single value map
func (sm *SafeList) SetMultiple(slist []string) {
	sm.Lock()
	defer sm.Unlock()

	for _, s := range slist {
		sm.set(s)
	}

	return
}

// Get returns true if a value is set and false if it is not
func (sm *SafeList) Get(s string) bool {
	sm.RLock()
	defer sm.RUnlock()
	return sm.get(s)
}

func (sm *SafeList) get(s string) bool {
	_, ok := sm.s[s]
	return ok
}

// GetList retrieves a list of values from the SafeList map
func (sm *SafeList) GetList(s string) []string {
	sm.RLock()
	defer sm.RUnlock()
	return sm.getList(s)
}

func (sm *SafeList) getList(s string) []string {
	if _, ok := sm.m[s]; ok {
		v := make([]string, 0, len(sm.m[s]))
		for k := range sm.m[s] {
			v = append(v, k)
		}
		return v
	}
	return []string{}
}

// GetOrSet tries to retrieve a value from a SafeList, and if it doesn't exist, it adds it
func (sm *SafeList) GetOrSet(s string) bool {
	sm.Lock()
	defer sm.Unlock()
	if sm.get(s) {
		return true
	}
	sm.set(s)
	return false
}

// GetOrSetList tries to retrieve a value from a SafeList list, and if it doesn't exist, it adds it
func (sm *SafeList) GetOrSetList(s string, v string) bool {
	sm.Lock()
	defer sm.Unlock()
	if _, ok := sm.m[s]; !ok {
		sm.setList(s, v)
		return false
	}
	if _, ok := sm.m[s][v]; !ok {
		sm.setList(s, v)
		return false
	}
	return true
}

// GetAndSetList tries to retrieve a value from a SafeList list, and returns that value while setting it
func (sm *SafeList) GetAndSetList(s string, v string) ([]string, bool) {
	sm.Lock()
	defer sm.Unlock()

	if _, ok := sm.m[s]; !ok {
		sm.setList(s, v)
		return []string{v}, false
	}
	if _, ok := sm.m[s][v]; !ok {
		sm.setList(s, v)
		return sm.getList(s), false
	}
	return sm.getList(s), true
}

// Len returns the length of a SafeList
func (sm *SafeList) Len() int {
	sm.RLock()
	defer sm.RUnlock()
	return sm.len()
}

func (sm *SafeList) len() int {
	if len(sm.m) != 0 {
		return len(sm.m)
	}
	return len(sm.s)
}

//Delete removes a value for a key
func (sm *SafeList) Delete(s string) {
	sm.Lock()
	defer sm.Unlock()

	sm.del(s)
	return
}

func (sm *SafeList) del(s string) {
	delete(sm.s, s)
	return
}

// Delete list deletes a list for a key
func (sm *SafeList) DeleteList(s string) {
	sm.Lock()
	defer sm.Unlock()

	sm.delList(s)
	return
}

func (sm *SafeList) delList(s string) {
	delete(sm.m, s)
	return
}

func (sm *SafeList) DeleteListValue(s, v string) {
	sm.Lock()
	defer sm.Unlock()

	sm.deleteListValue(s, v)
	return
}

func (sm *SafeList) deleteListValue(s, v string) {
	if _, ok := sm.m[s]; ok {
		delete(sm.m[s], v)
	}
	return
}

// GetValues gets the values of a SafeList and return them in a slice
func (sm *SafeList) GetValues() (vals []string) {
	sm.Lock()
	defer sm.Unlock()
	for id, _ := range sm.s {
		vals = append(vals, id)
	}

	return vals
}

func (sm *SafeList) GetListValues() (vals map[string][]string) {
	sm.Lock()
	defer sm.Unlock()

	vals = make(map[string][]string)
	for s, vs := range sm.m {
		if _, ok := vals[s]; !ok {
			vals[s] = []string{}
		}
		for v, _ := range vs {
			vals[s] = append(vals[s], v)
		}
	}

	return vals
}

// PrintListToFile prints the SafeList to a file for debugging
func (sm *SafeList) PrintListToFile(filename string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	for str := range sm.s {
		_, err = f.WriteString(fmt.Sprintf("%s\n", str))
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	return nil
}
