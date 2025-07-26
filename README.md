# Termtalk
A chat application running in your terminal

## Why Termtalk
Termtalk can run as a self-hosted server and a client inside of your terminal. The client provides two interfaces: cli and tui. 
>[!warning]
> The server uses ephemeral storage. Message history is lost when:
> - It is overwritten (buffer limit: 50 messages)
> - The server shuts down
 
> History size is currently unconfigurable.


### Requirements
- Go version 1.24.4

## Installation and usage
```bash
go install github.com/sternerr/termtalk@latest
```
### Run as server
```bash
termtalk server --host <ip> --port <port>
```
### Run as client
```bash
termtalk client --mode tui --host <ip> --port <port>
termtalk client --mode cli --host <ip> --port <port>
```
