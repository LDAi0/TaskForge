package task
import "time"

type Task struct{
	ID 			string
	Report_type string
	Parameters 	Parameters
	Status 		string
	Result 		Result
	Error 		string
	Attempts 	int32
	Created_at 	time.Time
	Started_at 	time.Time
	Finished_at time.Time
}

type HealthStatus struct{
	Status 			 string
	Response_Time_Ms int32
}

type Parameters struct{
	PeriodDays int16
}

type PostInput struct{
	ReportType string
	Parameters Parameters
}

type Result struct{

}