package basic

// sendAck 给父节点发送当前节点已经关闭的消息
// p 父节点
func sendAck(p *Object) {
	if p == nil {
		return
	}
	p.Send(CommandWrapper(func(p *Object) error {
		if p.ack > 0 {
			p.ack--
		}
		return nil
	}))
}
