package repositories

type RetryQueueRepository interface {
	EnqueueForRetry(taskID string, retryAfterSeconds int) error
	GetTasksReadyForRetry(limit int) ([]string, error)
	RemoveFromRetryQueue(taskID string) error
	CountPendingRetries() (int64, error)
	IsInRetryQueue(taskID string) (bool, error)
}
