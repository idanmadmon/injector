# injector
sender - receiver cli injector

### help
```
Usage:
  injector [flags]
  injector [command]

Available Commands:
  completion        Generate the autocompletion script for the specified shell
  help              Help about any command
  receiver          create a http receiver
  sender            create a http sender
  stream-duplicator create a stream-duplicator http receiver sender

Flags:
  -h, --help      help for injector
  -v, --version   version for injector
```

# Usages

## Normal Usage
### Sender
```
.\injector.exe sender --url http://127.0.0.1:1238/ --amount 10 --size 10MB --timeout 10s
```

### Receiver
```
.\injector.exe receiver --address 0.0.0.0:1238
```

## Rate Limit Usage
### Sender
```
.\injector.exe sender --url http://127.0.0.1:1238/ --amount 10 --size 10MB --timeout 10s
```

### Receiver
```
.\injector.exe receiver --address 0.0.0.0:1238
```

## Stream Duplicator Usage
### Sender
```
.\injector.exe sender --url http://127.0.0.1:1238/ --amount 10 --size 10MB --timeout 10s
```

### Stream Duplicator
```
.\injector.exe stream-duplicator --address 0.0.0.0:1238 --dest-url http://127.0.0.1:1236/  --dest-url http://127.0.0.1:1237/
```

### Receiver 1
```
.\injector.exe receiver --address 0.0.0.0:1236
```

### Receiver 2
```
.\injector.exe receiver --address 0.0.0.0:1237
```

# Lab Setup
### Environment Setup

Setup: <br>
sender rhel/debian virutal machine<br>
receiver rhel/debian virutal machine

You can download VMware workstation from here for free:<br>
[VMWare Workstation Pro](https://support.broadcom.com/group/ecx/productfiles?subFamily=VMware%20Workstation%20Pro&displayGroup=VMware%20Workstation%20Pro%2017.0%20for%20Windows&release=17.6.2&os=&servicePk=526672&language=EN)

And you can download Red Hat Enterprise Linux 9 64-Bit from here for free:<br>
[Red Hat Enterprise Linux 9.5 64-Bit](https://developers.redhat.com/products/rhel/download#publicandprivatecloudreadyrhelimages)

Install the Virtual Machine and setup 2 machine, sender and receiver

### Utils

Run update-injector to build injector for linux and put it in the lab (sender and receiver)

# Lab Usage
#### Run Injector in the sender and receiver

### Monitor
Sender's Cockpit's monitor:<br>
[https://192.168.108.128:9090/network](https://192.168.108.128:9090/network)

Receiver's Cockpit's monitor:<br>
[https://192.168.108.129:9090/network](https://192.168.108.129:9090/network)