package main

import (
	"github.com/faiface/pixel"
	"main/event"
)

type Location struct {
	x        int32
	y        int32
	offset_x float64
	offset_y float64
	offset   []pixel.Vec       //偏移
	canMove  map[string]*MoveS //可以移动的点
}

type UserStatus int

type ViewHero struct {
	*pixel.Sprite
	status int
}

type HeroBuffManger struct {
	buffProfit map[string]interface{} //记录buff 触发的属性

	//todo 每个人的buff 会订阅 频道 的 消息

	eventContext map[string]interface{}

	*BuffContext //英雄Buff

}

type Benefit struct {
	buffBenefit  map[string]int
	skillBenefit map[string]int
}

type Hero struct {
	Hp             int      //人物血量
	Mp             int      //人物MP
	attack         int      //人物攻击
	MovePoint      int      //人物移动点
	location       Location //人物当前坐标
	ViewHero                //
	HeroBuffManger          //Buff 管理器
	HeroSkillManger

	breakaction chan bool

	curRoom *Room
	skills  map[int]BaicSkill
	event.Topic
}

func (hero *Hero) AddHp(val int, event int) {
	hero.Hp += val
	//发送 加攻击力 事件

	//room 事件里丢事件消息
	switch val >= 0 {

	case true:
		//加血 可能来自队友的技能 buff 的技能 或者其它

	case false:
		//扣血可能来自buff 敌方的攻击 或者其它

	}

}

func (hero *Hero) AddAttack(val int, event int) {
	hero.attack += val
	//发送 加攻击力 事件

	//room 事件里丢事件消息
	switch val >= 0 {

	case true:
		//加血 可能来自队友的技能 buff 的技能 或者其它

	case false:
		//扣血可能来自buff 敌方的攻击 或者其它

	}

}

func (h *Hero) MoveToNextGrid(step float64) {
	if len(h.location.offset) < 1 {
		return
	}
	switch CaseDirection(h.location.offset[0]) {
	case UP:
		h.location.offset_y += step
		if h.location.offset_y >= 1 {
			h.location.y += 1
			h.location.offset_y = 0
			h.location.offset = h.location.offset[1:]
		}
	case DOWN:
		h.location.offset_y -= step
		if h.location.offset_y <= -1 {
			h.location.y -= 1
			h.location.offset_y = 0
			h.location.offset = h.location.offset[1:]
		}
	case RIGHT:
		h.location.offset_x += step
		if h.location.offset_x >= 1 {
			h.location.x += 1
			h.location.offset_x = 0
			h.location.offset = h.location.offset[1:]
		}
	case LEFT:
		h.location.offset_x -= step
		if h.location.offset_x <= -1 {
			h.location.x -= 1
			h.location.offset_x = 0
			h.location.offset = h.location.offset[1:]
		}
	}
}
