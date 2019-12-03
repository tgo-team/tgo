package tgo

import "sync"

// ClientManager 客户端管理
type ClientManager interface {
	// Add 添加客户端
	Add(client Client)
	// Remove 移除客户端
	Remove(client Client)
	// Get 获取客户端
	Get(clientID string) Client
}

// DefaultClientManager 默认客户端管理者
type DefaultClientManager struct {
	connMap map[string]Client
	sync.RWMutex
}

// NewDefaultClientManager NewDefaultClientManager
func NewDefaultClientManager() *DefaultClientManager {
	return &DefaultClientManager{connMap: map[string]Client{}}
}

// Add 添加客户端
func (d *DefaultClientManager) Add(client Client) {
	d.Lock()
	defer d.Unlock()
	d.connMap[client.GetID()] = client
}

// Remove 移除一个客户端
func (d *DefaultClientManager) Remove(client Client) {
	d.Lock()
	defer d.Unlock()

	delete(d.connMap, client.GetID())
}

// Get 获取指定id的客户端
func (d *DefaultClientManager) Get(id string) Client {
	d.RLock()
	defer d.RUnlock()

	return d.connMap[id]
}
