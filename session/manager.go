package session

import "web_copy"

type Manager struct {
	Store
	Propagator
	// 在上下文缓存Session	map[SessCtxKey] = Session
	SessCtxKey string
}

// InitSession 创建一个session
// 先根据id创一个session再注入到响应中
// Generate Inject
func (m *Manager) InitSession(ctx *web_copy.Context, id string) (Session, error) {
	sess, err := m.Store.Generate(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}
	if err = m.Propagator.Inject(id, ctx.Resp); err != nil {
		return nil, err
	}
	return sess, nil
}

// GetSession 获取session
// 先尝试从缓存中拿
func (m *Manager) GetSession(ctx *web_copy.Context) (Session, error) {
	if ctx.UserValues == nil {
		ctx.UserValues = make(map[string]any, 10)
	}

	// 1.先判断缓存
	sess, ok := ctx.UserValues[m.SessCtxKey]
	if ok {
		return sess.(Session), nil
	}
	// 2.在请求中获取id
	id, err := m.Extract(ctx.Req)
	if err != nil {
		return nil, err
	}
	// 3.在Store中获取一个
	sess, err = m.Get(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}
	// 缓存住Session
	ctx.UserValues[m.SessCtxKey] = sess
	return sess.(Session), nil
}

// RefreshSession 刷新session
// 先得到session再刷新过期时间然后重新注入到响应中
// GetSession	Refresh		Inject
func (m *Manager) RefreshSession(ctx *web_copy.Context) (Session, error) {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return nil, err
	}
	// 刷新存储的过期时间
	err = m.Refresh(ctx.Req.Context(), sess.ID())
	if err != nil {
		return nil, err
	}
	// 重新注入到HTTP中
	if err = m.Inject(sess.ID(), ctx.Resp); err != nil {
		return nil, err
	}
	return sess, nil
}

// RemoveSession 删除 Session
func (m *Manager) RemoveSession(ctx *web_copy.Context) error {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	err = m.Store.Remove(ctx.Req.Context(), sess.ID())
	if err != nil {
		return err
	}
	return m.Propagator.Remove(ctx.Resp)
}
