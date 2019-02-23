# Go TCP/UDP server and client

## Features

1. Start listening for tcp/udp connections.
2. Be able to accept connections.
3. Read json payload ({"user_id": 1, "friends": [2, 3, 4]})
3. After establishing successful connection - "store" it in memory.
4. When another connection established with the user_id from the list of any other user's "friends" section, notified about it with message {"user_id": <user_id>, "online": true}
5. When the user goes offline, his "friends" (if it has any and any of them online) receives a message {"user_id": <user_id>, "online": false}


## Getting Started

Start a server
```shell
make server
```

Start 3 clients in separate terminals.

```shell
make client1
make client2
make client3
```

To change protocol from TCP to UDP change PROTOCOL = tcp to PROTOCOL = udp inside Makefile.

## Running tests


Run unit tests
```shell
make test
```

Run all tests including integration
```
make test-all
```



## Demo

![demo](https://github.com/anjmao/friends/blob/master/demo.gif)