package mnms

import (
	"sync"
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
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			t.Log("TestRunMqttSubscribe")
			subtoken := client.Subscribe("testtopic", 1, func(client MQTTClient.Client, msg MQTTClient.Message) {
				t.Log("Received message: ", msg.Payload(), msg.Topic())
			})
			subtoken.Wait()
			if subtoken.Error() != nil {
				t.Error(subtoken.Error())
			} else {
				t.Log("subscripe success")
			}
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			t.Log("TestRunMqttPublish")
			receipt := client.Publish("testtopic", 1, false, "test messages")
			receipt.Wait()
			if receipt.Error() != nil {
				t.Error(receipt.Error())
			} else {
				t.Log("publish success")
			}
		}()
		wg.Wait()
	}

}
