package mnms

import (
	"errors"
	"sync"
	"time"

	MQTTClient "github.com/eclipse/paho.mqtt.golang"
	MQTTBroker "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/events"
	"github.com/mochi-co/mqtt/server/listeners"

	"github.com/qeof/q"
)

func init() {
}

const (
	keepAlive   = 60
	pingTimeout = 1
)

type MqttClient struct {
	client MQTTClient.Client
}

var mqttclient MqttClient

func RunMqttPublish(topicname, messages string) error {
	if mqttclient.client == nil {
		return errors.New("mqtt client not found")
	}
	receipt := mqttclient.client.Publish(topicname, 1, false, messages)
	receipt.Wait()
	return nil
}

func RunMqttSubscribe(topicname string, timeout int) error {
	if mqttclient.client == nil {
		return errors.New("mqtt client not found")
	}
	go func() {
		subtoken := mqttclient.client.Subscribe(topicname, 0, func(client MQTTClient.Client, msg MQTTClient.Message) {
			q.Q("Received message: ", msg.Payload(), msg.Topic())
			receivemsg := "topic: " + msg.Topic() + " message: " + string(msg.Payload()[:])
			syslogerr := SendSyslog(LOG_INFO, "mqttclient", receivemsg)
			if syslogerr != nil {
				q.Q(syslogerr)
			}
		})
		subtoken.Wait()
		for i := 0; i < timeout; i++ {
			q.Q(i)
			time.Sleep(pingTimeout * time.Second)
		}
		q.Q("Subscribe done")
	}()
	return nil
}

func RunMqttBroker(servername string) error {
	q.Q(QC.MqttBrokerAddr)
	broker := MQTTBroker.NewServer(nil)
	tcp := listeners.NewTCP("t1", QC.MqttBrokerAddr)
	err := broker.AddListener(tcp, nil)
	if err != nil {
		return err
	}

	broker.Events.OnMessage = func(client events.Client, pk events.Packet) (pkx events.Packet, err error) {
		q.Q("OnMessage : ", client.ID, pk.TopicName, pk.Payload)
		msg := "client id: " + client.ID + ", topic: " + pk.TopicName + ", message: " + string(pk.Payload[:])
		syslogerr := SendSyslog(LOG_INFO, "mqttbroker", msg)
		if syslogerr != nil {
			q.Q(syslogerr)
		}
		return pk, nil
	}

	err = broker.Serve()
	if err != nil {
		return err
	}
	defer broker.Close()

	// mqtt client init
	opts := MQTTClient.NewClientOptions().AddBroker(QC.MqttBrokerAddr).SetClientID(servername)
	opts.SetKeepAlive(keepAlive * time.Second)
	opts.SetPingTimeout(pingTimeout * time.Second)
	client := MQTTClient.NewClient(opts)
	token := client.Connect()
	if token.Wait() {
		err := token.Error()
		if err != nil {
			return err
		}
	}
	defer client.Disconnect(2)

	mqttclient = MqttClient{client: client}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			time.Sleep(pingTimeout * time.Second)
		}
	}()
	wg.Wait()
	return nil
}
