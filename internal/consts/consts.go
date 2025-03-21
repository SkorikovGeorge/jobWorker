package consts

// Main errors
const (
	ErrStartingServer = "Error starting server"
	ErrShutdownServer = "Error shutting down"
)

// JSON errors
const (
	ErrEncodingJSON    = "Error encoding JSON"
	ErrClosingBodyJSON = "Error closing JSON body"
	ErrParsingJSON     = "Error parsing JSON"
)

// Redis errors
const (
	ErrConnectingRedis = "Error connecting to redis"
	ErrEnqueueRedis    = "Error adding new job in redis queue"
	ErrDequeueRedis    = "Error removing job from redis queue"
	ErrSetRedis        = "Error setting job status in redis queue"
	ErrGetRedis        = "Error getting job status from redis queue"
	ErrKeyNotFound     = "Error finding key in redis queue"
	ErrEmptyQueue      = "Dequeue on empty queue"
	ErrIncrementRedis  = "Error incrementing retries counter"
	ErrTTLRedis        = "Error setting job TTL"
)

// Job errors
const (
	ErrRunningJob = "Error running job"
	ErrSomething  = "Something went wrong"
	ErrPriority   = "Error. Priority should be from 1 to 10"
)

// Job statuses
const (
	StatusPending = "pending"
	StatusRunning = "running"
	StatusDone    = "done"
	StatusFailed  = "failed"
)
