package userdata

import (
	"database/sql"
	"log"
	"time"
)

type DBConnector interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	Close()
	Connect() error
}

var DBCONN DBConnector

type FastingState int

const (
	NoFasting FastingState = iota
	LimitedCalories
	FullFasting
)

func (fs FastingState) String() string {
	return []string{"No Fasting", "Limited Calories", "Full Fasting"}[fs]
}

type UserData struct {
	UserId             string
	Active             bool
	LastActive         int64
	CurrentFastingData *FastingData
	Achievements       map[int]int64
	FastedHours        int
	FastedDays         int
	FastingStartedAt   int64
	IsPremium          bool
	Timezone           string
	ChatId             string
}

func NewUser(userId string, chatId string, tz string) *UserData {
	user := UserData{
		UserId:           userId,
		ChatId:           chatId,
		Timezone:         tz,
		Active:           true,
		LastActive:       time.Now().Unix(),
		FastingStartedAt: 0,
	}
	_, err := DBCONN.Exec("INSERT INTO Users (UserId, Active, LastActive, ChatId, Timezone) VALUES (?, ?, ?, ?, ?)", user.UserId, user.Active, user.LastActive, user.ChatId, user.Timezone)
	if err != nil {
		log.Print(err)
		return nil
	}
	return &user
}

func (ud *UserData) CheckAchievements(achievements map[int]*Achievement) {
	for id, achievement := range achievements {
		if ud.Achievements[id] == 0 {
			if achievement.checkUserDeserveThisAchievement(ud) {
				ud.Achievements[id] = time.Now().Unix()
			}
		}
	}
}

func (ud *UserData) GetNoBreakingFastingPlanDays() int {
	secondsFasting := time.Now().Unix() - ud.FastingStartedAt
	day, _ := time.ParseDuration("1d")
	return int(secondsFasting / int64(day.Seconds()))
}

func (ud *UserData) GetUserAchievements() map[int]int64 {
	rows, err := DBCONN.Query("SELECT AchievementId, GotAt FROM AchievementsGot WHERE UserId=?;", ud.UserId)
	if err != nil {
		log.Fatal(err)
	}
	userAchievements := make(map[int]int64)
	for rows.Next() {
		var id int
		var gotAt int64
		err = rows.Scan(id, gotAt)
		if err != nil {
			log.Fatal(err)
		}
		userAchievements[id] = gotAt
	}
	return userAchievements
}

func (ud *UserData) GetFastingState() FastingState {
	if ud.CurrentFastingData == nil {
		return NoFasting
	}
	return ud.CurrentFastingData.CurrentPhase.State
}
