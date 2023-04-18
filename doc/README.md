# Introduction to mnms software

This contains basic information about `mnms` software which is a network management software inteneded to be used with Atop produced switches and other devices. 

The `mnms` software is intended to be secure, scalable and extensible.

## Installation

Make a directory and change into that directory.

```
mkdir mydir
cd mydir
```

Copy the mnms zip file provided into this directory and unzip the content there.

If you are on a Windows machine and you don't have PCAP installed, please
run the npcap installer that is provided in the zip file.

You can run on a single machine but it is recommended that you install the provided zip file on each machine you intend to run mnms.

The main binary executable for the backend is called `mnmsctl` and the frontend is called `frontend`.

## Version

Run with -version flag to see the details of the software release version.

## mnms distributed model

The `mnms` services run in a distributed environment.  This allows handling many devices that are scattered in many networks.  For example, in a large deployment, there may be thousands of Atop produced devices such as switches and serial devices.  Some of the devices may be connected to a LAN while some of the other devices many be connected to another LAN.

Users can run a client node service on a machine directly attached to a given LAN.  If there are distinct LAN-1 and LAN-2 networks, and different sets of devices are connected to LAN-1 and LAN-2, you may run one client node service on a machine that is connected to LAN-1. And run another client node service on a machine on LAN-2.  

The client node service will scan and configure the devices on the directly attached LAN.  They also report details about the devices to the root service and communicate configuration details as directed by the root.

The root service runs on a machine that is reachable by client node services.  Usually the root is expected to run on a machine that is hosted in the cloud or a server that is accessible from all client node machines.

The root service aggregates information from client node services.  It also directs and distributes work to the client node services.  The root provides API to the web UI via REST API.  It also provides support for command line tools (CLI tools) which are mapped to APIs.  

Many client service nodes can be supported per root.  This architecture can scale up to support thousands of devices.


## How to run mnms


Typically, users will login to a machine that is capable running a service software such as a Linux PC or a Windows PC.  Users can unzip the installation package, and run the software in the directory where software is unzipped.

This can be done by using a command line shell such as bash or Powershell.


### Run root service

The basic command to run a root service instance named `root` is:

```
mnmsctl -n root -R
```

The -n allows users to specify any name to be given to the instance of this mnms service which will run in root mode as specified via -R flag.

You may need to specify where mnmsctl program is.  On Linux, if you unzipped the release into a directory /abc/def, then

```
$ cd /abc/def
$ ./mnmsctl -n root -R
```

On windows
```
$ cd C:/abc/def
$ ./mnmsctl.exe -n root -R
```

Note that mnmsctl services should run in super user mode.  On linux,

```
$ sudo ./mnmsctl -n -root -R
```

One windows, create Administrator mode powershell or use `gsudo.exe`.  Gsudo can be downloaded from Chocolatey.

```
$ choco install gsudo
```

In production clusters mnmsctl will need to run with -M flag. The -M flag must be the first flag. It indicates to mnmsctl to restart itself when there is a critical error.  In a clustered environments where there are many machines running many instances of distributed mnmsctl programs in Root and client modes, it is advisable to run with -M mode enabled so as to continue to have high up times even when faced with occasional failures.

The -M is "monitor mode" and it will also log any stack traces for post mortem analysis between runs.

The mnms system allows for various levels of logs to be turned on and off during runtime for debugging and analysis. This is critical for managing complex cluster deployments.  There are API commands and flags to enable log output files and patterns.

### Run a client node service on another machine

```
mnmsctl -n client1 -s -r http://10.10.10.1:27182 -rs 10.10.10.1:5514
```

if root service is running at 10.10.10.1

### Run another client node service on another machine
```
mnmsctl -n client2 -s -r http://10.10.10.1:27182 -rs 10.10.10.1:5514
```

### Run a web UI

```
frontend
```

Use a web browser and connect to localhost:9000 using username admin and password default. Change the password as soon as possible to a more secure one. By default UI will connect to the backend at http://localhost:27182 which can be changed in the UI configuration menu.


## Getting help

```
mnmsctl help
```

Will produce information about various API commands and features available for the users.

```
mnmsctl -h
```

Will list flag options that can be specified when running the `mnmsctl` program.

## APIs and SDK

Both the CLI and REST APIs share consistent set of methods to perform different configuration and maintenance features for Atop devices.

There is a single tool `mnmsctl` which can do most of the work. It can be used to run root and client node services.  It can also be used to run CLI commands.  It can even be used to create RSA public and private key pairs and encrypt and decrypt data to be used with `mnms`.

Use of APIs allow users to script set a set of actions that can be performed on a group of devices. Instead of configuring one device at a time, it is possible to use scripts to affect changes to many machines in parallel.

The latest configuration actions per device are recorded in the history which can be viewed.

The `mnms` software is designed so that it is possible to customize and extend the API and features quickly for different use cases.  Because mnms is implemented as a Go language package, it is possible to create custom versions of code that uses mnms package as SDK to implement custom actions and behaviors.



## Running TLS

The `mnms.zip` contains caddy server. We use caddy to provide TLS termination. The certificate is automatically obtained from Let's Encrypt and caddy reverse-proxy to the root mnms service provides https.

For example,

```
sudo caddy reverse-proxy --from 10.10.10.1.sslip.io --to :27182
```

## Cluster high availability

Cluster of nodes that run mnms Root and client services can be made reilient to failures by using client service monitor mode (-M flag) and caddy load balancing for Root services.

Services can be deployed in Google cloud, Azure, AWS and other clouds as well as inside docker or Kubernetes clusters.

## SNMP MIB Browser

The Web UI frontend includes a MIB browser feature which can be used to manage SNMP compatible devices.


## Syslog aggregation

When devices generate syslog messages, they can be forwarded to a syslog forwarder which runs inside a mnms client node service.  The syslog messages will be forwarded to the remote syslog service.  The client node services typically will forward to a root syslog service as specified via -rs flag.  The root can further forward syslog to the ultimate destination such as a rsyslog aggregation service or other commerical syslog aggregation services.  The root service can also sink the syslog if the remote  syslog service is not specified via -rs.  In this mode, mnms root service will act as a syslog aggregator and save the incoming syslogs in the local disk, roll and compress the logs as configured.

## Alerts and events

Alerts and event messages are forwarded to UI via websocket.  They are also recorded in syslog for aggregation and analytics.

## MQTT message service

Basic support for mqtt publish and subscribe messaging.

## OPC UA 

Basic support for OPC UA client features.

## Experimental secure custom agent

Code exists for implementation of a secure custom agent that can be embedded in target devices.  This code requires device specific customizations as per customer use cases.

## Additional documentation

Additional text and markdown files may be included in this package. Please refer to the additional documentation for different features of the software. 
