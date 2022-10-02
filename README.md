# TinyLB

[![Go Report Card](https://goreportcard.com/badge/github.com/ductnn/tinylb)](https://goreportcard.com/report/github.com/ductnn/tinylb)

**TinyLB** is a simple Load Balancer written in Go. Inspired by [simplelb](https://github.com/kasvith/simplelb)
- Use `RoundRobin` algorithm to send requests into set of backends and support retries too.
- `Healthcheck` to recovery for unhealthy servers in every 1 min.

## Usage

First, create 3 servers test loadbalancing using [docker-compose-test.yml](docker-compose-test.yml):

```sh
➜  tinylb git:(develop) ✗ docker-compose -f docker-compose-test.yml up -d  
Creating network "tinylb_default" with the default driver
Creating tinylb_web2_1 ... done
Creating tinylb_web1_1 ... done
Creating tinylb_web3_1 ... done

# Check
➜  tinylb git:(develop) ✗ docker ps
CONTAINER ID   IMAGE                  COMMAND      CREATED          STATUS         PORTS                                   NAMES
163112a3b548   strm/helloworld-http   "/main.sh"   12 seconds ago   Up 7 seconds   0.0.0.0:3002->80/tcp, :::3002->80/tcp   tinylb_web3_1
d6d1ae94d714   strm/helloworld-http   "/main.sh"   12 seconds ago   Up 6 seconds   0.0.0.0:3001->80/tcp, :::3001->80/tcp   tinylb_web2_1
f717cac76995   strm/helloworld-http   "/main.sh"   12 seconds ago   Up 5 seconds   0.0.0.0:3000->80/tcp, :::3000->80/tcp   tinylb_web1_1
```

To add followings as load balanced backends:
- [http://localhost:3000](http://localhost:3000)
- [http://localhost:3001](http://localhost:3001)
- [http://localhost:3002](http://localhost:3002)

Then, usage:

```sh
➜  tinylb git:(develop) ✗ go run main.go --backends "http://localhost:3000,http://localhost:3001,http://localhost:3002"
2022/10/03 00:14:12 Configured server: http://localhost:3000
2022/10/03 00:14:12 Configured server: http://localhost:3001
2022/10/03 00:14:12 Configured server: http://localhost:3002
2022/10/03 00:14:12 Started listening on localhost:4000
...
# After 1 min
2022/10/03 00:15:12 Health check starting...
2022/10/03 00:15:12 Health check completed
```

Checking on [localhost:4000](http://localhost:4000)

Flag:
```sh
Usage:
  -backends string
        Load balanced backends, use commas to separate
  -port int
        Port to serve (default 4000)
```
