package mymap

import (
	"sync"
	"fmt"
)

type SafeMap struct {
	sync.RWMutex
	Map map[int]int
}

func TestSafe() {
	safeMap := newSafeMap(10)

	for i := 0; i < 100000; i++ {
		go safeMap.writeMap(i, i)
		go safeMap.readMap(i)

	}
	fmt.Println("success")

}

func TestunSafe() {
	Map := make(map[int]int)

	for i := 0; i < 100000; i++ {
		go writeMap(Map, i, i)
		//go readMap(Map, i)
		go deleteMap(Map,i)
	}

	fmt.Println("success")


}

func readMap(Map map[int]int, key int) int {
	return Map[key]
}

func writeMap(Map map[int]int, key int, value int) {
	Map[key] = value
}

func deleteMap(Map map[int]int, key int) {
	delete(Map,key)
}


func newSafeMap(size int) *SafeMap {
	sm := new(SafeMap)
	sm.Map = make(map[int]int)
	return sm

}

func (sm *SafeMap) readMap(key int) int {
	sm.RLock()
	value := sm.Map[key]
	sm.RUnlock()
	return value
}

func (sm *SafeMap) writeMap(key int, value int) {
	sm.Lock()
	sm.Map[key] = value
	sm.Unlock()
}

func (sm *SafeMap) deleteMap(key int) {
	sm.Lock()
	delete(sm.Map,key)
	sm.Unlock()
}

