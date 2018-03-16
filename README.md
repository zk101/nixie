# nixie
A websocket based server platform written in golang.  Nixie is intended for use as the server component to a client server modeled web application.  This project uses a mono repo design which includes all the libraries, vendored packages and build scripts.  

## Components
Nixie consists of two primary applications, auth and ws.  The other applications prove out specific use cases and provide testing.

### Auth

This is a rest like api built using golangs net/http and uses protobuf for data packaging.  It primary purpose is authentication.  Features such as un/register were added to assist dynamic loadtesting.  Couchbase is used for session (presence) storage, with ldap being used for user object storage.  Messages use sha3 hashing.

The authentication model is as follows:

    1.  User sends login protobuf containing their username/password combo.
    2.  An LDAP lookup/bind process is done
    3.  A Session key, sign and cipher data is created, packaged and returned to the user.
    4.  The user uses the session key, sign and cipher data for all future requests in this session

Sessions expire after 5minutes, so tapping refresh is required to extend their session.  This is done automatically in the ws component every 150 secs of successful messages on an active connection.

### WS

This is the websocket component and primary message router.  Once a user has logged in, they can then form a websocket connection to ws.  They must send a null message to the server which sets up the users data.  The null message is special in that its security level is set to none (this prevents the ws server attempting to check the hash) however it still requires the msg hash to be set (which is used to validate the user).  Details of the client side implementation can be found in the loadtest application.

The websocket implementation is not standard by design.  As with all other components, effort was made to harden the services.  As such, the ws server will only accept binary websocket messages, all others (including ping/pong) will cause the connection to be dropped.  A custom ping/pong using signed messages exists.  The theory here is pretty simple, its not a generic websocket server and anything can connect to, its the backend to an applicaiton which is designed to connect to it.  It would be simple enough to switch over to supporting a more standardised websocket implementation.

Once a connection is established, a user can send messages using the wsmsg wrapper protobuf.  This provides capabilities to do plaintext, signed and signed encrypted payloads.  The protobuf is specifically seperate to the payload ([]byte) because it doesn't need to care about what its routing, just a msgtype to tell it what code to run for routing and a seclevel to tell it how to check the message is good.

In addition to server message type; pingpong, server time, latency; there is support for async message via rabbitmq.  Both rpc and single direction messaging is supported.  sync based messaging can be simply added, however, at the time of this writing, code to prove this out is not included.  There is enough code already to make adding this simple.

### Loadtest

The loadtest application was created to prove the client side implemtation and to also place the server components under load to generate metrics.  Grafana graph templates are included.  I created a couple of profiles for testing;  auth, just tests the auth server; all, tests all the components without a advanced chat tests; peer, was an attempt to implement a fully dynamic chat server tests with loadtest users adding each other as friends and sending messages.  It worked some of the time.  I believe its just a timing issue, and is  the last bit I was working on before halting this project.

### Chat, Telemetry

The chat server is a full rpc system using rabbitmq.  It is designed to pass messages between users.  Telemetry is a simple worker to show the processing of messages which don't need a reply (such as telemetry).  While this is backed into sql, a more production capable server would use something else (cloudera, redshift, ...).  As a proof of concept, it worked fine. 

### Configload

I added a tool to simply load config into consul.  All the applications will register into a consul cluster which was used to configure prometheus.  They can also use consul to load there config (or override file config).  This tool was to allow for consul to be loaded simply.  The config library is built using a struct of structs (they must be structs not pointers to structs), and provides a recursive reflect function to both both generate the consul (or environment) key, to do overwrite the config struct value which an appropriately typed key.  Its simple, and works rather well.  

#### Final notes

This project was largely made to learn golang.  The intent was to also create a server backend for future projects.  While I have moved on to working on other things (projects for which this server will form the basis of the backend), I'm placing this out for others to rumage around in and steal what they like (as i have stolen from the interweebs).  Its not really intended to be included in other projects, however, I'm open to moving some of the libraries into there own repos for just such a purpose.

cheers, zk
