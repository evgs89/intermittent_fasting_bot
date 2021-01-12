package userdata

import (
	"testing"
	"time"

	"github.com/dchest/uniuri"
)

func createTestUser() *UserData {
	userId := uniuri.NewLen(20)
	chatId := uniuri.NewLen(20)
	tz := "+03:00"
	return NewUser(userId, chatId, tz)
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
	user.FastingStartedAt = time.Now().AddDate(0, 0, -1).Unix()
	warriorPlanDay := []*FastingPhase{
		{Duration: 20, State: FullFasting},
		{Duration: 4, State: NoFasting},
	}
	fastingPlan := NewFastingPlan("Warrior's diet", warriorPlanDay)
	warriorDayAchievement := Achievement{Name: "Warrior day!", DaysToGet: 1, FastingPlanToGet: fastingPlan}
	if warriorDayAchievement.checkUserDeserveThisAchievement(user) {
		t.Errorf("%v achievement obtained cause of error", warriorDayAchievement.Name)
	}
	fastingData := FastingData{
		Active:      true,
		FastingPlan: fastingPlan,
		FastedHours: 4,
		FastedDays:  1,
	}
	user.CurrentFastingData = &fastingData
	user.CurrentFastingData.FastingPlan = fastingPlan
	if !warriorDayAchievement.checkUserDeserveThisAchievement(user) {
		t.Errorf("%v achievement can't be obtained cause of error", warriorDayAchievement.Name)
	}
	standardDay := []*FastingPhase{
		{Duration: 16, State: FullFasting},
		{Duration: 8, State: NoFasting},
	}
	standardPlan := NewFastingPlan("16/8", standardDay)
	user.CurrentFastingData.FastingPlan = standardPlan
	standardPlanAchievement := Achievement{Name: "Fasting week!", WeeksToGet: 1}
	if standardPlanAchievement.checkUserDeserveThisAchievement(user) {
		t.Errorf("%v achievement obtained by error!", standardPlanAchievement.Name)
	}
	user.FastingStartedAt = time.Now().AddDate(0, 0, -8).Unix()
	if !standardPlanAchievement.checkUserDeserveThisAchievement(user) {
		t.Errorf("%v achievement can't be obtained cause of error", standardPlanAchievement.Name)
	}
}
