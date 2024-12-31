package ranklist

import (
	"fmt"
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

// Node RankList中的节点结构
type Node[K Ordered, V Ordered] struct {
	key     K             // 节点的键
	value   V             // 节点的分数
	forward []*Node[K, V] // 指向下一个节点的指针数组
	span    []int         // 存储每层跨越的节点数量，是前一个节点到当前节点的中间的节点数量
	level   int           // 节点的层级
}

// RankList RankList结构
type RankList[K Ordered, V Ordered] struct {
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

// New 创建一个新的跳表
func New[K Ordered, V Ordered]() *RankList[K, V] {
	return &RankList[K, V]{
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
func (sl *RankList[K, V]) Set(key K, value V) {
	// 如果节点已经存在，先删除旧节点
	if node, exists := sl.dict[key]; exists {
		sl.Del(node.key)
	}

	prev := make([]*Node[K, V], MAXLEVEL) // 用于存储每一层的前驱节点
	rank := make([]int, MAXLEVEL)         // 用于记录每一层的rank值
	curr := sl.header                     // 从头节点开始遍历

	// 随机生成节点的层级
	level := randomLevel()
	if level > sl.level {
		// 如果生成的层级大于当前跳表的层级，需要增加跳表的层级
		for i := sl.level; i < level; i++ {
			prev[i] = sl.header
		}
		sl.level = level
	}

	sum := 0 // 累计的span数量
	// 遍历每一层
	for i := sl.level - 1; i >= 0; i-- {

		// 寻找合适的位置，如果value相等，后来的在后面
		for curr.forward[i] != nil {
			if curr.forward[i].value > value {
				break
			}
			if curr.forward[i].value == value && curr.forward[i].key > key {
				break
			}

			// 如果score相同，则按照key升序排列
			sum += curr.forward[i].span[i] // rank是前驱节点的排名值
			curr = curr.forward[i]
		}
		rank[i] = sum
		prev[i] = curr
	}

	// 创建新节点
	newNode := NewNode[K, V](key, value, level)
	sl.dict[key] = newNode

	// 插入节点并更新前驱节点的forward指针
	for i := 0; i < level; i++ {
		newNode.forward[i] = prev[i].forward[i]
		prev[i].forward[i] = newNode

		if i == 0 {
			// 最底下一层的span一定是1
			newNode.span[i] = 1
		} else {
			// 当前节点的span是第0层前驱节点的排名值，减去当前层前驱节点的排名 + 1
			newNode.span[i] = rank[0] - rank[i] + 1

			// 下一个节点的被新节点切割了，要重现计算下一个节点的span
			if newNode.forward[i] != nil {
				newNode.forward[i].span[i] = newNode.span[i] + 1
			}
		}
	}

	// 插入节点上面还有层级，上面层级当前节点后面的节点的span都要+1
	for i := level; i < sl.level; i++ {
		if prev[i].forward[i] != nil {
			prev[i].forward[i].span[i]++
		}
	}

	// 跳表节点数量增加
	sl.length++
}

// Del 从跳表中删除指定的节点
func (sl *RankList[K, V]) Del(key K) bool {
	node, exists := sl.dict[key]
	if !exists {
		return false // 如果节点不存在，返回false
	}

	prev := make([]*Node[K, V], MAXLEVEL) // 用于存储每一层的前驱节点
	curr := sl.header                     // 从头节点开始遍历

	// 遍历每一层，找到删除节点的位置
	for i := sl.level - 1; i >= 0; i-- {
		for curr.forward[i] != nil {
			if curr.forward[i].value == node.value && curr.forward[i].key >= key {
				break
			}
			if curr.forward[i].value > node.value {
				break
			}
			curr = curr.forward[i]
		}
		prev[i] = curr
	}

	// 删除节点
	for i := 0; i < sl.level; i++ {
		if prev[i].forward[i] != nil && prev[i].forward[i].key == key {
			prev[i].forward[i] = prev[i].forward[i].forward[i]
		}
		prev[i].forward[i].span[i]--
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
func (sl *RankList[K, V]) Get(key K) V {
	if node, exists := sl.dict[key]; exists {
		return node.value // 返回节点的score
	}
	return Zero[V]() // 如果节点不存在，返回0
}

// GetRank 获取指定key的rank（排名）
func (sl *RankList[K, V]) Rank(key K) int {
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

// PrintSkipList 打印整个跳表的结构
func (sl *RankList[K, V]) Print() {
	fmt.Printf("SkipList Level: %d, Length: %d\n", sl.level, sl.length)
	for i := sl.level - 1; i >= 0; i-- {
		current := sl.header
		fmt.Printf("L%d:", i+1)
		for current != nil {
			if current == sl.header {
				fmt.Printf(" ")
			} else {
				fmt.Printf("[%v:%v:%v] -> ", current.key, current.value, current.span[i])
			}
			current = current.forward[i]
		}
		fmt.Println("NIL")
	}
	fmt.Println("===================================")
}
