# How to use

Run this by using `go build *.go` and run the built executable (called `Ivy` by default).

**NOTE**: Keys can only be `int`s and all manager and client IDs are 0-indexed. The first manager (`manager[0]`) is the master server by default.

The program will open a REPL interface. The available commands are

1. Insert new KV pair - `[client ID] write [key] [value]`. For example if Client 1 wants to write the pair `{101: "world"}`, the command will be `1 write 101 World`.
2. Get KV pair - `[client ID] read [key]`. If Client 3 wants to read the data from key `101`, the command will be `3 read 101`.
3. Take down server (_Part 2 only_) - `[master ID] down`. Note the maximum number of servers is capped according to the initialisation loop in `func main`.
4. Bring up server (_Part 2 only_) - `[master ID] up`.

Example Usage

```bash
0 write 1 hello # write new value
0 down          # take down current master
1 read 1        # read value from before server was down
1 write 1 new   # overwrite old value
0 up            # bring back old server (will not become master)
2 write 2 new   # write new value to new key
```

## Experiments

1. Performing 10 simultaneous reads and writes, we can see that the fault tolerant version is much slower than the one without. In general, reading is about 1.5x slower (varies highly depending on interleaving), and writing is 3x slower.
2. Both scenarios are accounted for, and the system can continue functioning under both circumstances. (b) is faster since election is not required and the new replica is added to the available replicas. However, (a) requires an election since the master server was taken down.
3. The application can continue functioning as long as there is atleast 1 server up. If the last server is taken down, the application panics and shuts down.

Due to contstraints of the REPL interface, no reads/writes can be performed until the election is complete, so there is no consistency errors possible. However in the current implementation, if a server dies while it is still collecting election ACKs, it is possible for the system to become hung as the election was not complete and everyone is waiting for the server to acknowledge its master status.
