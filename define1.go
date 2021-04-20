package main

// 技能状态触发条件
type SkillBuffTrigger int

const (
	// 技能状态触发条件
	SkillBuffTrigger_NONE                             SkillBuffTrigger = iota //无
	SkillBuffTrigger_SKILLSELECT                                              //技能选择时
	SkillBuffTrigger_SKILLACTIVE                                              //技能释放时
	SkillBuffTrigger_BEHIT_NORMALATTACK                                       //被攻击时
	SkillBuffTrigger_HIT_NORMALATTACK                                         //击中目标时
	SkillBuffTrigger_ENDTURN                                                  //战斗回合结束时
	SkillBuffTrigger_BLOCK                                                    //战斗阻挡时
	SkillBuffTrigger_DODGE                                                    //战斗闪避时
	SkillBuffTrigger_START                                                    //人物登场时
	SkillBuffTrigger_MOVEDONE                                                 //结束移动时
	SkillBuffTrigger_NEWROUND                                                 //新ROUND开始
	SkillBuffTrigger_BEFORE_SKILL_DAMAGE_COUNT                                //技能伤害结算前
	SkillBuffTrigger_AFTER_SKILL_DAMAGE_COUNT                                 //技能伤害结算后
	SkillBuffTrigger_BEHIT_SKILLATTACK                                        //被技能攻击时
	SkillBuffTrigger_HIT_SKILLATTACK                                          //技能击中目标时
	SkillBuffTrigger_BEFORE_NORMALATTACK_DAMAGE_COUNT                         //普攻伤害结算前
	SkillBuffTrigger_AFTER_NORMALATTACK_DAMAGE_COUNT                          //普攻伤害结算后
	SkillBuffTrigger_PASSIVE_SKILL_ENABLE
	SkillBuffTrigger_PASSIVE_SKILL_DISABLE
	SkillBuffTrigger_SKILL_DONE
	SkillBuffTrigger_ENTER_OVERWATCH
	SkillBuffTrigger_LOCKSHOOT_DONE
	SkillBuffTrigger_LOCKSHOOT_OR_SKILL_DONE
)
