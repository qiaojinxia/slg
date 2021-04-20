package main

import (
	"errors"
	"fmt"
	"github.com/faiface/pixel"
	"math"
	"sort"
)

type Node struct {
	Parent *Node
	Pos    []int32
	GDir   []int32
	G      float32
	H      float32
	F      float32
	B      int // 障碍Block
	T      int // 转弯turn
}

func (node *Node) String() string {
	return fmt.Sprintf("%#v", node)
}

// 优先障碍最少，其次路径最短，再次转弯最少
func (best *Node) IsBetterThan(node *Node) bool {
	if best.B < node.B {
		return true
	}
	if best.B > node.B {
		return false
	}
	if best.F < node.F {
		return true
	}
	if best.F > node.F {
		return false
	}
	if best.T < node.T {
		return true
	}
	if best.T > node.T {
		return false
	}
	return false
}

type NodeList []*Node

func (nl NodeList) Len() int      { return len(nl) }
func (nl NodeList) Swap(i, j int) { nl[i], nl[j] = nl[j], nl[i] }
func (nl NodeList) Less(i, j int) bool {
	best, node := nl[i], nl[j]
	return best.IsBetterThan(node)
}

func (nl NodeList) String() string {
	s := "*******NodeList********\n"
	for i, node := range nl {
		s = fmt.Sprintf("%s%d = %s\n", s, i, node.String())
	}
	return s
}

func (nl NodeList) NodeIdx(node *Node) int {
	for i, n := range nl {
		if n.Pos[0] == node.Pos[0] && n.Pos[1] == node.Pos[1] {
			return i
		}
	}
	return -1
}

type AStar struct {
	W        int32
	H        int32
	SPos     []int32
	EPos     []int32
	Open     NodeList // 使用list，起点，终点相同时每次的路径都一样，换成map时，每次的路径可能不一样。
	CloseMap map[int32]struct{}
}

const (
	NODE_OPEN    int32 = iota // 可通行
	NODE_BARRIER              // 障碍，无法到达目的地时，返回障碍最少的路径
	NODE_NO_PASS              // 不可通行，边界等不能通行
)

func (a *AStar) newNode(parent *Node, x, y int32, dir []int32) *Node {
	state := a.GetNodeState(x, y)
	if state != NODE_OPEN && state != NODE_BARRIER {
		a.addClose(&Node{Pos: []int32{x, y}})
		return nil
	}
	b, t := 0, 0
	if parent != nil {
		b = parent.B
		if dir[0] == parent.GDir[0] && dir[1] == parent.GDir[1] {
			t = parent.T
		} else {
			t = parent.T + 1
		}
	}
	if state == NODE_BARRIER {
		b++
	}
	pos := []int32{x, y}
	g := a.getG(pos, parent)
	h := a.getH(pos, a.EPos)
	return &Node{
		Parent: parent,
		Pos:    pos,
		G:      g,
		H:      h,
		F:      g + h,
		B:      b,
		T:      t,
		GDir:   dir,
	}
}

func (a *AStar) GetNodeState(y, x int32) int32 {
	if x < 0 || x >= a.W || y < 0 || y >= a.H {
		return NODE_NO_PASS
	}
	//fmt.Println("GetNodeState", x, y, GMap[x][y])
	return GMap[x][y]
}

// 计算g值；直走=1；斜走=1.4
func (a *AStar) getG(pos1 []int32, parent *Node) float32 {
	if parent == nil {
		return 0
	}
	pos2 := parent.Pos
	g := parent.G
	if pos1[0] == pos2[0] {
		return g + float32(math.Abs(float64(pos1[1]-pos2[1])))
	} else if pos1[1] == pos2[1] {
		return g + float32(math.Abs(float64(pos1[0]-pos2[0])))
	} else {
		return g + float32(math.Abs(float64(pos1[0]-pos2[0])*1.4))
	}
}

// 计算h值
func (a *AStar) getH(pos1, pos2 []int32) float32 {
	return float32(math.Abs(float64(pos1[0]-pos2[0])) + math.Abs(float64(pos1[1]-pos2[1])))
}

func (a *AStar) addOpen(node *Node) {
	//idx := a.getNodeIdxInOpen(node)
	idx := a.Open.NodeIdx(node)
	//fmt.Println(a.Open)
	//fmt.Println("addOpen:", idx, "node:", node, "patrnt ", node.Parent)
	if idx >= 0 {
		n := a.Open[idx]
		if node.IsBetterThan(n) {
			a.Open[idx] = node
		}
	} else {
		a.Open = append(a.Open, node)
	}
}

// x, y 对应的node在open中的idx
func (a *AStar) getNodeIdxInOpen(node *Node) int {
	x, y := node.Pos[0], node.Pos[1]
	l := len(a.Open)
	if l == 0 {
		return -1
	}
	idx := sort.Search(len(a.Open),
		func(i int) bool {
			pos := a.Open[i].Pos
			return pos[0] == x && pos[1] == y
		})
	// sort.Search 对已排序的list使用二分查找，
	// sort.Search 的坑，未找到返回的值等于list长度
	if idx == l {
		return -1
	}
	return idx
}

// 从open中出栈，即最优的格子
func (a *AStar) getMinNode() *Node {
	sort.Sort(a.Open)
	node := a.Open[0]
	a.Open = a.Open[1:]
	return node
}

// 添加到close中,
func (a *AStar) addClose(n *Node) {
	x, y := n.Pos[0], n.Pos[1]
	key := x*a.W*10 + y
	a.CloseMap[key] = struct{}{}
}

// 判断是否在close中
func (a *AStar) isInClose(x, y int32) (ok bool) {
	key := x*a.W*10 + y
	_, ok = a.CloseMap[key]
	return
}

// 拓展周边方向的node
func (a *AStar) extendNeighbours(c *Node) {
	for _, dir := range GDir {
		x := c.Pos[0] + dir[0]
		y := c.Pos[1] + dir[1]
		if a.isInClose(x, y) {
			continue
		}
		node := a.newNode(c, x, y, dir)
		if node == nil {
			continue
		}
		a.addOpen(node)
	}
}

func (a *AStar) isTarget(n *Node, ePos []int32) bool {
	if n == nil {
		return false
	}
	if n.Pos[0] == ePos[0] && n.Pos[1] == ePos[1] {
		return true
	}
	return false
}

// 从结束点回溯到开始点，开始点的parent == None
func (a *AStar) makePath(p *Node) [][]int32 {
	path := make([][]int32, 0)
	for p != nil {
		path = append([][]int32{p.Pos}, path[:]...)
		p = p.Parent
	}
	fmt.Println("********makePath:", path)
	return path
}

// 路径查找主函数
func (a *AStar) findPath(sPos, ePos []int32) (path [][]int32, block, turn int, err error) {
	state := a.GetNodeState(sPos[0], sPos[1])
	if state != NODE_OPEN && state != NODE_BARRIER {
		err = errors.New(fmt.Sprintf("spos state is %d", state))
		return
	}

	estate := a.GetNodeState(ePos[0], ePos[1])
	if estate != NODE_OPEN && estate != NODE_BARRIER {
		err = errors.New(fmt.Sprintf("ePos state is %d", estate))
		return
	}

	a.SPos, a.EPos = sPos, ePos
	// 构建开始节点
	s := a.newNode(nil, sPos[0], sPos[1], []int32{0, 0})
	a.addOpen(s)
	for {
		if len(a.Open) <= 0 {
			err = errors.New("not find Open list is nil")
			return
		}
		p := a.getMinNode()
		//fmt.Println(fmt.Sprintf("---min p = %#v", p))
		if a.isTarget(p, ePos) {
			path = a.makePath(p)
			block, turn = p.B, p.T
			return
		}
		a.extendNeighbours(p)
		a.addClose(p)
	}
	return
}

// 查询过的所有格子
func (a *AStar) getSearched() []int32 {
	l := make([]int32, 0)
	for _, i := range a.Open {
		if i != nil {
			l = append(l, i.Pos[0], i.Pos[1])
		}
	}
	for k, _ := range a.CloseMap {
		x := k / a.W * 10
		y := k % a.W * 10
		l = append(l, x, y)
	}
	return l
}

// 清空open和close
func (a *AStar) clean() {
	a.Open = make(NodeList, 0)
	a.CloseMap = make(map[int32]struct{})
}

func NewAstar() *AStar {
	a := &AStar{
		W:        int32(len(GMap)),
		H:        int32(len(GMap[0])),
		Open:     make(NodeList, 0),
		CloseMap: make(map[int32]struct{}),
	}
	fmt.Println("NewAstar", a.W, a.H)
	return a
}

// 终点不唯一时找到最优的路径
func FindPathIgnoreBlock(sPos []int32, ePos [][]int32) [][]int32 {
	a := NewAstar()
	if a == nil {
		return nil
	}
	var path [][]int32
	var block, turn int
	for _, e := range ePos {
		a.clean()
		p, b, t, err := a.findPath(sPos, e)
		fmt.Println("FindPathIgnoreBlock:",
			fmt.Sprint("sPos", sPos),
			fmt.Sprint("ePos", e),
			fmt.Sprint("path", p),
			fmt.Sprint("Block", b),
			fmt.Sprint("turn", t),
			fmt.Sprint("err", err))
		if err != nil {
			continue
		}
		if path == nil {
			path, block, turn = p, b, t
			continue
		}
		stepNum, newStep := len(path), len(p)
		fmt.Println("stepNum", stepNum, "newStep", newStep)
		if stepNum < newStep {
			continue
		}
		if stepNum > newStep {
			path, block, turn = p, b, t
			continue
		}
		if block < b {
			continue
		}
		if block > b {
			path, block, turn = p, b, t
			continue
		}

		if turn < t {
			continue
		}

		if turn > b {
			path, block, turn = p, b, t
			continue
		}
	}
	fmt.Println("********end, FindPathIgnoreBlock", path, block, turn)
	return path
}

var (
	GDir [][]int32 //方向

)

type Direction int

const (
	UP     Direction = 1
	DOWN   Direction = 2
	LEFT   Direction = 3
	RIGHT  Direction = 4
	Origin Direction = 5
)

func equaleSlice(vec pixel.Vec, arr2 []int32) bool {
	return int32(vec.X) == arr2[0] && int32(vec.Y) == arr2[1]
}

func PointToDirection(location [][]int32) (pixel.Vec, Direction) {
	l1 := location[0]
	l2 := location[1]
	x := l1[0] - l2[0]
	y := l1[1] - l2[1]
	pos := pixel.Vec{X: float64(x), Y: float64(y)}
	if equaleSlice(pos, GDir[0]) {
		return pos, LEFT
	} else if equaleSlice(pos, GDir[1]) {
		return pos, UP
	} else if equaleSlice(pos, GDir[2]) {
		return pos, DOWN
	} else if equaleSlice(pos, GDir[3]) {
		return pos, RIGHT
	} else if equaleSlice(pos, []int32{0, 0}) {
		return pos, Origin
	}
	panic("error!")
}

func CaseDirection(dircection pixel.Vec) Direction {
	if equaleSlice(dircection, GDir[0]) {
		return LEFT
	} else if equaleSlice(dircection, GDir[1]) {
		return UP
	} else if equaleSlice(dircection, GDir[2]) {
		return DOWN
	} else if equaleSlice(dircection, GDir[3]) {
		return RIGHT
	} else if equaleSlice(dircection, []int32{0, 0}) {
		return Origin
	}
	panic("error!")
}

func init() {
	GDir = [][]int32{{1, 0}, {0, 1}, {0, -1}, {-1, 0}}
	fmt.Println("----init:", len(GMap[0]), len(GMap), NODE_OPEN, NODE_BARRIER, NODE_NO_PASS)
}
