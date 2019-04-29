# nbsnetwork module for nbs connect to each other peer

## 2019.4.11 support MessageSend and StreamSend base on udp
## 2019.4.16 support FileSend and support Resume Send File

## base on udp implements:
### 1. udp connection
### 2. udp reliable send message
### 3. udp reliable send stream
### 4. message rpc
### 5. file transfer
### 6. large file transfer resume

##  Wait to implements:
### 1. peer to peer to connection based on udp
### 2. if p2p failed, use a random agent to create connection
