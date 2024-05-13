package valueoutflowpoll

import (
	"usepolymer.co/background/constants"
	"usepolymer.co/background/logger"
	"usepolymer.co/background/poller"
)



func PollForValueOutflow() {
	logger.Info("beginning value out poll")
	poller.BeginCronPoll(constants.CHIMONEY_OUTFLOW_POLLING_INTERVAL_MINS, pollChimoney)
	poller.BeginCronPoll(constants.FLW_OUTFLOW_POLLING_INTERVAL_MINS, pollFlutterwave)
}