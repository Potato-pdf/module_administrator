package valueobjects

type TaskStatus string

const (
	Pending    TaskStatus = "PENDING"
	Processing TaskStatus = "PROCESSING"
	Done       TaskStatus = "DONE"
	Failed     TaskStatus = "FAILED"
)

var ValidTaskStatus = []TaskStatus{
	Pending,
	Processing,
	Done,
	Failed,
}
