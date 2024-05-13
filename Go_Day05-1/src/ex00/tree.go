package tree

type Node struct {
	HasToy bool
	Left   *Node
	Right  *Node
}

func Create(toy bool) *Node {
	return &Node{
		HasToy: toy,
		Left:   nil,
		Right:  nil,
	}
}

func AreToysBalanced(root *Node) bool {
	if root == nil {
		return true
	}

	leftToys := getToysCount(root.Left)
	rightToys := getToysCount(root.Right)

	return leftToys == rightToys
}

func getToysCount(tree *Node) int {
	if tree == nil {
		return 0
	}

	count := 0
	if tree.HasToy {
		count = 1
	}

	leftCount := getToysCount(tree.Left)
	rightCount := getToysCount(tree.Right)

	return count + leftCount + rightCount
}
