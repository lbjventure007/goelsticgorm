package inits

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"

	"log"
)

func InitSentinel() {
	err := sentinel.InitWithConfigFile("./sentinel.yaml")
	//err := sentinel.InitDefault()
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
	}
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource: "some-test",
			//MetricType:      flow.QPS,

			Threshold:              1,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			//StatIntervalInMs:       1000,
		},
	})
	if err != nil {
		log.Fatalf("eror:%+v", err)
		return
	}

}
