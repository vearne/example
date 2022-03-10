package main

import (
	"container/list"
	"fmt"
)

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func inorder(root *TreeNode) []int {
	stack := list.New()
	stack.PushBack(root)

	result := make([]int, 0)

	for stack.Len() > 0 {
		b := stack.Back()
		stack.Remove(b)
		if b.Value == nil {
			b = stack.Back()
			stack.Remove(b)
			node := b.Value.(*TreeNode)
			result = append(result, node.Val)

		} else {
			x := b.Value.(*TreeNode)

			if x.Right != nil {
				stack.PushBack(x.Right)
			}

			stack.PushBack(x)
			stack.PushBack(nil)

			if x.Left != nil {
				stack.PushBack(x.Left)
			}
		}
	}

	return result
}

func main() {
	nodeMap := make(map[int]*TreeNode)
	nodeMap[5] = &TreeNode{Val: 5}
	nodeMap[4] = &TreeNode{Val: 4}
	nodeMap[2] = &TreeNode{Val: 2, Left: nodeMap[4], Right: nodeMap[5]}
	nodeMap[3] = &TreeNode{Val: 3}
	nodeMap[1] = &TreeNode{Val: 1, Left: nodeMap[2], Right: nodeMap[3]}

	fmt.Println(inorder(nodeMap[1]))
}
