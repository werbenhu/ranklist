package skiplist

import (
	"math/rand/v2"
)

const (
	MAXLEVEL    = 18   // 最大层级
	PROBABILITY = 0.25 // 每层节点生成的概率
)

// Skiplist中的节点结构
type Node struct {
	key     int     // 节点的键
	score   float64 // 节点的分数
	forward []*Node // 指向下一个节点的指针数组
	span    []int   // 存储每层跨越的节点数量
	level   int     // 节点的层级
}

// Skiplist结构
type Skiplist struct {
	header *Node         // 跳表的头节点
	dict   map[int]*Node // 用于存储键值对的字典，方便快速查找
	level  int           // 当前跳表的最大层级
	length int           // 跳表中节点的总数
}

// 创建一个新的节点
func NewNode(key int, score float64, level int) *Node {
	return &Node{
		key:     key,
		score:   score,
		forward: make([]*Node, level), // 根据节点的层级分配指针数组
		span:    make([]int, level),   // 根据节点的层级分配跨度数组
		level:   level,
	}
}

// 创建一个新的跳表
func NewSkiplist() *Skiplist {
	return &Skiplist{
		header: NewNode(0, 0, MAXLEVEL), // 创建一个头节点
		dict:   make(map[int]*Node),     // 初始化字典
		level:  1,                       // 初始跳表层级为1
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

// 向跳表中插入一个新的节点
func (sl *Skiplist) Insert(key int, score float64) {
	// 如果节点已经存在，先删除旧节点
	if node, exists := sl.dict[key]; exists {
		sl.Delete(node.key)
	}

	update := make([]*Node, MAXLEVEL)  // 用于存储每一层的前驱节点
	rankArray := make([]int, MAXLEVEL) // 用于记录每一层的rank值
	current := sl.header               // 从头节点开始遍历

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
			(current.forward[i].score < score || // 按照score升序排列
				(current.forward[i].score == score && current.forward[i].key < key)) { // 如果score相同，则按照key升序排列
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
	newNode := NewNode(key, score, level)
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

// 从跳表中删除指定的节点
func (sl *Skiplist) Delete(key int) bool {
	node, exists := sl.dict[key]
	if !exists {
		return false // 如果节点不存在，返回false
	}

	update := make([]*Node, sl.level) // 用于存储每一层的前驱节点
	current := sl.header              // 从头节点开始遍历

	// 遍历每一层，找到删除节点的位置
	for i := sl.level - 1; i >= 0; i-- {
		for current.forward[i] != nil &&
			(current.forward[i].score < node.score ||
				(current.forward[i].score == node.score && current.forward[i].key < key)) {
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

// 获取指定key的score值
func (sl *Skiplist) GetScore(key int) float64 {
	if node, exists := sl.dict[key]; exists {
		return node.score // 返回节点的score
	}
	return 0 // 如果节点不存在，返回0
}

// 获取指定key的rank（排名）
func (sl *Skiplist) GetRank(key int) int {
	node, exists := sl.dict[key]
	if !exists {
		return 0 // 如果节点不存在，返回0
	}

	rank := 0
	current := sl.header
	score := node.score

	// 遍历每一层，计算节点的排名
	for i := sl.level - 1; i >= 0; i-- {
		for current.forward[i] != nil &&
			(current.forward[i].score < score ||
				(current.forward[i].score == score && current.forward[i].key <= key)) {
			rank += current.span[i]
			current = current.forward[i]
		}
		if current.key == key {
			return rank // 找到节点并返回排名
		}
	}
	return rank
}
