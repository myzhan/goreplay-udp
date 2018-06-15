[![Go Report Card](https://goreportcard.com/badge/github.com/myzhan/goreplay-udp)](https://goreportcard.com/report/github.com/myzhan/goreplay-udp)

# About

GoReplay-udp is copycat of [goreplay](https://github.com/buger/goreplay), works on UDP tracffic.

It's currently a toy project and not tested as well as goreplay.

All credit goes to Leonid Bugaev, [@buger](https://twitter.com/buger), https://leonsbox.com

# Build

```bash
sudo apt-get install libpcap libpcap-dev flex bison -y
git clone github.com/myzhan/goreplay-udp $GOPATH/src/github.com/myzhan/goreplay-udp
cd $GOPATH/src/github.com/myzhan/goreplay-udp
go build -ldflags '-extldflags "-static"'
```

# Usage

```
# Running as non root user
sudo etcap "cap_net_raw,cap_net_admin+eip" ./goreplay-udp
# Test
sudo ./goreplay-udp --input-udp :22 --output-stdout
# Capture
sudo ./goreplay-udp --input-udp :22 --output-file dns.req
# Replay Online
sudo ./goreplay-udp --input-udp :22 --output-dns localhost:2222
# Replay Offline
sudo ./goreplay-udp --input-file dns.req --output-udp localhost:2222
```