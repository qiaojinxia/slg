package main

import "fmt"

//所有的操作皆为技能
//每触发 一个技能投递到任务 hero 的事件消息中
//用于 buff 或者技能消费处理 事件

type BaicSkill struct {
	name        string //技能名称
	skillID     int32  //技能ID
	attackRange []*Vec //攻击范围
	maxLayers   int    //最大叠加层数
}

type HeroSkillManger struct {
}
type ISkill interface {
	GetSkillId() int32
	GetSkillRange() []Vec
	GetSkillTargetID() int32
}

type YINBOJIAN struct {
	BaicSkill
}

func (y *YINBOJIAN) GetSkillId() int32 {
	return y.skillID
}

func (y *YINBOJIAN) GetSkillRange() []*Vec {
	vecs := make([]*Vec, 0)
	vec1 := &Vec{
		x: 1,
		y: 2,
		z: 3,
	}
	vec2 := &Vec{
		x: 1,
		y: 2,
		z: 3,
	}
	vecs = append(vecs, vec1)
	vecs = append(vecs, vec2)
	return vecs
}

func (y *YINBOJIAN) GetSkillTargetID() int32 {
	y.GetSkillRange()
	for _, v := range y.name {
		fmt.Println(v)
	}
	return 0
}
