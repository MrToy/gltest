package render

type Renderer interface {
	Render()
}

type TreeNode struct {
	Object   Renderer
	Children []*TreeNode
	Parent   *TreeNode
}

func (this *TreeNode) Add(it Renderer) {
	node := &TreeNode{
		Object: it,
		Parent: this,
	}
	this.Children = append(this.Children, node)
}

func (this *TreeNode) GetAll() []*TreeNode {
	objs := []*TreeNode{}
	if this != nil {
		objs = append(objs, this)
		objs = recur(objs, this.Children)
	}
	return objs
}
func recur(objs []*TreeNode, items []*TreeNode) []*TreeNode {
	for _, item := range items {
		objs = append(objs, item)
		objs = recur(objs, item.Children)
	}
	return objs
}
