package requestMonitor

import (
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

//fmt.Println(proto.MarshalTextString(metric))
func ExampleRequestMonitor() {
	requestMonitor := New("test")

	for i := 0; i < 1000; i++ {
		requestMonitor.AddRequest("Success")
	}

	counter, err := requestMonitor.counter.GetMetricWithLabelValues("Success")
	if err != nil {
		panic("No counter object")
	}

	metric := &dto.Metric{}
	counter.Write(metric)

	// TP50
	for i := 0; i < 500; i++ {
		requestMonitor.RecordRequest("Success", time.Duration(100)*time.Millisecond)
	}
	// TP75
	for i := 0; i < 250; i++ {
		requestMonitor.RecordRequest("Success", time.Duration(200)*time.Millisecond)
	}
	// TP90
	for i := 0; i < 200; i++ {
		requestMonitor.RecordRequest("Success", time.Duration(300)*time.Millisecond)
	}
	// TP99
	for i := 0; i < 50; i++ {
		requestMonitor.RecordRequest("Success", time.Duration(500)*time.Millisecond)
	}

	reg := prometheus.NewRegistry()
	reg.MustRegister(requestMonitor.duration)
	metricFamilies, err := reg.Gather()
	if err != nil || len(metricFamilies) != 1 {
		panic("unexpected behavior of custom test registry")
	}

	fmt.Println(proto.MarshalTextString(metric))
	fmt.Println(proto.MarshalTextString(metricFamilies[0]))
	/**
	Output:
	label: <
	  name: "status"
	  value: "Success"
	>
	counter: <
	  value: 1000
	>

	name: "test_duration_milliseconds"
	help: "Summary the duration time of test request"
	type: SUMMARYSUMMARY
	metric: <
	  label: <
	    name: "status"
	    value: "Success"
	  >
	  summary: <
	    sample_count: 1000
	    sample_sum: 185000
	    quantile: <
	      quantile: 0.5
	      value: 100
	    >
	    quantile: <
	      quantile: 0.75
	      value: 200
	    >
	    quantile: <
	      quantile: 0.9
	      value: 300
	    >
	    quantile: <
	      quantile: 0.99
	      value: 500
	    >
	  >
	>
	*/
}
