package audit

import "context"

// Recorder 负责记录业务审计事件。
type Recorder interface {

	// RecordSuccess 记录成功的业务操作。
	RecordSuccess(
		ctx context.Context,
		event Event,
	) error

	// RecordFailure 记录失败的业务操作。
	RecordFailure(
		ctx context.Context,
		event Event,
		err error,
	) error
}
