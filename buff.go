package main

import "fmt"

type BasicBuff struct {
	ID            int //BufferID
	OwnerID       int //拥有者ID
	FromSkillID   int //来自技能ID
	BuffCondition func(...interface{}) bool

	TypeId    int
	Exerciser int //释放者
	Owner     int // 拥有者

}

type BuffAttr map[int]interface{}

type condition []int

type BuffContext struct {
	Handlers  []func(*BuffContext, *Hero, map[BUffMsg]interface{})
	index     int8
	BuffEvent map[int]map[string]interface{} //记录每个回合buff 的一些事件
}

func (c *BuffContext) TriggerBuff(hero Hero, buffMsg map[BUffMsg]interface{}) *Hero {
	c.Handlers[0](c, &hero, buffMsg)
	return &hero
}

func (c *BuffContext) Next(hero *Hero, buffMsg map[BUffMsg]interface{}) {
	c.index++
	for c.index < int8(len(c.Handlers)) {
		c.Handlers[c.index](c, hero, buffMsg)
		c.index++
	}
}

func init() {
	c := &BuffContext{}
	c.Handlers = make([]func(*BuffContext, *Hero, map[BUffMsg]interface{}), 0)
	// 注册中间件
	c.Handlers = append(c.Handlers, AddHpBuff)
	c.Handlers = append(c.Handlers, AddMpBuff)
	c.Handlers = append(c.Handlers, AddMovePointBuff)
	msg := make(map[BUffMsg]interface{})
	// 控制器函数
	c.TriggerBuff(Hero{}, msg)
	BuffManager = make(map[int]func(c *BuffContext, hero *Hero, buffMsg map[BUffMsg]interface{}))

	BuffManager[1] = AddHpBuff
	BuffManager[2] = AddMpBuff
	BuffManager[3] = AddMovePointBuff

}

var BuffManager map[int]func(c *BuffContext, hero *Hero, buffMsg map[BUffMsg]interface{})

type BUffMsg int

const (
	AddAttack BUffMsg = 1
	AddMove   BUffMsg = 2
)

func AddHpBuff(c *BuffContext, hero *Hero, buffMsg map[BUffMsg]interface{}) {
	fmt.Println("添加血量")

	//判断Buff 的触发实际
	buffMsg[AddAttack] = 1
	//判断 当前的回合 或者是 状态
	hero.Hp += 1
	c.Next(hero, buffMsg)
	fmt.Println("添加血量结束")
}

//如果之前添加过 属性 翻倍
func DoubleAdd(c *BuffContext, hero *Hero, buffMsg map[BUffMsg]interface{}) {
	fmt.Println("添加属性属性翻倍")
	val1, ok1 := buffMsg[AddAttack]
	val2, ok2 := buffMsg[AddMove]
	if ok1 {
		hero.AddAttack(val1.(int), 1)
	} else if ok2 {
		hero.AddHp(val2.(int), 2)
	}
	hero.Mp += 1
	c.Next(hero, buffMsg)
	fmt.Println("添加攻击力结束")
}

func AddMpBuff(c *BuffContext, hero *Hero, buffMsg map[BUffMsg]interface{}) {
	fmt.Println("添加攻击力")
	hero.Mp += 1
	c.Next(hero, buffMsg)
	fmt.Println("添加攻击力结束")
}

func AddMovePointBuff(c *BuffContext, hero *Hero, buffMsg map[BUffMsg]interface{}) {
	fmt.Println("添加移动力")
	hero.MovePoint += 1
	c.Next(hero, buffMsg)
	fmt.Println("添加移动力结束")
}
