package mymock

type People interface {
	Say(s string) string
}

// 业务逻辑
func BizLogic(p People, s string) string {
	return "People Say:" + p.Say(s)
}
