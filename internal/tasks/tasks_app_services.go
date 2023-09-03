package tasks

type TaskServices interface {
	// Retrieve all the uncompleted tasks
	GetAllUncompleted(taskGateway TaskProvider) ([]*Task, error)
	MarkAsCompleted([]*Task) error
}

// Task providers can be Microsoft, Google, etc.
// Microsoft provides tasks through the Microsoft Graph API and Google through the Google Tasks API
// Each of these providers will have their own implementation of the TaskProvider interface
type TaskProvider interface {
	// Initialise the task provider with its specific configuration
	Init(options ...func(*TaskProvider)) *TaskProvider

	GetTasks(status string) ([]*Task, error)
}

type MicrosoftTaskProvider struct{}

// func (m *TaskProvider) Init(options ...func(*TaskProvider)) *TaskProvider {
// 	provider := &MicrosoftTaskProvider{}
// 	for _, option := range options {
// 		option(*provider)
// 	}
// 	return provider
// }

// func (m *MicrosoftTaskProvider) GetTasks(status string) ([]*Task, error) {
// 	return []*Task{}, nil
// }
