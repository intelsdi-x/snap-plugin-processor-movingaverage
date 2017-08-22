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
	"fmt"
	"reflect"

	log "github.com/Sirupsen/logrus"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

var (
	floatType = reflect.TypeOf(float64(0))
)

const (
	// Name of the plugin
	Name = "movingaverage"
	// Version of the plugin
	Version = 8

	defaultBufferLength = 10
)

type average struct {
	// buffer with past values
	buffer []float64
	// index of current element
	counter int
	// indicator if our buffer was already filled
	full bool
}

// MAProcessor struct hold map with average object for each unique namespace passed to plugin
type MAProcessor struct {
	movingAverageMap map[string]*average
	bufferLength     int
	initialized      bool
}

// New returns initialized plugin instance
func New() *MAProcessor {
	return &MAProcessor{
		movingAverageMap: make(map[string]*average),
		bufferLength:     defaultBufferLength,
		initialized:      false,
	}
}

func newAverage(length int) *average {
	avg := average{
		buffer:  make([]float64, length),
		counter: 0,
		full:    false,
	}
	return &avg
}

// castToFloat is used to convert data hidden under interface{} to float64
// reflect is used, so that any compatible type can be cast
func castToFloat(data interface{}) (float64, error) {
	val := reflect.ValueOf(data)
	val = reflect.Indirect(val)
	if !val.Type().ConvertibleTo(floatType) {
		return 0, fmt.Errorf("Cannot convert %v to float64", val.Type())
	}
	return val.Convert(floatType).Float(), nil
}

func (p *MAProcessor) calculateMovingAverage(m plugin.Metric) (float64, error) {
	namespace := m.Namespace.String()
	value, err := castToFloat(m.Data)
	if err != nil {
		return 0, err
	}
	avg, ok := p.movingAverageMap[namespace]
	if !ok {
		p.movingAverageMap[namespace] = newAverage(p.bufferLength)
		avg = p.movingAverageMap[namespace]
	}
	counter := avg.counter % p.bufferLength
	avg.buffer[counter] = value

	sum := float64(0)
	for _, data := range avg.buffer {
		sum += data
	}

	count := p.bufferLength
	// Divide sum by counter instead if buffer is not full yet
	if !avg.full {
		count = avg.counter + 1
	}
	movingAvg := sum / float64(count)

	counter++
	if counter == p.bufferLength {
		avg.full = true
	}
	avg.counter = counter
	return movingAvg, nil
}

// GetConfigPolicy returns plugin's config policy
func (p *MAProcessor) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	cp := plugin.NewConfigPolicy()
	cp.AddNewIntRule([]string{""}, "MovingAvgBufLength", false,
		plugin.SetDefaultInt(defaultBufferLength), plugin.SetMinInt(1))
	return *cp, nil
}

// Process processes metrics
func (p *MAProcessor) Process(metrics []plugin.Metric, config plugin.Config) ([]plugin.Metric, error) {
	logger := log.New()
	logger.Debug("movingAverage Processor started")

	if !p.initialized {
		bufLen, err := config.GetInt("MovingAvgBufLength")
		if err == nil {
			if bufLen >= 1 {
				logger.Debugf("Using buffer length: %v\n", bufLen)
				p.bufferLength = int(bufLen)
			}
		}
		p.initialized = true
	}

	for i, m := range metrics {
		logger.Debugf("Data received %v", m.Data)
		var err error
		metrics[i].Data, err = p.calculateMovingAverage(m)
		if err != nil {
			logger.Warnf(err.Error())
		}
		logger.Debugf("Moving Average: %v", metrics[i].Data)
	}
	return metrics, nil
}
