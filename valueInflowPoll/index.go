package valueinflowpoll

import (
	"usepolymer.co/background/constants"
	"usepolymer.co/background/logger"
	"usepolymer.co/background/poller"
)


func PollForValueInflow() {
	logger.Info("beginning value in poll")
	poller.BeginCronPoll(constants.FLW_INFLOW_POLLING_INTERVAL_MINS, pollFlutterwave)
}