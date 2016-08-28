package renderer

type Object interface {
	Render()
}

type TreeNode struct {
	Object
	Children []*TreeNode
	Parent   *TreeNode
}

func (this *TreeNode) Add(object Object) {
	node := &TreeNode{
		Object: object,
		Parent: this,
	}
	this.Children = append(this.Children, node)
}
