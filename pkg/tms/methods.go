package tms

func AddMessage(m string) {
	if node, ok := ctxMgr.GetValue(testResultKey); ok {
		tr := node.(*testResult)
		tr.message = m
	}
}

func AddLinks(l Link) {
	if node, ok := ctxMgr.GetValue(testResultKey); ok {
		tr := node.(*testResult)
		tr.resultLinks = append(tr.resultLinks, l)
	}
}