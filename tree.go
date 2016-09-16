package render

type TreeNode struct {
	Data     interface{}
	Children []*TreeNode
	Parent   *TreeNode
}

func (this *TreeNode) Add(it *TreeNode) {
	it.Parent = this
	this.Children = append(this.Children, it)
}

func (this *TreeNode) Walk() chan *TreeNode {
	c := make(chan *TreeNode)
	go func() {
		recur(this, c)
		close(c)
	}()
	return c
}
func recur(it *TreeNode, c chan *TreeNode) {
	c <- it
	for _, child := range it.Children {
		recur(child, c)
	}
}
