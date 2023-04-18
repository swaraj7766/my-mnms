package mnms

import (
	"errors"
	"strings"
	"sync"
	"time"

	MQTTClient "github.com/eclipse/paho.mqtt.golang"
	MQTTBroker "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/events"
	"github.com/mochi-co/mqtt/server/listeners"

	"github.com/qeof/q"
)

const (
	keepAlive     = 60
	pingTimeout   = 1
	clientTimeout = 5
)

type SubscribeStatus struct {
	inSubscribe bool
}

type MqttClient struct {
	clientHandle  MQTTClient.Client
	tcp           string
	isRun         bool
	subscribeList map[string]SubscribeStatus
}

type MqttClientList struct {
	client map[string]MqttClient
}

var mqttclient MqttClientList

func DisplayAllSubscribeTopic() string {
	result := ""
	count := 0
	QC.DevMutex.Lock()
	allsubclient := mqttclient
	QC.DevMutex.Unlock()
	for k, v := range allsubclient.client {
		count++
		result = result + " " + k + ":"
		for kv, vv := range v.subscribeList {
			if vv.inSubscribe {
				result = result + " " + kv
			}
		}
		if len(allsubclient.client) < count {
			result = result + ","
		}
	}
	return result
}

func RunMqttPublish(tcpaddr string, topicname string, messages string) error {
	CreateMqttClient(tcpaddr, QC.Name+":"+tcpaddr)
	time.Sleep(500 * time.Millisecond)
	if mqttclient.client[tcpaddr].clientHandle == nil {
		return errors.New("mqtt client not found")
	}
	QC.DevMutex.Lock()
	receipt := mqttclient.client[tcpaddr].clientHandle.Publish(topicname, 1, false, messages)
	receipt.Wait()
	QC.DevMutex.Unlock()
	return nil
}

func RunMqttSubscribe(tcpaddr string, topicname string) error {
	CreateMqttClient(tcpaddr, QC.Name+":"+tcpaddr)
	time.Sleep(500 * time.Millisecond)
	q.Q(mqttclient)
	if mqttclient.client[tcpaddr].clientHandle == nil {
		return errors.New("mqtt client not found")
	}

	QC.DevMutex.Lock()
	subList, ok := mqttclient.client[tcpaddr].subscribeList[topicname]
	QC.DevMutex.Unlock()
	if ok {
		if subList.inSubscribe {
			return errors.New("topic " + topicname + " in using.")
		}
	}
	//add new subscribe
	newSubscribe := SubscribeStatus{inSubscribe: true}
	QC.DevMutex.Lock()
	mqttclient.client[tcpaddr].subscribeList[topicname] = newSubscribe
	QC.DevMutex.Unlock()
	subtoken := mqttclient.client[tcpaddr].clientHandle.Subscribe(topicname, 0, func(client MQTTClient.Client, msg MQTTClient.Message) {
		q.Q("Received message: ", msg.Payload(), msg.Topic())
		receivemsg := "tcp: " + tcpaddr + " topic: " + msg.Topic() + " message: " + string(msg.Payload()[:])
		syslogerr := SendSyslog(LOG_INFO, "mqttclient", receivemsg)
		if syslogerr != nil {
			q.Q(syslogerr)
		}
	})
	subtoken.Wait()
	return nil
}

func RunMqttUnSubscribe(tcpaddr string, topicname string) error {
	if mqttclient.client[tcpaddr].clientHandle == nil {
		return errors.New("mqtt client not found")
	}
	_, ok := mqttclient.client[tcpaddr].subscribeList[topicname]
	if !ok {
		return errors.New("topic " + topicname + " not found.")
	}
	//unsubscribe
	unSubscribe := SubscribeStatus{inSubscribe: false}
	QC.DevMutex.Lock()
	mqttclient.client[tcpaddr].subscribeList[topicname] = unSubscribe
	QC.DevMutex.Unlock()
	unsubtoken := mqttclient.client[tcpaddr].clientHandle.Unsubscribe(topicname)
	if unsubtoken.Wait() && unsubtoken.Error() != nil {
		return unsubtoken.Error()
	}
	return nil
}

func CreateMqttClient(tcpaddr string, servername string) {

	v, ok := mqttclient.client[tcpaddr]
	if ok {
		// true mean client could use.
		if v.isRun {
			return
		}
	}
	q.Q("create mqtt client", tcpaddr)
	brokertcp := ""
	checkIP := strings.Split(tcpaddr, ":")
	if checkIP[0] != "" {
		brokertcp = "tcp://"
	}
	q.Q(brokertcp, tcpaddr)
	go func() {
		// mqtt client init
		opts := MQTTClient.NewClientOptions().AddBroker(brokertcp + tcpaddr).SetClientID(servername)
		opts.SetKeepAlive(keepAlive * time.Second)
		opts.SetPingTimeout(pingTimeout * time.Second)
		client := MQTTClient.NewClient(opts)
		token := client.Connect()
		if token.Wait() {
			err := token.Error()
			if err != nil {
				q.Q(err)
				return
			}
		}
		defer client.Disconnect(2)

		q.Q("new create mqtt client", tcpaddr)
		newclient := MqttClient{clientHandle: client, tcp: tcpaddr, isRun: true}
		newclient.subscribeList = make(map[string]SubscribeStatus)
		QC.DevMutex.Lock()
		mqttclient.client[tcpaddr] = newclient
		QC.DevMutex.Unlock()

		for {
			isUse := false
			time.Sleep(clientTimeout * time.Second)
			QC.DevMutex.Lock()
			client := mqttclient.client[tcpaddr]
			QC.DevMutex.Unlock()
			for _, v := range client.subscribeList {
				if v.inSubscribe {
					isUse = true
					break
				}
			}
			if !isUse {
				client.isRun = false
				QC.DevMutex.Lock()
				mqttclient.client[tcpaddr] = client
				QC.DevMutex.Unlock()
			}
			q.Q(client.tcp, client.isRun)
			if !client.isRun {
				break
			}
		}
	}()
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

	mqttClientList := make(map[string]MqttClient)
	QC.DevMutex.Lock()
	mqttclient = MqttClientList{client: mqttClientList}
	QC.DevMutex.Unlock()

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
