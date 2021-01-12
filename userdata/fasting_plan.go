package userdata

import "log"

var FastingPlans map[int]*FastingPlan

func LoadData() {
	if FastingPlans == nil || len(FastingPlans) == 0 {
		FastingPlans = GetFastingPlans()
	}
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

type FastingPhase struct {
	Duration          int
	State             FastingState
	Description       string
	NotifyAtBeginning bool
	Picture           string
}

func GetFastingPhases(fastingPlanId int) []*FastingPhase {
	rows, err := DBCONN.Query("SELECT Duration, Description, FastingState, NotifyAtBeginning, Picture FROM FastingPhase WHERE FastingPlan=? ORDER BY OrderNum;", fastingPlanId)
	if err != nil {
		log.Fatal(err)
	}
	phases := make([]*FastingPhase, 0)
	for rows.Next() {
		phase := FastingPhase{}
		err := rows.Scan(&phase.Duration, &phase.Description, &phase.State, &phase.NotifyAtBeginning, &phase.Picture)
		if err != nil {
			log.Fatal(err)
		}
		phases = append(phases, &phase)
	}
	return phases
}

type FastingPlan struct {
	Name            string
	Description     string
	Phases          []*FastingPhase
	Duration        int
	Level           FastingPlanHardnessLevel
	GraphicsPicture string
}

func NewFastingPlan(name string, phases []*FastingPhase) *FastingPlan {
	duration := 0
	for _, fp := range phases {
		duration += fp.Duration
	}
	return &FastingPlan{
		Name:     name,
		Phases:   phases,
		Duration: duration,
	}
}

func GetFastingPlans() map[int]*FastingPlan {
	rows, err := DBCONN.Query("SELECT id, Name, Description, Duration, Level, Picture FROM FastingPlan;")
	if err != nil {
		log.Fatal(err)
	}
	FastingPlans := make(map[int]*FastingPlan)
	for rows.Next() {
		var id int
		fp := FastingPlan{}
		err := rows.Scan(id, &fp.Name, &fp.Description, &fp.Duration, &fp.Level, &fp.GraphicsPicture)
		if err != nil {
			log.Fatal(err)
		}
		phases := GetFastingPhases(id)
		fp.Phases = phases
		FastingPlans[id] = &fp
	}
	return FastingPlans
}

type FastingData struct {
	Active          bool
	FastingPlan     *FastingPlan
	CurrentPhase    *FastingPhase
	PhaseStartedAt  int64
	PhaseEndsAt     int64
	RemindersWorkId string
}

func GetFastingDataFromDB(userid string) *FastingData {
	LoadData()
	rows, err := DBCONN.Query("SELECT FastingPlan, OrderNum, StartAt, EndAt, WorkId FROM FastingsView WHERE UNIX_TIMESTAMP() BETWEEN StartAt AND EndAt AND UserId=?;", userid)
	if err != nil {
		log.Fatal(err)
	}
	fd := FastingData{
		Active:       false,
		FastingPlan:  nil,
		CurrentPhase: nil,
	}
	for rows.Next() {
		fd.Active = true
		var fpId, fphNum int
		err = rows.Scan(fpId, fphNum, &fd.PhaseStartedAt, &fd.PhaseEndsAt, &fd.RemindersWorkId)
		if err != nil {
			log.Fatal(err)
		}
		fd.FastingPlan = FastingPlans[fpId]
		fd.CurrentPhase = FastingPlans[fpId].Phases[fphNum]
		return &fd
	}
	return nil
}
