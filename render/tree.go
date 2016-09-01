package render

type TreeGeter interface {
	GetTree() *TreeNode
}

type TreeNode struct {
	Children []TreeGeter
	Parent   TreeGeter
}

func (this *TreeNode) GetTree() *TreeNode {
	return this
}

func AddToTree(parent, child TreeGeter) bool {
	child.GetTree().Parent = parent
	tree := parent.GetTree()
	tree.Children = append(tree.Children, child)
	return true
}

func (this *TreeNode) WalkTree() chan interface{} {
	c := make(chan interface{})
	go func() {
		recur(this, c)
		close(c)
	}()
	return c
}
func recur(it TreeGeter, c chan interface{}) {
	for _, child := range it.GetTree().Children {
		c <- child
		recur(child, c)
	}
}
