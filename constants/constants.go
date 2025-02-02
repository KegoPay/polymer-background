package constants

import (
	"os"
	"strconv"
)

var FLW_INFLOW_POLLING_INTERVAL_MINS int
var FLW_OUTFLOW_POLLING_INTERVAL_MINS int
var CHIMONEY_OUTFLOW_POLLING_INTERVAL_MINS int

func InitialisePollingIntervals() {
	FLW_INFLOW_POLLING_INTERVAL_MINS, _ = strconv.Atoi(os.Getenv("FLW_INFLOW_POLLING_INTERVAL_MINS"))
	FLW_OUTFLOW_POLLING_INTERVAL_MINS, _ = strconv.Atoi(os.Getenv("FLW_OUTFLOW_POLLING_INTERVAL_MINS"))
	CHIMONEY_OUTFLOW_POLLING_INTERVAL_MINS, _ = strconv.Atoi(os.Getenv("CHIMONEY_OUTFLOW_POLLING_INTERVAL_MINS"))
}