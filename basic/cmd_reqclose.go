package basic

// sendReqClose 关闭子节点
// 关闭一个有父节点的子节点，需要给它的父节点发送消息来关闭这个子节点
// p 父节点
// c 子节点
func sendReqClose(p, c *Object) {
	if p == nil || c == nil {
		return
	}
	p.Send(CommandWrapper(func(p *Object) error {
		if _, ok := p.child.Load(c.ID); ok {
			p.ack++
			p.child.Delete(c.ID)
			sendClose(c)
		}
		return nil
	}))
}
