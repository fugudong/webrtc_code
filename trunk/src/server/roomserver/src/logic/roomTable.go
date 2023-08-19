package logic

import (
	Xlog "github.com/cihub/seelog"
	"sync"
	"time"
)


// 需要考虑加锁的时机
/**

房间字典的锁
1. 新建房间
2. 删除房间

房间内部锁
1. join
2. leave等

 */
type RoomTable struct {
	lock  	sync.Mutex		// s锁
	rooms   map[string]*Room		// 房间
	registerTimeout time.Duration
}


// 创建房间字典
func NewRoomTable(to time.Duration) *RoomTable {
	return &RoomTable{rooms: make(map[string]*Room), registerTimeout: to}
}

// room returns the room specified by |id|, or creates the room if it does not exist.
func (rt *RoomTable) room(id string, name string) *Room {
	rt.lock.Lock()
	defer rt.lock.Unlock()

	return rt.roomLocked(id, name)
}

// 查找房间，如果房间不存在则创建房间
func (rt *RoomTable) roomLocked(id string, name string) *Room {
	if r, ok := rt.rooms[id]; ok {
		Xlog.Debugf("Find the existing room %s", id)
		return r
	}
	rt.rooms[id] = NewRoom(rt, id, name,rt.registerTimeout)
	size := len(rt.rooms)
	Xlog.Infof("Created new room, rid = %s, total rooms = %d", id, size)

	return rt.rooms[id]
}

// remove removes the client. If the room becomes empty, it also removes the room.
func (rt *RoomTable) remove(rid string, cid string) {
	rt.lock.Lock()
	defer rt.lock.Unlock()

	rt.removeLocked(rid, cid)
}

// removeLocked removes the client without acquiring the lock. Used when the caller already acquired the lock.
func (rt *RoomTable) removeLocked(rid string, cid string) {
	if r := rt.rooms[rid]; r != nil {
		r.remove(cid)
		if r.ClientsNumber() == 0 {
			delete(rt.rooms, rid)
			Xlog.Infof("Removed room %s", rid)
		}
	}
}

func (rt *RoomTable) GetRoomLocked(rid string) *Room {
	rt.lock.Lock()
	defer rt.lock.Unlock()
	if r, ok := rt.rooms[rid]; ok {
		return r
	}
	return nil
}

// 如果房间为空则删除
func (rt *RoomTable) removeRoom(rid string) {
	rt.lock.Lock()
	defer rt.lock.Unlock()

	if r := rt.rooms[rid]; r != nil {
		num := r.ClientsNumber()
		if num == 0 {
			delete(rt.rooms, rid)
			Xlog.Infof("Removed room %s", rid)
		} else {
			Xlog.Infof("Room %s is't empty, pepole is %d", rid, num)
		}
	}
}
