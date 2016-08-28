package render

type Renderer interface {
	Render()
}

type TreeNode struct {
	Renderer
	Children []*TreeNode
	Parent   *TreeNode
}

func (this *TreeNode) Add(it Renderer) {
	node := &TreeNode{
		Renderer: it,
		Parent:   this,
	}
	this.Children = append(this.Children, node)
}
