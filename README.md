## rtc-terminal (CLI)

1. Install rtc-ssh on remote device/server from the repository: https://github.com/mxseba/rtc-ssh

Note the uuid key generated in the rtc-ssh application.

2. Get source rtc-terminal on local computer using the Go compilator:
```
go get -u github.com/mxseba/rtc-terminal
cd $GOPATH/bin
rtc-terminal -uuid=<UUID key remote device>
```

