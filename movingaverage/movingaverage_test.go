//
// +build unit

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
	"bytes"
	"encoding/gob"
	"math/rand"
	"testing"
	"time"

	"github.com/intelsdi-x/pulse/control/plugin"
	"github.com/intelsdi-x/pulse/control/plugin/cpolicy"
	"github.com/intelsdi-x/pulse/core/ctypes"
	. "github.com/smartystreets/goconvey/convey"
)

//Random number generator
func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func TestMovingAverageProcessor(t *testing.T) {
	meta := Meta()
	Convey("Meta should return metadata for the plugin", t, func() {
		Convey("So meta.Name should equal movingaverage", func() {
			So(meta.Name, ShouldEqual, "movingaverage")
		})
		Convey("So meta.Version should equal 3", func() {
			So(meta.Version, ShouldEqual, 3)
		})
		Convey("So meta.Type should be of type plugin.ProcessorPluginType", func() {
			So(meta.Type, ShouldResemble, plugin.ProcessorPluginType)
		})
	})

	proc := NewMovingaverageProcessor()
	Convey("Create Movingaverage processor", t, func() {
		Convey("So proc should not be nil", func() {
			So(proc, ShouldNotBeNil)
		})
		Convey("So proc should be of type movingAverageProcessor", func() {
			So(proc, ShouldHaveSameTypeAs, &movingAverageProcessor{})
		})
		Convey("proc.GetConfigPolicy should return a config policy", func() {
			configPolicy, _ := proc.GetConfigPolicy()
			Convey("So config policy should be a cpolicy.ConfigPolicy", func() {
				So(configPolicy, ShouldHaveSameTypeAs, &cpolicy.ConfigPolicy{})
			})
			testConfig := make(map[string]ctypes.ConfigValue)
			testConfig["MovingAvgBufLength"] = ctypes.ConfigValueInt{Value: 10}
			cfg, errs := configPolicy.Get([]string{""}).Process(testConfig)
			Convey("So config policy should process testConfig and return a config", func() {
				So(cfg, ShouldNotBeNil)
			})
			Convey("So testConfig processing should return no errors", func() {
				So(errs.HasErrors(), ShouldBeFalse)
			})
		})
	})
}

func TestMovingAverageProcessorMetrics(t *testing.T) {
	Convey("Moving Average Processor tests", t, func() {
		metrics := make([]plugin.PluginMetricType, 10)
		config := make(map[string]ctypes.ConfigValue)

		config["MovingAvgBufLength"] = ctypes.ConfigValueInt{Value: -1}

		Convey("Moving average for int data", func() {
			for i, _ := range metrics {
				time.Sleep(3)
				rand.Seed(time.Now().UTC().UnixNano())
				data := randInt(65, 90)
				metrics[i] = *plugin.NewPluginMetricType([]string{"foo", "bar"}, time.Now(), "", data)
			}
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			enc.Encode(metrics)
			movingAverageObj := NewMovingaverageProcessor()

			_, received_data, _ := movingAverageObj.Process("pulse.gob", buf.Bytes(), config)

			var metrics_new []plugin.PluginMetricType

			//Decodes the content into pluginMetricType
			dec := gob.NewDecoder(bytes.NewBuffer(received_data))
			dec.Decode(&metrics_new)
			So(metrics, ShouldNotResemble, metrics_new)

		})

		Convey("Moving average for float32 data", func() {
			config["MovingAvgBufLength"] = ctypes.ConfigValueInt{Value: 40}
			for i, _ := range metrics {
				time.Sleep(3)
				rand.Seed(time.Now().UTC().UnixNano())
				data := randInt(65, 90)
				metrics[i] = *plugin.NewPluginMetricType([]string{"foo", "bar"}, time.Now(), "", float32(data))
			}
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			enc.Encode(metrics)

			movingAverageObj := NewMovingaverageProcessor()

			_, received_data, _ := movingAverageObj.Process("pulse.gob", buf.Bytes(), config)

			var metrics_new []plugin.PluginMetricType

			//Decodes the content into pluginMetricType
			dec := gob.NewDecoder(bytes.NewBuffer(received_data))
			dec.Decode(&metrics_new)
			So(metrics, ShouldNotResemble, metrics_new)

		})
		Convey("Moving average for float64 data", func() {
			for i, _ := range metrics {
				time.Sleep(3)
				rand.Seed(time.Now().UTC().UnixNano())
				data := randInt(65, 90)
				metrics[i] = *plugin.NewPluginMetricType([]string{"foo", "bar"}, time.Now(), "", float64(data))
			}
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			enc.Encode(metrics)

			movingAverageObj := NewMovingaverageProcessor()

			_, received_data, _ := movingAverageObj.Process("pulse.gob", buf.Bytes(), nil)

			var metrics_new []plugin.PluginMetricType

			//Decodes the content into pluginMetricType
			dec := gob.NewDecoder(bytes.NewBuffer(received_data))
			dec.Decode(&metrics_new)
			So(metrics, ShouldNotResemble, metrics_new)

		})

		Convey("Moving average for uint32 data", func() {
			for i, _ := range metrics {
				time.Sleep(3)
				rand.Seed(time.Now().UTC().UnixNano())
				data := randInt(65, 90)
				metrics[i] = *plugin.NewPluginMetricType([]string{"foo", "bar"}, time.Now(), "", uint32(data))
			}
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			enc.Encode(metrics)
			movingAverageObj := NewMovingaverageProcessor()

			_, received_data, _ := movingAverageObj.Process("pulse.gob", buf.Bytes(), nil)

			var metrics_new []plugin.PluginMetricType

			//Decodes the content into pluginMetricType
			dec := gob.NewDecoder(bytes.NewBuffer(received_data))
			dec.Decode(&metrics_new)
			So(metrics, ShouldNotResemble, metrics_new)

		})

		Convey("Moving average for uint64 data", func() {
			for i, _ := range metrics {
				time.Sleep(3)
				rand.Seed(time.Now().UTC().UnixNano())
				data := randInt(65, 90)
				metrics[i] = *plugin.NewPluginMetricType([]string{"foo", "bar"}, time.Now(), "", uint64(data))
			}
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			enc.Encode(metrics)
			movingAverageObj := NewMovingaverageProcessor()

			_, received_data, _ := movingAverageObj.Process("pulse.gob", buf.Bytes(), nil)

			var metrics_new []plugin.PluginMetricType

			//Decodes the content into pluginMetricType
			dec := gob.NewDecoder(bytes.NewBuffer(received_data))
			dec.Decode(&metrics_new)
			So(metrics, ShouldNotResemble, metrics_new)

		})

		Convey("Moving average for unknown data type", func() {
			for i, _ := range metrics {

				data := "I am an unknow data Type"
				metrics[i] = *plugin.NewPluginMetricType([]string{"foo", "bar"}, time.Now(), "", data)
			}
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			enc.Encode(metrics)

			movingAverageObj := NewMovingaverageProcessor()

			_, received_data, _ := movingAverageObj.Process("pulse.gob", buf.Bytes(), nil)

			var metrics_new []plugin.PluginMetricType

			//Decodes the content into pluginMetricType
			dec := gob.NewDecoder(bytes.NewBuffer(received_data))
			dec.Decode(&metrics_new)
			So(metrics, ShouldNotResemble, metrics_new)

		})

	})
}
