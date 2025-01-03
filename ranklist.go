package ranklist

import (
	"math/rand"
	"sync"
)

const (
	// 跳表的最大层数。对于1000万数据量，根据公式 h = log₂(n)/2，
	// 大约需要 log₂(10⁷)/2 ≈ 12 层。请根据实际业务需求和数据量大小适当调整该值，以优化性能和空间利用。
	// Maximum number of levels in the skip list. For 10 million elements,
	// using formula h = log₂(n)/2, we need approximately log₂(10⁷)/2 ≈ 12 levels.
	// Please adjust this value based on your actual business needs and data size to optimize performance and space utilization.
	MaxLevel = 12

	// 用于随机层级生成的概率值，设置为0.25
	// Probability used for random level generation, set to 0.25
	Probability = 0.25
)

// Ordered 接口定义了可用作键或值的类型约束
// Ordered interface defines type constraints for keys and values
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

// ZeroValue 返回指定类型的零值
// Zero returns the zero value for the specified type
func ZeroValue[K Ordered]() K {
	var zero K
	return zero
}

// Entry  represents a key-value pair
type Entry[K Ordered, V Ordered] struct {
	Key   K
	Value V
}

// Node 定义跳表节点的结构
// Node defines the structure of a skip list node
type Node[K Ordered, V Ordered] struct {
	// 节点的键值对
	// Key-value pair of the node
	data Entry[K, V]

	// 每一层对应的前向指针数组
	// Array of forward pointers for each level
	forward [MaxLevel]*Node[K, V]

	// 每一层对应的跨度数组，记录到下一个节点的距离
	// Array of spans for each level, recording distance to next node
	span [MaxLevel]int

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
	dict map[K]V

	// 当前跳表的最大层级
	// Current maximum level of the skip list
	level int

	// 跳表中的节点总数
	// Total number of nodes in the skip list
	length int
}

// NewNode 创建一个新的跳表节点
// NewNode creates a new skip list node
func NewNode[K Ordered, V Ordered](key K, value V, level int) *Node[K, V] {
	return &Node[K, V]{
		data:  Entry[K, V]{Key: key, Value: value},
		level: level,
	}
}

// New 创建一个新的跳表
// New creates a new skip list
func New[K Ordered, V Ordered]() *RankList[K, V] {
	return &RankList[K, V]{
		header: NewNode[K, V](ZeroValue[K](), ZeroValue[V](), MaxLevel),
		dict:   make(map[K]V),
		level:  1,
	}
}

// randomLevel 随机生成节点的层级
// 使用概率Probability来决定是否增加层级，最高不超过MaxLevel
// randomLevel generates a random level for a new node
// Uses Probability to decide level increment, not exceeding MaxLevel
func randomLevel() int {
	level := 1
	for rand.Float64() < Probability && level < MaxLevel {
		level++
	}
	return level
}

// Set 向跳表中插入数据
// 如果键已存在，则先删除旧节点再插入新节点
// Set inserts or updates a key-value pair
// If the key exists, removes the old node before inserting the new one
func (sl *RankList[K, V]) Set(key K, value V) {
	sl.Lock()
	defer sl.Unlock()

	// 如果节点已存在，先删除旧节点
	// If node exists, remove old node first
	if _, exists := sl.dict[key]; exists {
		sl.del(key)
	}

	// 用于记录每层的前驱节点
	// Records predecessor nodes at each level
	var prev [MaxLevel]*Node[K, V]

	// 用于记录每层的排名值
	// Records rank values at each level
	var rank [MaxLevel]int

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

			if curr.forward[i].data.Value > value ||
				(curr.forward[i].data.Value == value && curr.forward[i].data.Key > key) {
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
	sl.dict[key] = value
	for i := 0; i < level; i++ {
		newNode.forward[i] = prev[i].forward[i]
		prev[i].forward[i] = newNode
		if i == 0 {
			newNode.span[i] = 1
		} else {
			newNode.span[i] = rank[0] - rank[i] + 1
			if newNode.forward[i] != nil {
				// 后面节点的span，被插入的节点切割了
				newNode.forward[i].span[i] = newNode.forward[i].span[i] - newNode.span[i] + 1
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

// Length 返回跳表中当前元素的数量。
// Length returns the current number of elements in the skip list.
func (sl *RankList[K, V]) Length() int {
	sl.RLock()
	defer sl.RUnlock()
	return sl.length
}

// Del 从跳表中删除指定键的节点。
// 如果键存在并且节点被删除，返回true；如果键不存在，返回false。
// Del removes the node with the specified key from the skip list.
// Returns true if the key exists and the node is deleted, false if the key does not exist.
func (sl *RankList[K, V]) Del(key K) bool {
	sl.Lock()
	defer sl.Unlock()
	return sl.del(key)
}

// 删除操作实际执行跳表节点的删除。
// 它搜索指定的节点，更新前向指针，并相应地调整跨度值。
// del performs the actual deletion of a node from the skip list.
// It searches for the node, updates the forward pointers, and adjusts the span values accordingly.
func (sl *RankList[K, V]) del(key K) bool {
	value, exists := sl.dict[key]
	if !exists {
		return false
	}

	// 记录每层的前驱节点
	// Record predecessor nodes at each level
	var prev [MaxLevel]*Node[K, V]
	curr := sl.header

	// 查找要删除的节点
	// Find the node to be deleted
	for i := sl.level - 1; i >= 0; i-- {
		for curr.forward[i] != nil &&
			(curr.forward[i].data.Value < value ||
				(curr.forward[i].data.Value == value && curr.forward[i].data.Key < key)) {
			curr = curr.forward[i]
		}
		prev[i] = curr
	}

	// 更新前向指针和跨度
	// Update forward pointers and spans
	for i := 0; i < sl.level; i++ {
		curr = prev[i].forward[i]

		if curr != nil && curr.data.Key == key && curr.data.Value == value {
			// 如果这一层找到了删除的节点，那么将删除节点清除，并将删除节点的 span 甩给后面的节点
			// If the node to be deleted is found at this level, remove the node and pass its span to the next node
			prev[i].forward[i] = curr.forward[i]
			if curr.forward[i] != nil {
				curr.forward[i].span[i] += (curr.span[i] - 1)
			}

		} else if curr != nil {
			// 如果没有找到节点，说明这些层级比删除的节点的层级高，这些比删除节点高的节点，span 要 -1
			// If the node is not found, it means these levels are higher than the level of the node to be deleted,
			// so the span of these higher-level nodes needs to be decremented by 1
			curr.span[i]--
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
// 如果键存在并且节点被删除，返回true；如果键不存在，返回false。
// Get retrieves the value associated with the key
// Returns true if the key exists and the node is deleted, false if the key does not exist.
func (sl *RankList[K, V]) Get(key K) (V, bool) {
	sl.RLock()
	defer sl.RUnlock()

	if value, exists := sl.dict[key]; exists {
		return value, true
	}
	return ZeroValue[V](), false
}

// Rank 获取节点的排名
// 如果键存在并且节点被删除，返回true；如果键不存在，返回false。
// Rank gets the rank of a node
// Returns true if the key exists and the node is deleted, false if the key does not exist.
func (sl *RankList[K, V]) Rank(key K) (int, bool) {
	sl.RLock()
	defer sl.RUnlock()

	value, exists := sl.dict[key]
	if !exists {
		return 0, false
	}

	// 计算节点的排名
	// Calculate node's rank
	rank := 0
	curr := sl.header

	for i := sl.level - 1; i >= 0; i-- {
		for curr.forward[i] != nil {

			if curr.forward[i].data.Value == value && curr.forward[i].data.Key == key {
				rank += curr.forward[i].span[i]
				return rank, true
			}

			if curr.forward[i].data.Value > value ||
				(curr.forward[i].data.Value == value && curr.forward[i].data.Key > key) {
				break
			}

			rank += curr.forward[i].span[i]
			curr = curr.forward[i]
		}
	}
	return 0, false
}

// Range 获取指定排名区间内的榜单项（不包含END）
// 返回指定范围内的条目列表。
// Range retrieves the entries within the specified rank range (excluding END)
// Returns a list of entries within the specified range.
func (sl *RankList[K, V]) Range(start int, end int) []Entry[K, V] {
	sl.RLock()
	defer sl.RUnlock()

	rank := 0
	curr := sl.header
	entries := make([]Entry[K, V], 0)

	for i := sl.level - 1; i >= 0; i-- {
		for curr.forward[i] != nil {
			forwardRank := rank + curr.forward[i].span[i]
			if forwardRank >= start {
				break
			}
			rank = forwardRank
			curr = curr.forward[i]
		}
	}

	total := 0
	for curr.forward[0] != nil && start+total < end {
		entries = append(entries, curr.forward[0].data)
		curr = curr.forward[0]
		total++
	}
	return entries
}

// Print for test
// func (sl *RankList[K, V]) Print() {
// 	fmt.Printf("SkipList Level: %d, Length: %d\n", sl.level, sl.length)
// 	for i := sl.level - 1; i >= 0; i-- {
// 		curr := sl.header
// 		fmt.Printf("L%d -> ", i+1)
// 		for curr != nil {
// 			if curr != sl.header {
// 				fmt.Printf("[%v:%v:%v] -> ", curr.data.Key, curr.data.Value, curr.span[i])
// 			}
// 			curr = curr.forward[i]
// 		}
// 		fmt.Println("NIL")
// 	}
// 	fmt.Println("===================================")
// }
