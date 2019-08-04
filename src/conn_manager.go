package tgo

import "sync"

type ConnManager interface {
	// Add 添加连接
	Add(conn StatefulConn)
	// Remove 移除连接
	Remove(conn StatefulConn)
	// Get 获取连接
	Get(connId string) StatefulConn
}

type DefaultConnManager struct {
	connMap map[string]StatefulConn
	sync.RWMutex
}

func NewDefaultConnManager() *DefaultConnManager  {
	return &DefaultConnManager{connMap: map[string]StatefulConn{}}
}

func (d *DefaultConnManager) Add(conn StatefulConn) {
	d.Lock()
	defer  d.Unlock()
	d.connMap[conn.GetId()] = conn
}

func (d *DefaultConnManager) Remove(conn StatefulConn) {
	d.Lock()
	defer  d.Unlock()

	delete(d.connMap,conn.GetId())
}

func (d *DefaultConnManager) Get(connId string) StatefulConn{
	d.RLock()
	defer d.RUnlock()

	return d.connMap[connId]
}