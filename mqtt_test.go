package mnms

import (
	"testing"
	"time"

	MQTTClient "github.com/eclipse/paho.mqtt.golang"
	MQTTBroker "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/listeners"
)

func TestRunMqttBroker(t *testing.T) {
	broker := MQTTBroker.NewServer(nil)
	tcp := listeners.NewTCP("t1", QC.MqttBrokerAddr)
	err := broker.AddListener(tcp, nil)
	if err != nil {
		t.Error(err)
	}
	err = broker.Serve()
	if err != nil {
		t.Error(err)
		return
	}
	defer broker.Close()
}

func TestRunMqttClient(t *testing.T) {
	go func() {
		_ = RunMqttBroker("TestRunMqttClient")
	}()
	time.Sleep(500 * time.Millisecond)

	opts := MQTTClient.NewClientOptions().AddBroker(QC.MqttBrokerAddr).SetClientID("testmqtt")
	opts.SetKeepAlive(keepAlive * time.Second)
	opts.SetPingTimeout(pingTimeout * time.Second)
	client := MQTTClient.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
	} else {
		t.Log("subscribe")

		err := RunMqttSubscribe(":11883", "testtopic1")
		if err != nil {
			t.Error(err)
		}
		err = RunMqttSubscribe(":11883", "testtopic2")
		if err != nil {
			t.Error(err)
		}
		err = RunMqttPublish(":11883", "testtopic1", "test messages1")
		if err != nil {
			t.Error(err)
		}
		err = RunMqttPublish(":11883", "testtopic2", "test messages2")
		if err != nil {
			t.Error(err)
		}
		time.Sleep(250 * time.Millisecond)
		err = RunMqttUnSubscribe(":11883", "testtopic1")
		if err != nil {
			t.Error(err)
		}
		err = RunMqttUnSubscribe(":11883", "testtopic2")
		if err != nil {
			t.Error(err)
		}
		err = RunMqttPublish(":11883", "testtopic1", "test messages1")
		if err != nil {
			t.Error(err)
		}
		err = RunMqttPublish(":11883", "testtopic2", "test messages2")
		if err != nil {
			t.Error(err)
		}
		time.Sleep(250 * time.Millisecond)

	}

}
