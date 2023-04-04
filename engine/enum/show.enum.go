package enum

type ShowTaskStatusID = uint32

var ShowTaskStatus = struct {
	Absent     ShowTaskStatusID
	Monitoring ShowTaskStatusID
	Recording  ShowTaskStatusID
}{
	Absent:     0,
	Monitoring: 1,
	Recording:  2,
}
