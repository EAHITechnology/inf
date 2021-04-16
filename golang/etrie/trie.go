package etrie

import (
	"fmt"
	"strings"
	"sync"
)

type trieNode struct {
	pattern  string
	mutex    sync.RWMutex
	children map[string]*trieNode
}

func NewTrie() *trieNode {
	return &trieNode{
		pattern:  "/",
		children: make(map[string]*trieNode),
	}
}

func (this *trieNode) checkRouting(routing string) bool {
	if !strings.HasPrefix(routing, "/") {
		return false
	}
	if strings.HasSuffix(routing, "/") {
		return false
	}
	return true
}

func (this *trieNode) insert(routing string, parts []string, idx int) {
	if _, ok := this.children["/"+parts[idx]]; !ok {
		this.children["/"+parts[idx]] = NewTrie()
	}
	prouting := this.children["/"+parts[idx]]

	if len(parts)-1 == idx {
		prouting.pattern = routing
		return
	}
	newIdx := idx + 1
	prouting.insert(routing, parts, newIdx)
}

func (this *trieNode) Insert(routing string) error {
	if !this.checkRouting(routing) {
		return fmt.Errorf("illegal routing")
	}
	this.mutex.Lock()
	defer this.mutex.Unlock()

	parts := strings.Split(routing, "/")
	parts = parts[1:]
	this.insert(routing, parts, 0)
	return nil
}

func (this *trieNode) delete(parts []string, idx int) (bool, error) {
	if len(parts) == idx {
		if this.pattern == "/" {
			return true, fmt.Errorf("routing not exit!")
		}
		this.pattern = "/"

		if len(this.children) != 0 {
			return false, nil
		}

		return true, nil
	}
	_, ok := this.children["/"+parts[idx]]
	if !ok {
		return false, fmt.Errorf(fmt.Sprintf("%s not exit!", parts[idx]))
	}
	prouting := this.children["/"+parts[idx]]
	newIdx := idx + 1
	Ifdel, err := prouting.delete(parts, newIdx)
	if err != nil {
		return Ifdel, err
	}
	if Ifdel {
		delete(this.children, "/"+parts[idx])
		if len(this.children) == 0 && this.pattern == "/" {
			return true, nil
		}
	}
	return false, nil
}

func (this *trieNode) Delete(routing string) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if !this.checkRouting(routing) {
		return fmt.Errorf("illegal routing")
	}

	parts := strings.Split(routing, "/")
	parts = parts[1:]
	_, err := this.delete(parts, 0)
	if err != nil {
		return err
	}
	return nil
}

func (this *trieNode) search(routing string, parts []string, idx int) (bool, error) {
	if len(parts) == idx {
		if this.pattern != routing {
			return false, nil
		}
		return true, nil
	}

	_, ok := this.children["/"+parts[idx]]
	if !ok {
		return false, nil
	}
	prouting := this.children["/"+parts[idx]]
	newIdx := idx + 1
	return prouting.search(routing, parts, newIdx)
}

func (this *trieNode) Search(routing string) (bool, error) {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	if !this.checkRouting(routing) {
		return false, fmt.Errorf("illegal routing")
	}
	parts := strings.Split(routing, "/")
	parts = parts[1:]

	return this.search(routing, parts, 0)
}
