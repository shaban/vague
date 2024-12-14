package vague

type TemplateTreeNode struct {
	RenderFunc func(data interface{}) string
	CachedHTML string
	IsDirty    bool
}
