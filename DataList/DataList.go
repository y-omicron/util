package DataList

import (
	"ConUi/VT100C"
	"container/list"
	"fmt"
)

// DataNode 单个数据，最后利用 Key 来寻找参数
type DataNode struct {
	Key   string
	Value string
}

var maxLine int
var tmpLine int

func (dn DataNode) String() string {
	return fmt.Sprintf("%s: %s", dn.Key, dn.Value)
}
func (dn DataNode) Len() int {
	return len(dn.String())
}
func (dn *DataNode) New(key, value string) {
	dn.Set(key, value)
}
func (dn *DataNode) Set(key, value string) {
	dn.Key = key
	dn.Value = value
}

// DataList 真正用来存储数据的结构体，有序列表
type DataList struct {
	dataMap  map[string]*list.Element
	dataList *list.List
}

// New 初始化一个 DataList
func New() *DataList {
	return &DataList{
		dataMap:  make(map[string]*list.Element),
		dataList: list.New(),
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

func (dl DataList) Show() {
	maxLine = 0
	dl.Walk(func(node *DataNode) {
		fmt.Printf("%s\n", node.String())
		maxLine++
	})
}
func (dl DataList) Update() {
	VT100C.Move(VT100C.Up, maxLine)
	dl.print()
	if maxLine > tmpLine {
		var xLine = tmpLine - maxLine
		for i := 0; i < xLine; i++ {
			VT100C.CleanLine()
			fmt.Printf("\n")
		}
		VT100C.Move(VT100C.Up, xLine)
	}
	maxLine = tmpLine
	return
}
func (dl DataList) print() {
	tmpLine = 0
	dl.Walk(func(node *DataNode) {
		VT100C.CleanLine()
		fmt.Printf("%s\n", node.String())
		tmpLine++
	})

}
