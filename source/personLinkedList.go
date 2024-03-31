// COVID-19 Vaccination Appointment Application : Built by Pallavi Limaye - 06/03/2021
package main

import (
	"errors"
	"fmt"
	"sync"
)

type adminStruct struct {
	Message    []string
	Users      []string
	Deleteuser string
	ApptAdd    string
}

type Node struct {
	item person
	next *Node
}

type linkedList struct {
	head *Node
	size int
	mu   sync.Mutex
}

func (p *linkedList) get(index int) (person, error) {
	emptyItem := person{}
	if p.head == nil {
		return emptyItem, errors.New("Empty Linked list!")
	}
	if index > 0 && index <= p.size {
		currentNode := p.head
		for i := 1; i <= index-1; i++ {
			currentNode = currentNode.next
		}
		item := currentNode.item
		return item, nil

	}
	return emptyItem, errors.New("Invalid Index")
}

func (p *linkedList) addNode(name person, wglocal *sync.WaitGroup) error {
	defer wglocal.Done()

	defer func() {
		if r := recover(); r != nil {
			println("Panic:" + r.(string))
		}
	}()
	p.mu.Lock()
	{
		newNode := &Node{
			item: name,
			next: nil,
		}
		if p.head == nil {
			p.head = newNode
		} else {
			currentNode := p.head
			for currentNode.next != nil {
				currentNode = currentNode.next
			}
			currentNode.next = newNode
		}
		p.size++

	}
	p.mu.Unlock()
	return nil
}

func (p *linkedList) addAtPos(index int, name person) error {
	newNode := &Node{
		item: name,
		next: nil,
	}

	if index > 0 && index <= p.size+1 {
		if index == 1 {
			newNode.next = p.head
			p.head = newNode

		} else {

			currentNode := p.head
			var prevNode *Node
			for i := 1; i <= index-1; i++ {
				prevNode = currentNode
				currentNode = currentNode.next
			}
			newNode.next = currentNode
			prevNode.next = newNode

		}
		p.size++
		return nil
	} else {
		return errors.New("Invalid Index")
	}
}

func (p *linkedList) remove(index int) (person, error) {
	var item person
	emptyItem := person{}

	if p.head == nil {
		return emptyItem, errors.New("Empty Linked list!")
	}
	if index > 0 && index <= p.size {
		if index == 1 {
			item = p.head.item
			p.head = p.head.next
		} else {
			var currentNode *Node = p.head
			var prevNode *Node
			for i := 1; i <= index-1; i++ {
				prevNode = currentNode
				currentNode = currentNode.next

			}
			item = currentNode.item
			prevNode.next = currentNode.next
		}
	}
	p.size--
	return item, nil
}

func (p *linkedList) printAllNodes() error {
	currentNode := p.head
	if currentNode == nil {
		return nil
	}
	fmt.Printf("%+v\n", currentNode.item)
	for currentNode.next != nil {
		currentNode = currentNode.next
		fmt.Printf("%+v\n", currentNode.item)
	}
	return nil
}

func (p *linkedList) printAllUsers() adminStruct {
	adminToTemplate := adminStruct{}

	count := 1
	currentNode := p.head
	if currentNode == nil {
		adminToTemplate.Message = append(adminToTemplate.Message, fmt.Sprintf("No users found."))
		return adminToTemplate
	}
	adminToTemplate.Message = append(adminToTemplate.Message, fmt.Sprintf("\nListing all usernames:"))
	thisPerson := currentNode.item
	//adminToTemplate.Users = append(adminToTemplate.Users, fmt.Sprintf("%d.\t%s\n", count, thisPerson.Username))
	adminToTemplate.Users = append(adminToTemplate.Users, fmt.Sprintf("%s", thisPerson.Username))

	count++
	for currentNode.next != nil {
		currentNode = currentNode.next
		//fmt.Printf("%+v\n", currentNode.item)
		thisPerson := currentNode.item
		//adminToTemplate.Users = append(adminToTemplate.Users, fmt.Sprintf("%d.\t%s\n", count, thisPerson.Username))
		adminToTemplate.Users = append(adminToTemplate.Users, fmt.Sprintf("%s", thisPerson.Username))
		count++
	}
	return adminToTemplate
}

func (p *linkedList) searchUserName(username string) (person, int, error) {
	emptyItem := person{}
	index := 1
	if p.head == nil {
		return emptyItem, -1, errors.New("Empty Linked list!")
	}
	currentNode := p.head
	for i := 1; i <= p.size; i++ {
		if currentNode.item.Username != username {
			currentNode = currentNode.next
			index++
		} else {
			item := currentNode.item
			return item, index, nil
		}
	}
	return emptyItem, -1, errors.New("Invalid Username")
}

func (p *linkedList) writeAtIndex(index int, thisPerson person) error {

	if p.head == nil {
		return errors.New("Empty Linked list!")
	}
	if index > 0 && index <= p.size {
		currentNode := p.head
		for i := 1; i <= index-1; i++ {
			currentNode = currentNode.next
		}
		currentNode.item = thisPerson
		return nil

	}
	return errors.New("Invalid Index")
}

func (p *linkedList) writePersonData(currentPerson person, currentPersonIndex int) {
	err := p.writeAtIndex(currentPersonIndex, currentPerson)
	if err != nil {
		fmt.Printf("Unable to store data for %s %s\n", currentPerson.Firstname, currentPerson.Lastname)
	}
	return
}
