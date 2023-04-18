# MNMS mqtt guide

`mnms` provide basic support for mqtt publish and subscribe messaging. Both Mqtt broker and client run on the client node. When a message comes in, the message will be sent to the server node by syslog. And `mnms` does not have its own subscription item, we only provide basic subscription and publishing.

## Mqtt broker ##

`mnms` contains an mqtt broker that can accept subscriptions and publish messages. The default port is :11883. If you want to change port, you can add commmand `mb`:
```
mnmsctl -s -n client -mb ":1883" -r http://192.168.12.1:27182
```
And when mqtt broker receive the publish messages, it would send the publish messages to server node by syslog.

## Mqtt client ##

`mnms` supports subscribing and publishing messages to other remote mqtt brokers. If you need to subscribe messages to the remote mqtt broker, you can write the command:

In command line,
```
mnmsctl mqtt sub [tcp address] [topic] 
Example:
	mnmsctl mqtt sub 192.168.12.1:1883 topictest
```
In UI script,
```
mqtt sub [tcp address] [topic] 
Example:
	mqtt sub 192.168.12.1:1883 topictest
```

If you need to publish messages to the remote mqtt broker, you can write the command:

In command line,
```
mnmsctl mqtt pub [tcp address] [topic] [message]
Example:
	mnmsctl mqtt pub 192.168.12.1:1883 topictest "this is message."
```
In UI script,
```
mqtt pub [tcp address] [topic] [message]
Example:
	mqtt pub 192.168.12.1:1883 topictest "this is message."
```

If you don't need to subscribe messages to the remote mqtt broker, you can write the command:

In command line,
```
mnmsctl mqtt unsub [tcp address] [topic]
Example:
	mnmsctl mqtt unsub 192.168.12.1:1883 topictest
```
In UI script,
```
mqtt unsub [tcp address] [topic] 
Example:
	mqtt unsub 192.168.12.1:1883 topictest
```

And, when the remote mqtt broker does not exist, it will not work to subscribe and publish messages. User must carefully check remote mqtt broker status.

## Not yet implemented mqtt feature ##

1. User can get more information about the machine by subscribing to the `network`, `config` and `device` topic of `mnms`.

2. User can use the UI to observe all topics, and easily add/remove subscribed topics without commands.

3. The message published/subscribed by the user will be displayed on the UI, and the user can easily get the message.