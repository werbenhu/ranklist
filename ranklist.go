package ranklist

import (
	"fmt"
	"math/rand"
)

const (
	MAXLEVEL    = 18   // 最大层级
	PROBABILITY = 0.25 // 每层节点生成的概率
)

// Ordered 定义可比较的类型集合
// 支持所有基本数字类型和字符串类型
// 泛型约束定义为支持 ~ 运算符的类型
// 即底层类型是指定类型的用户自定义类型也可以
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

// Zero 返回类型 K 的默认零值
// 利用泛型的特性，返回一个零值
func Zero[K Ordered]() K {
	var zero K // 声明一个零值变量
	return zero
}

// Node 定义跳表中的节点结构
// 包含键、值、前向指针数组、跨距数组和节点层级
type Node[K Ordered, V Ordered] struct {
	key     K             // 节点的键
	value   V             // 节点的值
	forward []*Node[K, V] // 每层的前向指针数组
	span    []int         // 每层跨越节点的数量
	level   int           // 节点的层级
}

// RankList 定义跳表结构
// 包含头节点、键值对字典、当前最大层级和节点总数
type RankList[K Ordered, V Ordered] struct {
	header *Node[K, V]       // 跳表的头节点
	dict   map[K]*Node[K, V] // 用于快速查找键值对的字典
	level  int               // 当前跳表的最大层级
	length int               // 跳表中节点的总数
}

// NewNode 创建一个新的节点
// 初始化节点的键、值和层级
func NewNode[K Ordered, V Ordered](key K, value V, level int) *Node[K, V] {
	return &Node[K, V]{
		key:     key,
		value:   value,
		forward: make([]*Node[K, V], level),
		span:    make([]int, level),
		level:   level,
	}
}

// New 创建一个新的跳表
// 初始化头节点和字典
func New[K Ordered, V Ordered]() *RankList[K, V] {
	return &RankList[K, V]{
		header: NewNode[K, V](Zero[K](), Zero[V](), MAXLEVEL),
		dict:   make(map[K]*Node[K, V]),
		level:  1, // 初始层级为 1
	}
}

// randomLevel 随机生成节点的层级
// 根据 PROBABILITY 概率生成层级，最大不超过 MAXLEVEL
func randomLevel() int {
	level := 1
	for rand.Float64() < PROBABILITY && level < MAXLEVEL {
		level++
	}
	return level
}

// Set 向跳表中插入或更新节点
// 如果键已存在，则先删除旧节点再插入新节点
func (sl *RankList[K, V]) Set(key K, value V) {
	// 如果节点已存在，先删除旧节点
	if node, exists := sl.dict[key]; exists {
		sl.Del(node.key)
	}

	prev := make([]*Node[K, V], MAXLEVEL) // 每层的前驱节点
	rank := make([]int, MAXLEVEL)         // 每层的 rank 值
	curr := sl.header                     // 从头节点开始

	// 随机生成新节点的层级
	level := randomLevel()
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			prev[i] = sl.header
		}
		sl.level = level
	}

	sum := 0
	for i := sl.level - 1; i >= 0; i-- {
		for curr.forward[i] != nil {
			if curr.forward[i].value > value ||
				(curr.forward[i].value == value && curr.forward[i].key > key) {
				break
			}
			sum += curr.span[i]
			curr = curr.forward[i]
		}
		rank[i] = sum
		prev[i] = curr
	}

	// 创建并插入新节点
	newNode := NewNode(key, value, level)
	sl.dict[key] = newNode
	for i := 0; i < level; i++ {
		newNode.forward[i] = prev[i].forward[i]
		prev[i].forward[i] = newNode
		if i == 0 {
			newNode.span[i] = 1
		} else {
			newNode.span[i] = rank[0] - rank[i] + 1
			if newNode.forward[i] != nil {
				newNode.forward[i].span[i] = newNode.span[i] + 1
			}
		}
	}

	for i := level; i < sl.level; i++ {
		if prev[i].forward[i] != nil {
			prev[i].forward[i].span[i]++
		}
	}
	sl.length++
}

// Del 从跳表中删除指定键的节点
func (sl *RankList[K, V]) Del(key K) bool {
	node, exists := sl.dict[key]
	if !exists {
		return false
	}

	prev := make([]*Node[K, V], MAXLEVEL)
	curr := sl.header

	for i := sl.level - 1; i >= 0; i-- {
		for curr.forward[i] != nil &&
			(curr.forward[i].value < node.value ||
				(curr.forward[i].value == node.value && curr.forward[i].key < key)) {
			curr = curr.forward[i]
		}
		prev[i] = curr
	}

	for i := 0; i < sl.level; i++ {
		if prev[i].forward[i] == node {
			prev[i].forward[i] = node.forward[i]
			if prev[i].forward[i] != nil {
				prev[i].forward[i].span[i] += node.span[i] - 1
			}
		}
	}

	for sl.level > 1 && sl.header.forward[sl.level-1] == nil {
		sl.level--
	}

	delete(sl.dict, key)
	sl.length--
	return true
}

// Get 根据键获取节点的值
func (sl *RankList[K, V]) Get(key K) (V, bool) {
	if node, exists := sl.dict[key]; exists {
		return node.value, true
	}
	return Zero[V](), false
}

// Rank 获取节点的排名
func (sl *RankList[K, V]) Rank(key K) (int, bool) {
	node, exists := sl.dict[key]
	if !exists {
		return 0, false
	}

	rank := 0
	curr := sl.header

	for i := sl.level - 1; i >= 0; i-- {
		for curr.forward[i] != nil {

			if curr.forward[i].value == node.value && curr.forward[i].key == key {
				rank += curr.forward[i].span[i]
				return rank, true
			}

			if curr.forward[i].value > node.value ||
				(curr.forward[i].value == node.value && curr.forward[i].key > key) {
				break
			}

			rank += curr.forward[i].span[i]
			curr = curr.forward[i]
		}
	}
	return 0, false
}

// Print 打印跳表结构
func (sl *RankList[K, V]) Print() {
	fmt.Printf("SkipList Level: %d, Length: %d\n", sl.level, sl.length)
	for i := sl.level - 1; i >= 0; i-- {
		current := sl.header
		fmt.Printf("L%d -> ", i+1)
		for current != nil {
			if current != sl.header {
				fmt.Printf("[%v:%v:%v] -> ", current.key, current.value, current.span[i])
			}
			current = current.forward[i]
		}
		fmt.Println("NIL")
	}
	fmt.Println("===================================")
}
