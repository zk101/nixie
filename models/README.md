### Models

Models is a complex topic, so this will largely break down the methods used in Nixie.

There is no specific Interface for the various types of models, however, each model is built around the interfaces listed below incase at some point there becomes reason to implement an interface.  They also hold to specific file naming in order to make navigation easier, as the number of Models will increase.

The typical common layout for each model is tho provide a Model struct, that holds all the data and any annotations that describe the data which the model holds.  There may be several "models" that describe the same object (such as the User model), to allow for specific types of interactions.  In the User case, password needs specific handling while most other attributes can be interacted with in a more generalised approach.  There may also be times that a single Model abstractions more complex interactions with two or more other models in differant storage systems.  Hence, the interfaces listed below are guidelines rather than programatic structures.

As a generalised overview, a Model folder holds the following files and content:

    model.go   # This file holds the Model struct, and a New function to return an empty Model ptr
    create.go  # This file holds the Create function, which accepts the required storage interfaces and returns an error
    fetch.go   # This file holds the Fetch function, which accepts the required storage interfaces, a key string and returns an error
    remove.go  # This file holds the Remove function, which accepts the required storage interfaces and returns an error
    update.go  # This file holds the Update function, which accepts the required storage interfaces and returns an error
    util.go    # This file holds any utility functions specific to the Model

For Remove functions, it could be argued that a key string is also required, so the logic of get can be abstracted into it, due to there being no formal interface to adher to, the choice to use this method is left open.  The view point at this stage is that a Fetch would be run before calling Remove, so the models key data is already available within the model.

As a generalised overview, the Model Interface follows the following structure:

    type Model interface {
      Create(storage.Interface) error
      Fetch(storage.Interface, error) error
      Remove(storage.Interface) error
      Update(storage.Interface) error
    }

Where storage.Interface can be found the the storage lib.  As mentioned above, this is a pseudo interface, as some models may require multiple storage interfaces, or additional parameters.

#### Protobuffers

As an additional note, which naturally extents Models, is the use of Protobuffers.  This is used primarily as a transport mechanism as well as a storage structure.  While protobuffers generated code does not provide direct pack and unpack functions, as an extension on the above Model interface, and as a pseudo interface of its own, consider the following interface and as standalone functions:

    type Model interface {
      ...
      Pack() ([]byte, error)
      Unpack([]byte) error
    }

    Pack(*protobuf.Structure) ([]byte, error)
    Unpack([]byte) (*protobuf.Structure, error)

In some instances, it may not even be required to pass a *protobuf.Structure to Pack, as the message content cab be abstracted into the Pack function, as is the case with srvtime.  They may also require many more parameters, however, in general, this is the approach used to structure these actions.

The files used to hold these functions are as follows:

    pack.go    # This file holds the Pack function
    unpack.go  # This file holds the Unpack function

#### Models folder structure

The models folder is structured to allow natrual seperation and logically location of data and transport models.  The first level is based on the data storage system that holds the object, or, in the case of transport, the serialisation method used.  This has an obvious downfall where two or more storage systems are required, in which case, the primary reference storage should be used, or a new folder created that combines them.  i.e. A Model that uses data from both ldap and sql combined, which is then stored in couchbase should use couchbase as its folder.  As this is where the Model is primarily referenced from.  If the same model was not stored, a new folder (ldapsql) could be created for this model.

The next level is named for the service that controls the primary write for the Model.  This is a micro-service application, and while its architecture does not follow the more draconian views of micro service; namely that data should be only read and write from the micro-service that owns the data and all access to that data is via that micro-service, which imposes unneeded network overhead and complexity; we still hold that there is a microservice primarily responsible for the initial creation Model.  That service name is used here.

There third level is the name of the model itself.  Generally, this is the final level for models that have only a single set of interactions.  If a model requires more than one, such as the user model, then this level is the overall name of the model (user) and the next level is the more specific name for the models function; i.e. create, passwd, update.  Where create is the model for creating a new user, passwd is the model to update the users passwd, and update is the model to update all other attributes.  This may be simplify once the model is fleshed out and better logic can be used, however, it does service as a useful example of handling more complex models.  Another example could be related to the mixed storage model mentioned above, where having a failed storage update could make rolling back near impossible, updates to the overall model could be made storage specific.  Often this times of models are really just a useful combination to make retrieval easier and so have simplier models already.

    ./models/ldap/auth/usernew/
    ./models/ldap/auth/userpasswd/
    ./models/ldap/auth/userupdate/
    ./models/pb/auth/login/
    ./models/pb/auth/register/
    ./models/pb/ws/pingpong/

The third level has simplified, however I'm leaving the above thought in place as it may still be relavant for a future problem.  The simplification comes from golang package naming, where it seems saner to include the base model name in the seperate submodels, as the above (modified) example shows.
