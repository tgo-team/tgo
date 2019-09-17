package tgo

import "sync"

type ClientManager interface {
	// Add 添加客户端
	Add(client Client)
	// Remove 移除客户端
	Remove(client Client)
	// Get 获取客户端
	Get(clientId uint64) Client
}

type DefaultClientManager struct {
	connMap map[uint64]Client
	sync.RWMutex
}

func NewDefaultClientManager() *DefaultClientManager  {
	return &DefaultClientManager{connMap: map[uint64]Client{}}
}

func (d *DefaultClientManager) Add(client Client) {
	d.Lock()
	defer  d.Unlock()
	d.connMap[client.GetId()] = client
}

func (d *DefaultClientManager) Remove(client Client) {
	d.Lock()
	defer  d.Unlock()

	delete(d.connMap,client.GetId())
}

func (d *DefaultClientManager) Get(connId uint64) Client{
	d.RLock()
	defer d.RUnlock()

	return d.connMap[connId]
}