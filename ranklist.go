package ranklist

import (
	"math/rand"
	"sync"
)

const (
	// 跳表的最大层数，设置为18层
	// Maximum number of levels in the skip list, set to 18
	MAXLEVEL = 18

	// 用于随机层级生成的概率值，设置为0.25
	// Probability used for random level generation, set to 0.25
	PROBABILITY = 0.25
)

// Ordered 接口定义了可用作键或值的类型约束
// ~符号表示包含所有以这些基本类型为底层类型的用户定义类型
// Ordered interface defines type constraints for keys and values
// The ~ symbol includes all user-defined types with these base types
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

// ZeroValue 返回指定类型的零值
// 这在需要返回默认值的场景中很有用
// Zero returns the zero value for the specified type
// This is useful in scenarios where a default value needs to be returned
func ZeroValue[K Ordered]() K {
	var zero K
	return zero
}

// Node 定义跳表节点的结构
// 包含键值对、前向指针数组、跨度数组和节点层级信息
// Node defines the structure of a skip list node
// Contains key-value pair, forward pointers array, span array, and node level
type Node[K Ordered, V Ordered] struct {
	// 节点的键
	// Key of the node
	key K

	// 节点的值
	// Value of the node
	value V

	// 每一层对应的前向指针数组
	// Array of forward pointers for each level
	forward []*Node[K, V]

	// 每一层对应的跨度数组，记录到下一个节点的距离
	// Array of spans for each level, recording distance to next node
	span []int

	// 当前节点的层级
	// Current level of the node
	level int
}

// RankList 定义跳表的核心结构
// 提供线程安全的节点管理，支持插入、删除、查找、排名等功能
// RankList defines the core structure of the skip list
// Provides thread-safe node management and supports insertion, deletion, retrieval, and ranking functionalities
type RankList[K Ordered, V Ordered] struct {
	sync.RWMutex

	// 跳表的头节点
	// Header node of the skip list
	header *Node[K, V]

	// 用于快速查找的键值对字典
	// Dictionary for fast key-value lookup
	dict map[K]*Node[K, V]

	// 当前跳表的最大层级
	// Current maximum level of the skip list
	level int

	// 跳表中的节点总数
	// Total number of nodes in the skip list
	length int
}

// NewNode 创建一个新的跳表节点
// 初始化节点的键、值、前向指针数组和跨度数组
// NewNode creates a new skip list node
// Initializes node's key, value, forward pointer array and span array
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
// 初始化头节点和键值对字典
// New creates a new skip list
// Initializes header node and key-value dictionary
func New[K Ordered, V Ordered]() *RankList[K, V] {
	return &RankList[K, V]{
		header: NewNode[K, V](ZeroValue[K](), ZeroValue[V](), MAXLEVEL),
		dict:   make(map[K]*Node[K, V]),
		level:  1,
	}
}

// randomLevel 随机生成节点的层级
// 使用概率PROBABILITY来决定是否增加层级，最高不超过MAXLEVEL
// randomLevel generates a random level for a new node
// Uses PROBABILITY to decide level increment, not exceeding MAXLEVEL
func randomLevel() int {
	level := 1
	for rand.Float64() < PROBABILITY && level < MAXLEVEL {
		level++
	}
	return level
}

// Set 向跳表中插入或更新节点
// 如果键已存在，则先删除旧节点再插入新节点
// Set inserts or updates a node in the skip list
// If the key exists, removes the old node before inserting the new one
func (sl *RankList[K, V]) Set(key K, value V) {
	sl.Lock()
	defer sl.Unlock()

	// 如果节点已存在，先删除旧节点
	// If node exists, remove old node first
	if node, exists := sl.dict[key]; exists {
		sl.del(node.key)
	}

	// 用于记录每层的前驱节点
	// Records predecessor nodes at each level
	prev := make([]*Node[K, V], MAXLEVEL)

	// 用于记录每层的排名值
	// Records rank values at each level
	rank := make([]int, MAXLEVEL)

	curr := sl.header

	// 生成新节点的随机层级
	// Generate random level for new node
	level := randomLevel()
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			prev[i] = sl.header
		}
		sl.level = level
	}

	// 查找插入位置并更新排名信息
	// Find insertion position and update rank information
	sum := 0
	for i := sl.level - 1; i >= 0; i-- {
		for curr.forward[i] != nil {

			if curr.forward[i].value > value ||
				(curr.forward[i].value == value && curr.forward[i].key > key) {
				break
			}
			sum += curr.forward[i].span[i]
			curr = curr.forward[i]
		}

		rank[i] = sum
		prev[i] = curr
	}

	// 创建并插入新节点
	// Create and insert new node
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

	// 更新高于新节点的层级的跨度
	// Update spans for levels above new node
	for i := level; i < sl.level; i++ {
		if prev[i].forward[i] != nil {
			prev[i].forward[i].span[i]++
		}
	}
	sl.length++
}

// Del 从跳表中删除指定键的节点
// 返回key是否存在
// Del removes a node with the specified key from the skip list
// Returns whether the deletion was successful
func (sl *RankList[K, V]) Del(key K) bool {
	sl.Lock()
	defer sl.Unlock()
	return sl.del(key)
}

// del performs the actual deletion of a node from the skip list.
// It searches for the node and updates the forward pointers and spans accordingly.
func (sl *RankList[K, V]) del(key K) bool {
	node, exists := sl.dict[key]
	if !exists {
		return false
	}

	// 记录每层的前驱节点
	// Record predecessor nodes at each level
	prev := make([]*Node[K, V], MAXLEVEL)
	curr := sl.header

	// 查找要删除的节点
	// Find the node to be deleted
	for i := sl.level - 1; i >= 0; i-- {
		for curr.forward[i] != nil &&
			(curr.forward[i].value < node.value ||
				(curr.forward[i].value == node.value && curr.forward[i].key < key)) {
			curr = curr.forward[i]
		}
		prev[i] = curr
	}

	// 更新前向指针和跨度
	// Update forward pointers and spans
	for i := 0; i < sl.level; i++ {
		if prev[i].forward[i] == node {
			prev[i].forward[i] = node.forward[i]
			if prev[i].forward[i] != nil {
				prev[i].forward[i].span[i] += node.span[i] - 1
			}
		}
	}

	// 更新跳表的最大层级
	// Update maximum level of skip list
	for sl.level > 1 && sl.header.forward[sl.level-1] == nil {
		sl.level--
	}

	delete(sl.dict, key)
	sl.length--
	return true
}

// Get 根据键获取节点的值
// 返回值和是否存在的标志
// Get retrieves the value associated with the key
// Returns the value and whether it exists
func (sl *RankList[K, V]) Get(key K) (V, bool) {
	sl.RLock()
	defer sl.RUnlock()

	if node, exists := sl.dict[key]; exists {
		return node.value, true
	}
	return ZeroValue[V](), false
}

// Rank 获取节点的排名
// 返回排名和是否存在的标志
// Rank gets the rank of a node
// Returns the rank and whether the node exists
func (sl *RankList[K, V]) Rank(key K) (int, bool) {
	sl.RLock()
	defer sl.RUnlock()

	node, exists := sl.dict[key]
	if !exists {
		return 0, false
	}

	// 计算节点的排名
	// Calculate node's rank
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

// Print 打印跳表的结构
// 用于调试和可视化跳表
// Print prints the structure of the skip list
// Used for debugging and visualization
// func (sl *RankList[K, V]) Print() {
// 	fmt.Printf("SkipList Level: %d, Length: %d\n", sl.level, sl.length)
// 	for i := sl.level - 1; i >= 0; i-- {
// 		current := sl.header
// 		fmt.Printf("L%d -> ", i+1)
// 		for current != nil {
// 			if current != sl.header {
// 				fmt.Printf("[%v:%v:%v] -> ", current.key, current.value, current.span[i])
// 			}
// 			current = current.forward[i]
// 		}
// 		fmt.Println("NIL")
// 	}
// 	fmt.Println("===================================")
// }
