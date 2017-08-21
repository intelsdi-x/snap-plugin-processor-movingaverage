//
// +build small

/*
http://www.apache.org/licenses/LICENSE-2.0.txt

Copyright 2015 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package movingaverage

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

func TestMovingAverageProcessor(t *testing.T) {
	proc := New()
	Convey("Create Movingaverage processor", t, func() {
		Convey("So proc should not be nil", func() {
			So(proc, ShouldNotBeNil)
		})
		Convey("So proc should be of type MAProcessor", func() {
			So(proc, ShouldHaveSameTypeAs, &MAProcessor{})
		})
		Convey("proc.GetConfigPolicy should return a config policy", func() {
			configPolicy, err := proc.GetConfigPolicy()
			Convey("So there should be no error returned", func() {
				So(err, ShouldBeNil)
			})
			Convey("So config policy should not be nil", func() {
				So(configPolicy, ShouldNotBeNil)
			})
			Convey("So config policy should be a plugin.ConfigPolicy", func() {
				So(configPolicy, ShouldHaveSameTypeAs, plugin.ConfigPolicy{})
			})
		})
	})
}

func TestMovingAverageProcessorMetrics(t *testing.T) {
	Convey("Moving Average Processor tests", t, func() {
		data := []int{1, 10, 20, 30, 50, 10, 20, 30}
		averages := []float64{1, 11.0 / 2, 31.0 / 3, 61.0 / 4, 111.0 / 5, 120.0 / 5, 130.0 / 5, 140.0 / 5}
		emptySlice := make([]float64, len(data))
		config := plugin.Config{
			"MovingAvgBufLength": int64(5),
		}

		Convey("Moving average for int data", func() {
			metrics := make([]plugin.Metric, len(data))
			for i := range metrics {
				metrics[i] = plugin.Metric{
					Namespace: plugin.NewNamespace("foo", "bar"),
					Timestamp: time.Now(),
					Data:      data[i],
				}
			}
			proc := New()
			mts, err := proc.Process(metrics, config)
			So(err, ShouldBeNil)
			So(mts, ShouldHaveLength, len(data))
			values := make([]float64, len(data))
			for i := range values {
				values[i] = mts[i].Data.(float64)
			}
			So(values, ShouldResemble, averages)
		})

		Convey("Moving average for unknown data type", func() {
			metrics := make([]plugin.Metric, len(data))
			for i := range metrics {
				metrics[i] = plugin.Metric{
					Namespace: plugin.NewNamespace("foo", "bar"),
					Timestamp: time.Now(),
					Data:      "some string",
				}
			}
			proc := New()
			mts, err := proc.Process(metrics, config)
			So(err, ShouldBeNil)
			So(mts, ShouldHaveLength, len(data))
			values := make([]float64, len(data))
			for i := range values {
				values[i] = mts[i].Data.(float64)
			}
			So(values, ShouldResemble, emptySlice)
		})
	})
}
