package task
import "time"

type Task struct{
	ID 			string
	Report_type string
	Parameters 	Parameters
	Status 		string
	Result 		Result
	Error 		Error
	Attempts 	int32
	Created_at 	time.Time
	Started_at 	time.Time
	Finished_at time.Time
}

type HealthStatus struct{
	Status 			 string
	Response_Time_Ms int
}

type Parameters struct{
	PeriodDays int
}

type PostInput struct{
	Parameters Parameters
	ReportType string
}

type Result struct{
	Summary string
}

type Error struct{
	Details *string
	Message string
}