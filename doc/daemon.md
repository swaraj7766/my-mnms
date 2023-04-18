# Daemon

 run a program as a service (daemon)

### required

run as root or administrator

## Command

run | start | stop | restart | install |  uninstall | status

-n: as service name

### run

just run as usual

```sh
./mnmsctl/mnmsctl.exe -R -n root -s -P ".*" -daemon run
```

### install

install will run automatically as service

```sh
./mnmsctl/mnmsctl.exe -R -n root -s -P ".*" -daemon install
```

### uninstall

uninstall and stop service

```sh
./mnmsctl/mnmsctl.exe -n root -daemon uninstall 
```

### start

if installed, start service

```sh
./mnmsctl/mnmsctl.exe -n root -daemon start 
```

### stop

if service is running, stop it

```sh
 ./mnmsctl/mnmsctl.exe -n root -daemon stop 
```

### restart

restart service

```sh
./mnmsctl/mnmsctl.exe -n root -daemon restart 
```

### status

show service status

```sh
./mnmsctl/mnmsctl.exe -n root -daemon status 
```

