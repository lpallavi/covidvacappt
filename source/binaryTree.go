// COVID-19 Vaccination Appointment Application : Built by Pallavi Limaye - 06/03/2021
package main

import (
	"fmt"
	"sync"
)

type BinaryNode struct {
	item  string      // to store the data item
	left  *BinaryNode // pointer to point to left node
	right *BinaryNode // pointer to point to right node
}

type BST struct {
	root *BinaryNode
	mu   sync.Mutex
}

func (bst *BST) reset() {
	bst.root = nil
}

func (bst *BST) search(item string) *BinaryNode {
	return bst.searchNode(bst.root, item)
}

func (bst *BST) searchNode(t *BinaryNode, item string) *BinaryNode {
	if t == nil {
		return nil
	} else {
		if t.item == item {
			return t
		} else {
			if item < t.item {
				return bst.searchNode(t.left, item)
			} else {
				return bst.searchNode(t.right, item)
			}
		}
	}
}

func (bst *BST) insertNode(t **BinaryNode, item string) error {

	if *t == nil {
		// that is if the value inside the memory pointed to by t which is address of bst.root
		// &bst.root is nil, then there are no nodes to the tree
		newNode := &BinaryNode{
			item:  item,
			left:  nil,
			right: nil,
		}
		*t = newNode
		return nil
	}

	if item < (*t).item {
		bst.insertNode(&((*t).left), item)
	} else {
		bst.insertNode(&((*t).right), item)
	}

	return nil
}

func (bst *BST) insert(item string, wglocal *sync.WaitGroup) {
	defer wglocal.Done()
	defer func() {
		if r := recover(); r != nil {
			println("Panic:" + r.(string))
		}
	}()
	bst.mu.Lock()
	{
		bst.insertNode(&bst.root, item)
	}
	bst.mu.Unlock()
	return
}

func (bst *BST) inOrderTraverse(t *BinaryNode) {
	if t != nil {
		bst.inOrderTraverse(t.left)
		fmt.Println(t.item)
		bst.inOrderTraverse(t.right)
	}
}

func (bst *BST) inOrder() {
	bst.inOrderTraverse(bst.root)
}

func (bst *BST) delete(item string, wglocal *sync.WaitGroup) {
	defer wglocal.Done()
	defer func() {
		if r := recover(); r != nil {
			println("Panic:" + r.(string))
		}
	}()
	bst.mu.Lock()
	{
		bst.deleteNode(bst.root, item)
	}
	bst.mu.Unlock()
	return
}

// internal recursive function to delete an item
func (bst *BST) deleteNode(t *BinaryNode, item string) *BinaryNode {
	if t == nil {
		return nil
	}
	if item < t.item {
		t.left = bst.deleteNode(t.left, item)
		return t
	}
	if item > t.item {
		t.right = bst.deleteNode(t.right, item)
		return t
	}
	// item == t.item
	if t.left == nil && t.right == nil {
		t = nil
		return nil
	}
	if t.left == nil {
		t = t.right
		return t
	}
	if t.right == nil {
		t = t.left
		return t
	}
	leftmostrightside := t.right
	for {
		//find smallest value on the right side
		if leftmostrightside != nil && leftmostrightside.left != nil {
			leftmostrightside = leftmostrightside.left
		} else {
			break
		}
	}
	t.item = leftmostrightside.item
	t.right = bst.deleteNode(t.right, t.item)
	return t
}
