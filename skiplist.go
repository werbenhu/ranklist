package skiplist

import (
	"math/rand/v2"
)

const (
	MAXLEVEL    = 18   // 最大层级
	PROBABILITY = 0.25 // 每层节点生成的概率
)

// Ordered 定义可比较的类型集合
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

// Zero 返回类型 K 的默认零值
func Zero[K Ordered]() K {
	var zero K // 声明一个 K 类型的零值变量
	return zero
}

// Node Skiplist中的节点结构
type Node[K Ordered, V Ordered] struct {
	key     K             // 节点的键
	value   V             // 节点的分数
	forward []*Node[K, V] // 指向下一个节点的指针数组
	span    []int         // 存储每层跨越的节点数量
	level   int           // 节点的层级
}

// Skiplist Skiplist结构
type Skiplist[K Ordered, V Ordered] struct {
	header *Node[K, V]       // 跳表的头节点
	dict   map[K]*Node[K, V] // 用于存储键值对的字典，方便快速查找
	level  int               // 当前跳表的最大层级
	length int               // 跳表中节点的总数
}

// NewNode 创建一个新的节点
func NewNode[K Ordered, V Ordered](key K, value V, level int) *Node[K, V] {
	return &Node[K, V]{
		key:     key,
		value:   value,
		forward: make([]*Node[K, V], level), // 根据节点的层级分配指针数组
		span:    make([]int, level),         // 根据节点的层级分配跨度数组
		level:   level,
	}
}

// NewSkiplist 创建一个新的跳表
func NewSkiplist[K Ordered, V Ordered]() *Skiplist[K, V] {
	return &Skiplist[K, V]{
		header: NewNode[K, V](Zero[K](), Zero[V](), MAXLEVEL), // 创建一个头节点
		dict:   make(map[K]*Node[K, V]),                       // 初始化字典
		level:  1,                                             // 初始跳表层级为1
	}
}

// 随机生成一个层级
func randomLevel() int {
	level := 1
	for rand.Float64() < PROBABILITY && level < MAXLEVEL {
		level++ // 随机生成层级，直到概率不满足或达到最大层级
	}
	return level
}

// Insert 向跳表中插入一个新的节点
func (sl *Skiplist[K, V]) Insert(key K, value V) {
	// 如果节点已经存在，先删除旧节点
	if node, exists := sl.dict[key]; exists {
		sl.Delete(node.key)
	}

	update := make([]*Node[K, V], MAXLEVEL) // 用于存储每一层的前驱节点
	rankArray := make([]int, MAXLEVEL)      // 用于记录每一层的rank值
	current := sl.header                    // 从头节点开始遍历

	// 遍历每一层
	for i := sl.level - 1; i >= 0; i-- {
		// 更新rankArray，用于计算当前节点的排名
		if i == sl.level-1 {
			rankArray[i] = 0
		} else {
			rankArray[i] = rankArray[i+1]
		}

		// 寻找合适的位置
		for current.forward[i] != nil &&
			(current.forward[i].value < value || // 按照score升序排列
				(current.forward[i].value == value && current.forward[i].key < key)) { // 如果score相同，则按照key升序排列
			rankArray[i] += current.span[i]
			current = current.forward[i]
		}
		update[i] = current
	}

	// 随机生成节点的层级
	level := randomLevel()
	if level > sl.level {
		// 如果生成的层级大于当前跳表的层级，需要增加跳表的层级
		for i := sl.level; i < level; i++ {
			update[i] = sl.header
			sl.header.span[i] = sl.length
		}
		sl.level = level
	}

	// 创建新节点
	newNode := NewNode[K, V](key, value, level)
	sl.dict[key] = newNode

	// 插入节点并更新前驱节点的forward指针
	for i := 0; i < level; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode

		// 更新span数组
		newNode.span[i] = update[i].span[i] - (rankArray[0] - rankArray[i])
		update[i].span[i] = rankArray[0] - rankArray[i] + 1
	}

	// 更新层级大于新节点层级的前驱节点的span
	for i := level; i < sl.level; i++ {
		if update[i] != nil {
			update[i].span[i]++
		}
	}

	// 跳表节点数量增加
	sl.length++
}

// Delete 从跳表中删除指定的节点
func (sl *Skiplist[K, V]) Delete(key K) bool {
	node, exists := sl.dict[key]
	if !exists {
		return false // 如果节点不存在，返回false
	}

	update := make([]*Node[K, V], sl.level) // 用于存储每一层的前驱节点
	current := sl.header                    // 从头节点开始遍历

	// 遍历每一层，找到删除节点的位置
	for i := sl.level - 1; i >= 0; i-- {
		for current.forward[i] != nil &&
			(current.forward[i].value < node.value ||
				(current.forward[i].value == node.value && current.forward[i].key < key)) {
			current = current.forward[i]
		}
		update[i] = current
	}

	// 删除节点
	for i := 0; i < sl.level; i++ {
		if update[i].forward[i] != nil && update[i].forward[i].key == key {
			update[i].span[i] += update[i].forward[i].span[i] - 1
			update[i].forward[i] = update[i].forward[i].forward[i]
		}
	}

	// 如果最上层没有节点，降低跳表的层级
	for i := sl.level - 1; i >= 0; i-- {
		if sl.header.forward[i] == nil {
			sl.level--
		} else {
			break
		}
	}

	// 删除字典中的键
	delete(sl.dict, key)
	sl.length-- // 跳表节点数量减少
	return true
}

// GetScore 获取指定key的score值
func (sl *Skiplist[K, V]) GetScore(key K) V {
	if node, exists := sl.dict[key]; exists {
		return node.value // 返回节点的score
	}
	return Zero[V]() // 如果节点不存在，返回0
}

// GetRank 获取指定key的rank（排名）
func (sl *Skiplist[K, V]) GetRank(key K) int {
	node, exists := sl.dict[key]
	if !exists {
		return 0 // 如果节点不存在，返回0
	}

	rank := 0
	current := sl.header
	value := node.value

	// 遍历每一层，计算节点的排名
	for i := sl.level - 1; i >= 0; i-- {
		for current.forward[i] != nil &&
			(current.forward[i].value < value ||
				(current.forward[i].value == value && current.forward[i].key <= key)) {

			rank += current.span[i]
			current = current.forward[i]
		}

		// 如果当前节点的 forward[i] 即为目标节点，直接返回排名
		if current.forward[i] != nil && current.forward[i].key == key {
			rank += current.span[i] // 累加跨度
			return rank
		}
	}
	return rank
}

// Print 打印跳表的结构，展示各层的key和span
func (sl *Skiplist[K, V]) Print() {
	for i := sl.level - 1; i >= 0; i-- {
		current := sl.header.forward[i] // 从该层的第一个节点开始
		print("Level ", i, ": ")

		// 遍历当前层的所有节点
		for current != nil {
			print("[", current.key, ",", current.span[i], "] ") // 打印key和span
			current = current.forward[i]                        // 移动到当前层的下一个节点
		}
		print("\n")
	}
}
