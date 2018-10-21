# tcp-mitm-proxy
A proxy app that prints out data between a TCP connection.

You may need to tinker a bit to get the app to work for your purposes.

I built this primarily for listening to connections between game client and server. 
It should work for any game based on TCP provided that you route the game servers host name (not ip address) to localhost in the hosts file.

By no means is this a "silver bullet" solution, but it should work for a majority of connections that don't have odd reconnection behaviors. 
