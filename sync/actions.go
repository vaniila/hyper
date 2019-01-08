package sync

// action codes
const (
	ActionUnknown int32 = iota
	ActionSubscribe
	ActionSubscribeSuccessful
	ActionSubscribeFailure
	ActionUnsubscribe
	ActionUnsubscribeSuccessful
	ActionUnsubscribeFailure
	ActionMessage
	ActionMessageSuccessful
	ActionMessageFailure
)
