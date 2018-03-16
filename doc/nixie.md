### Nixie Golang Applications

#### Auth (Authentication Service)

This is a RestAPI styled application that provides user services.  Its primary function is; user registeration and deregistration; user login, logout and refreshing of the current session.  Some functions listed are outside of the scope of Auth, and it would be best to call this User instead.  Its primary role to to create session objects which are used by all other services (Auth), with user registration/deregistration as an added function to allow dynamic users for loadtest.

#### Chat (Chat Service)

This feature demonstrates async rpc messaging via a message queue.  It is currently incomplete, and its feature list is fleshed out in the document.  The focus here was to show how messages can be passed in and out simply via a message queue for async style processing.

#### Telemtry (Telemetry Service)

This feature demonstrates async messaging in a single direction, such as that required for telemetry.  A client sending data to be processed that, does not require or should not have, a reply.  Currently telemetry is written to an sql database, however, any production application would use something more suitable to the task of mass telemetry collection.

#### WS (Web Socket Service)

This is the primary socket daemon to allow bi-directional message flow.  WebSockets was chossen as it is a well documented protocol, which is already widely supported by intermiary infrastructure.  It provides a simple read write socket, an is easily setup with 2 http messages.

### Nixie Golang Applications - Tools

#### ConfigLoad (Consul Configuration Loader)

This is a support application to allow for static file configuration to be easily loaded into Consul.
