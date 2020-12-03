package basic

// sendAddChild 将一个节点连接到另一个节点上，成为另一个节点的子节点
// p 父节点
// c 子节点
func sendAddChild(p, c *Object) {
	if p == nil || c == nil {
		return
	}
	p.Send(CommandWrapper(func(p *Object) error {
		p.child.Store(c.ID, c)
		return nil
	}))
}
