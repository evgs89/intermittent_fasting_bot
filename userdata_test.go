package main

import (
	"testing"
	"time"

	"github.com/dchest/uniuri"
)

func createTestUser() *UserData {
	userId := uniuri.NewLen(20)
	chatId := uniuri.NewLen(20)
	Language := "ru-ru"
	tz := "+03:00"
	return NewUser(userId, chatId, Language, tz)
}

func TestCreateUser(t *testing.T) {
	u := createTestUser()
	if u.UserId == "" {
		t.Error("Userdata creation error")
	}
}

func TestUserCanGetAchievement(t *testing.T) {
	user := createTestUser()
	user.FastedDays = 24
	user.FastedHours = 20
	user.NoBreakingPlanDays = 1
	warriorPlanDay := []*fastingPhase{
		{Duration: 20, State: FullFasting},
		{Duration: 4, State: NoFasting},
	}
	fastingPlan := NewFastingPlan("Warrior's diet", warriorPlanDay)
	fastingPlan.FixedBeginningHour = 0
	warriorDayAchievement := Achievment{Name: "Warrior day!", DaysToGet: 1, FastingPlanToGet: fastingPlan}
	if warriorDayAchievement.checkUserDeserveThisAchievment(user) {
		t.Errorf("%v achievement obtained cause of error", warriorDayAchievement.Name)
	}
	fastingData := FastingData{
		Active:         true,
		FastingPlan:    fastingPlan,
		FastingStarted: time.Now().Unix(),
		FastedHours:    4,
		FastedDays:     1,
	}
	user.CurrentFastingData = &fastingData
	user.CurrentFastingData.FastingPlan = fastingPlan
	if !warriorDayAchievement.checkUserDeserveThisAchievment(user) {
		t.Errorf("%v achievement can't be obtained cause of error", warriorDayAchievement.Name)
	}
	standartDay := []*fastingPhase{
		{Duration: 16, State: FullFasting},
		{Duration: 8, State: NoFasting},
	}
	standartPlan := NewFastingPlan("16/8", standartDay)
	user.CurrentFastingData.FastingPlan = standartPlan
	standartPlanAchievement := Achievment{Name: "Fasting week!", WeeksToGet: 1}
	if standartPlanAchievement.checkUserDeserveThisAchievment(user) {
		t.Errorf("%v achievement obtained by error!", standartPlanAchievement.Name)
	}
	user.NoBreakingPlanDays = 8
	if !standartPlanAchievement.checkUserDeserveThisAchievment(user) {
		t.Errorf("%v achievement can't be obtained cause of error", standartPlanAchievement.Name)
	}
}
