import Paho from "paho-mqtt";

export function connect(
  host,
  port,
  clientId,
  onConnectionLost,
  onMessageArrived
) {
  const client = new Paho.Client(host, port, clientId);
  // set callback handlers
  client.onConnectionLost = onConnectionLost;
  client.onMessageArrived = onMessageArrived;
  return client;
}

// called when sending a message
export function parsePayload(topic, payload, qos, retained) {
  payload = new Paho.Message(payload);
  payload.destinationName = topic ? topic : "";
  payload.qos = qos ? qos : 1;
  payload.retained = retained ? retained : false;
  return payload;
}
