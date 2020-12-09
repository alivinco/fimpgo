package transport

import (
	"github.com/futurehomeno/fimpgo"
	"testing"
	log "github.com/sirupsen/logrus"
	"time"
)

func TestBufferedStream_serializeBuffer(t *testing.T) {
	type fields struct {
		bufferMaxSize       int
		bufferInterval      time.Duration
		flushToFile         bool
		compressBeforeFlush bool
		filePrefix          string
		flushToSinkChannel  bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "test1", fields: fields{
			bufferMaxSize:       20,
			bufferInterval:      600,
			flushToFile:         true,
			compressBeforeFlush: true,
			filePrefix:          "test",
			flushToSinkChannel:  false,
		}, wantErr:false },
	}

	log.SetLevel(log.DebugLevel)
	mqtt := fimpgo.NewMqttTransport("tcp://cube.local:1884","fimpgotest","","",true,1,1)
	err := mqtt.Start()
	t.Log("Connected")
	time.Sleep(time.Second*1)
	mqtt.Subscribe("pt:j1/mt:evt/rt:dev/+/ad:1/sv:meter_elec/+")
	if err != nil {
		t.Error("Error connecting to broker ",err)
	}
	chan1 := make(fimpgo.MessageCh)
	mqtt.RegisterChannel("chan1",chan1)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			su := NewBufferedStream(tt.fields.bufferMaxSize,tt.fields.bufferInterval,tt.fields.compressBeforeFlush)
			su.SetSourceStream(chan1)
			su.ConfigureFileSink("test1",".")
			sink := su.ConfigureChanelSink(5)
			r := <- sink
			if len(r)==0 {
				t.Error("Empty buffer size")
			}
		})
	}
	mqtt.Stop()
}
