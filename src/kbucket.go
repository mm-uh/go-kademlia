package kademlia

import "sort"

type kademliaKBucket struct {
	start *linkedList
	last  *linkedList
	k     int
}

func (kB *kademliaKBucket) Update(c Contact) {
	nn := newLinkedListNode(c)

	// if the contact is already in the kBucket
	first := kB.start
	var prev *linkedList = nil
	for true {
		if first.value.GetNodeId() == c.GetNodeId() {
			kB.last.next = first
			kB.last = first
			if prev == nil {
				kB.start = first.next
			} else {
				prev.next = first.next
			}
			first.next = nil
			return
		}
		if first == kB.last {
			break
		}
		prev = first
		first = first.next
	}

	// if the kBucket is not full
	if kB.start.len() < kB.k {
		kB.last.next = nn
		kB.last = nn
		return
	}

	//if the kBucket is full
	head := kB.start
	if head.value.Ping() {
		kB.start = head.next
		head.next = nil
		kB.last.next = head
		kB.last = head
		return
	}

	kB.start = head.next
	kB.last.next = nn
	kB.last = nn
	return
}

func (kB *kademliaKBucket) GetClosestNodes(k int, nodeId uint64) []Contact {
	unorderedScl := sortableContactListFromLinkedList(kB.start, nodeId)
	sort.Sort(unorderedScl)
	scl := (*unorderedScl)[:k]
	contacts := make([]Contact, 0)
	for _, cd := range scl {
		contacts = append(contacts, cd.c)
	}
	return contacts
}

type distanceToContact struct {
	distance uint64
	c        Contact
}

func sortableContactListFromLinkedList(start *linkedList, nodeId uint64) *sortableContactList {
	scl := new(sortableContactList)
	for true {
		if start == nil {
			break
		}
		scl.append(&distanceToContact{
			c:        start.value,
			distance: XOR(nodeId, start.value.GetNodeId()),
		})
		start = start.next
	}
	return scl
}

type sortableContactList []*distanceToContact

func (scl *sortableContactList) append(dc *distanceToContact) {
	tscl := append(*scl, dc)
	scl = &tscl
}

func (scl *sortableContactList) Len() int {
	return len(*scl)
}

func (scl *sortableContactList) Swap(i, j int) {
	(*scl)[i], (*scl)[j] = (*scl)[j], (*scl)[i]
}

func (scl *sortableContactList) Less(i, j int) bool {
	return (*scl)[i].distance < (*scl)[j].distance
}

func newLinkedListNode(c Contact) *linkedList {
	return &linkedList{
		next:  nil,
		value: c,
	}
}

type linkedList struct {
	next  *linkedList
	value Contact
}

func (n *linkedList) len() int {
	if n.next == nil {
		return 1
	}
	return n.next.len() + 1
}
