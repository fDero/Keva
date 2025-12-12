# Keva: Distributed Key-Value Storage

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Raft](https://img.shields.io/badge/raft-consensus-darkgreen?style=for-the-badge)
![License](https://img.shields.io/badge/license-Apache%202.0-purple.svg?style=for-the-badge)

**Keva** is a lightweight, distributed key-value storage system written in Go, powered by the **Raft consensus algorithm** for fault tolerance and strong consistency. Every node in the Keva cluster exposes a minimal HTTP-based REST api. The master (leader) node can perform true write operations, while 
every other node (slaves/followers) can perform read operations and internally forward write-requests to the current master. Master nodes are 
randomly elected via democratic election performed by all the alive peers.

---

## Features

- **Strong Consistency**: All reads and writes are linearizable thanks to Raft
- **Fault Tolerance**: As long as at least one node is up the cluster survives
- **Simple API**: Easy-to-use key-value interface

---

## Getting started

The easiest way to test Keva is spinning up a cluster composed of three instances all running on the local machine.
To achive such setup, creating three separate working directories (e.g. `workdir1`, `workdir2`, `workdir3`) 
and a `keva.toml` config file with your cluster settings. Here there's an example of such file:

```toml
[[node]]
identity='host1'
address='localhost'
keva_port='370'
user_port='380'
wait_time='5'

[[node]]
identity='host2'
address='localhost'
keva_port='371'
user_port='381'
wait_time='4'

[[node]]
identity='host3'
address='localhost'
keva_port='372'
user_port='382'
wait_time='3'
```

Then you will need to spawn at least two out of the three instances (one more then half of the nodes) to 
get the cluster up and running. You can launch an instance like this:

```bash
keva --config-file './keva.toml' --node-identity 'host1' --working-directory './workdir1'
```

## API
Every node in the cluster exposes a REST api, that is internally implemented in such a way that 
reads are performed locally on the recieving node, while writes are forwarded to the current master
node (leader). The API is very user friendly, and all you need to do is send requests to 
the following route with `POST`, `GET` or `DELETE`. When trying to query for a non present key, you will get a `204` response (HTTP/No Content)

```endpoint
http://<NODE-IP>:<USER-PORT>/v1/storage/key/<KEY>
```
