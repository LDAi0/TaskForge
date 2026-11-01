package task
import "time"

type Task struct{
	ID 			string
	Report_type string
	Parameters 	[]byte
	Status 		string
	Result 		[]byte
	Error 		string
	Attempts 	int32
	Created_at 	time.Time
	Started_at 	time.Time
	Finished_at time.Time
}