package common

import "errors"

var (
	JOB_LOCK_ERROR = errors.New("任务锁已被占用")
)
