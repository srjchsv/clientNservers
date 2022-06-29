## Thats an educational golang project- client and servers

We can run one or mulpiple instances of servers, each would write its address and timezone to a txt file. 
Server accepts incoming connections, and writes local time to it.

And we got a client(consumer) that reads txt file with running servers that it would dial. It receives data, and logs it.

### Couple make commands that runs servers and client:

Runs first server:

```
make servergo

```
Runs second server:
```
make servergo2

```
Finally run a consumer with that command:
```
make clientgo
```
