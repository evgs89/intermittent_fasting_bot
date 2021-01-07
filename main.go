package main

import (
	"time"
)

var users = []*UserData{}

func main() {

	return
}

type FastingState int

const (
	NoFasting FastingState = iota
	LimitedCalories
	FullFasting
)

func (fs FastingState) String() string {
	return []string{"No Fasting", "Limited Calories", "Full Fasting"}[fs]
}

type FastingPlanHardnessLevel int

const (
	Beginner FastingPlanHardnessLevel = iota
	Intermediate
	Advanced
)

func (fpl FastingPlanHardnessLevel) String() string {
	return []string{"Beginner", "Intermediate", "Advanced"}[fpl]
}

type Achievment struct {
	Name             string
	Icon             string
	Description      string
	HoursToGet       int
	DaysToGet        int
	WeeksToGet       int
	FastingPlanToGet *FastingPlan
}

func (a *Achievment) checkUserDeserveThisAchievment(ud *UserData) bool {
	if a.FastingPlanToGet != nil {
		if ud.CurrentFastingData == nil || ud.CurrentFastingData.FastingPlan.Name != a.FastingPlanToGet.Name {
			return false
		}
		if ud.CurrentFastingData.FastedDays >= a.DaysToGet {
			return true
		}
	}
	if a.HoursToGet > 0 && ud.FastedHours >= a.HoursToGet {
		return true
	}
	if a.DaysToGet > 0 && ud.FastedDays >= a.HoursToGet {
		return true
	}
	if a.WeeksToGet > 0 && ud.NoBreakingPlanDays >= a.WeeksToGet*7 {
		return true
	}
	return false
}

type fastingPhase struct {
	Duration    int
	State       FastingState
	Description string
}

type FastingPlan struct {
	Name               string
	Description        string
	Phases             []*fastingPhase
	Duration           int
	FixedBeginningHour int
	Level              FastingPlanHardnessLevel
}

func NewFastingPlan(name string, phases []*fastingPhase) *FastingPlan {
	duration := 0
	for _, fp := range phases {
		duration += fp.Duration
	}
	return &FastingPlan{
		Name:               name,
		Phases:             phases,
		FixedBeginningHour: -1,
		Duration:           duration,
	}
}

type UserData struct {
	UserId             string
	Active             bool
	LastActive         int64
	FastingState       FastingState
	CurrentFastingData *FastingData
	Acievments         map[string]int64
	FastedHours        int
	FastedDays         int
	NoBreakingPlanDays int
	Settings           *UserSettings
}

func NewUser(user_id string, chat_id string, lang string, tz string) *UserData {
	userSettings := UserSettings{
		ChatId:   chat_id,
		Language: lang,
		Timezone: tz,
	}
	user := UserData{
		UserId:     user_id,
		Settings:   &userSettings,
		Active:     true,
		LastActive: time.Now().Unix(),
	}
	return &user
}

func (ud *UserData) CheckAchievments(achievments []*Achievment) {
	for _, achievment := range achievments {
		if ud.Acievments[achievment.Name] == 0 {
			if achievment.checkUserDeserveThisAchievment(ud) {
				ud.Acievments[achievment.Name] = time.Now().Unix()
			}
		}
	}
}

type FastingData struct {
	Active         bool
	FastingPlan    *FastingPlan
	FastingStarted int64
	FastedHours    int
	FastedDays     int
}

type UserSettings struct {
	ChatId      string
	Language    string
	PhoneNumber string
	Timezone    string
}
