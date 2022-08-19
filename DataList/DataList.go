package DataList

import (
	"container/list"
	"fmt"
	"github.com/y-omicron/util/VT100C"
	"strings"
	"sync"
)

// DataNode 单个数据，最后利用 Key 来寻找参数
type DataNode struct {
	Key   string
	Value string
}

func (dn DataNode) String() string {
	return dn.Value
}
func (dn DataNode) Len() int {
	return len(dn.String())
}
func (dn *DataNode) New(key, value string) {
	dn.Key = key
	dn.Value = value
}
func (dn *DataNode) Set(value string) {
	dn.Value = value
}
func (dn *DataNode) Status() bool {
	if strings.Contains(dn.Value, "Running") {
		return true
	}
	return false
}

var maxLine int

// DataList 真正用来存储数据的结构体，有序列表
type DataList struct {
	dataMap  map[string]*list.Element
	dataList *list.List
	rwMutex  *sync.RWMutex
}

func (dl *DataList) Lock() {
	dl.rwMutex.Lock()
}
func (dl *DataList) Unlock() {
	dl.rwMutex.Unlock()
}

// New 初始化一个 DataList
func New() *DataList {
	return &DataList{
		dataMap:  make(map[string]*list.Element),
		dataList: list.New(),
		rwMutex:  &sync.RWMutex{},
	}
}

// Exists 检查节点是否已经存在，指针
func (dl *DataList) Exists(node *DataNode) bool {
	for n := dl.dataList.Front(); n != nil; n = n.Next() {
		if n.Value.(*DataNode) == node {
			return true
		}
	}
	return false
}

// Add 存入一个节点
func (dl *DataList) Add(node *DataNode) bool {
	if dl.Exists(node) {
		return false
	}
	elem := dl.dataList.PushBack(node)
	dl.dataMap[node.Key] = elem
	return true
}

// Remove 删除一个节点，当连接中断后或者不需要的时候主动删除掉
func (dl *DataList) Remove(node *DataNode) {
	if !dl.Exists(node) {
		return
	}
	dl.dataList.Remove(dl.dataMap[node.Key])
	delete(dl.dataMap, node.Key)
}

// Walk List中的遍历，对每个节点执行cb方法
func (dl *DataList) Walk(cb func(node *DataNode)) {
	for node := dl.dataList.Front(); node != nil; node = node.Next() {
		cb(node.Value.(*DataNode))
	}
}

func (dl *DataList) Len() int {
	return dl.dataList.Len()
}

func (dl *DataList) Show() {
	dl.print()
}
func (dl *DataList) Update() {
	dl.Lock()
	VT100C.CleanLine(maxLine)
	for node := dl.dataList.Front(); node != nil && !node.Value.(*DataNode).Status(); node = node.Next() {
		fmt.Printf("%s\n", node.Value.(*DataNode).String())
		dl.Remove(node.Value.(*DataNode))
	}
	dl.print()
	dl.Unlock()
	return
}
func (dl *DataList) print() {
	maxLine = 0
	dl.Walk(func(node *DataNode) {
		fmt.Printf("%s\n", node.String())
		maxLine++
	})
}
