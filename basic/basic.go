package basic

const (
	RootId = iota
	TaskID
	ModuleID
)

// Root 根节点
// 主要作用，通知所有节点关闭
var Root = NewObject(RootId, "root", new(Options), nil)

func init() {
	Root.Run()
}
