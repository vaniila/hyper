package sync

import "github.com/samuelngs/hyper/fault"

// sync errors
var (
	InvalidPacket     = fault.Format("packet is Invalid")
	InvalidAction     = fault.Format("action [%d] is Invalid")
	NamespaceNotExist = fault.Format("namespace %s does not exist")

	ChannelUnauthorized      = fault.Format("no access permission to `%s:%s`")
	ChannelNotExist          = fault.Format("channel %s does not exist")
	ChannelAlreadySubscribed = fault.Format("`%s:%s` has already been subscribed")
	ChannelNotSubscribed     = fault.Format("`%s:%s` has not been subscribed")
)
