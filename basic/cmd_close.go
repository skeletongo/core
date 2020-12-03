package basic

// sendClose 关闭根节点
// o 根节点
func sendClose(o *Object) {
	if o == nil {
		return
	}
	o.Send(CommandWrapper(func(p *Object) error {
		if p.closing {
			return nil
		}
		p.closing = true
		p.child.Range(func(key, value interface{}) bool {
			if c, ok := value.(*Object); ok && c != nil {
				p.ack++
				sendClose(c)
			}
			return true
		})
		p.safeStop()
		return nil
	}))
}
