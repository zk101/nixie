### Nixie Chat Service

The chat service is not currently complete, and this document captures some of the concepts to guide later devolpement.  In essence, it is designed to not store information at all,  and to also not have state, leaving such concepts to the enduser chat applications.  The idea being, this is a simple chat service for those applications that only require simple chat.  More advanced chat services would likely require a completely different style of message parsing, and run as there own stadalone service.  So, if the concepts laid below are overally simplistic, that is by design.

#### Null (Implemented)

The Null message is a chat logon message.  This is implemented.  It checks the user has the needed ldap structure, and loads that into the users current presence, performing friend lookups as well.  This is a rather heavy and neccasary process.  Logout is handled when the user logs out (or the presence expires).  This is no "online" type states, though, the current setup could support that.

#### Friends

Messages can only be sent to Friends, or groups the user is a member of.  A Friend request can be completed while users are offline, however, such logic needs to be handled by the chat client.  The server will not cache requests for later transmission.  A user searches for another user, and then sends that user a friend request.  The server will update the users cn=friends,ou=chat group with the DN of the user, and forward the friend request message to the other user.  That user either accepts or declines, which sends are reply message.  On Decline, the server removes the users DN from the requestors friends group.  On Accept, the server adds the requests DN to the repliers friends group creating a two way link.  Messages can not be freely exchanged.  A friend delete message can be sent at any time, this will remove the deleted friend from the users friends group.  This will mean the delete user will see the friend as a failed request and can be delete or re-request the friendship.  At this point, the usual friend request logic applies.

#### Search

Allows a user to search for a valid user via email or userid.  Only users that have a valid ldap chat structure should be considered.  Reply should only contain the users Display name and not either the email or the userid.  Wild card searches could be considered, however, if they are Friends, then they should not be required.

#### Groups

This is currently the most problematic concept and may just be removed.  The idea was to allow users to create groups, where the membership of such a group can only be friends, or people that have been requested to be in the group.  It opens up to many side cases such as, can friends of friends be in a group, if so, to becomes possible to send messages to people who may not want this.  Another option is to allow for groups only at a global level, or to allow users to create groups to simply send a message to multple people, however, group creates an illusion that all members can reply to each other.

As such, groups are likely to only be global in nature and automatic in joining.  Such as that required for multiplayer games, or general chat where the option to join can be controlled by an enduser application va configuration.

#### Message

Messages can only be send freely between Friends.  When a message is received from another user, a lookup is performed against their friends group to check they are full friends (cached in presence with a null message), and another check to check the user is online.  Messages can only be send between online users.
