package userdata

import "log"

type Achievement struct {
	Name             string
	Icon             string
	Description      string
	HoursToGet       int
	DaysToGet        int
	WeeksToGet       int
	FastingPlanToGet *FastingPlan
}

func (a *Achievement) checkUserDeserveThisAchievement(ud *UserData) bool {
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
	if a.WeeksToGet > 0 && ud.GetNoBreakingFastingPlanDays() >= a.WeeksToGet*7 {
		return true
	}
	return false
}

func GetAchievementsList() map[int]*Achievement {
	rows, err := DBCONN.Query("SELECT id, Name, Icon, Description, HoursToGet, DaysToGet, FastingPlanToGet FROM Achievements")
	if err != nil {
		log.Fatal(err)
	}
	achievements := make(map[int]*Achievement)
	for rows.Next() {
		var id int
		a := Achievement{}
		err = rows.Scan(id, &a.Name, &a.Icon, &a.Description, &a.HoursToGet, &a.DaysToGet, &a.FastingPlanToGet)
		if err != nil {
			log.Fatal(err)
		}
		achievements[id] = &a
	}
	return achievements
}
