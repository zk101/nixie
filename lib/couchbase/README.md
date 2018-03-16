# Couchbase Lib

## Usage

### Setup Couchbase Config

    cfg := couchbase.DefaultConfig()

    cfg := couchbase.Config{
        Cluster: "couchbase://localhost",
        Bucket:  "bucket",
        Pass:    "password",
    }

### Create a new Client

    cbclient := couchbase.NewClient(nil)  // This creates a new client with a DefaultConfig

    cbclient := couchbase.NewClient(&cfg)

### Method functions use the same signature as found in gocb library

    cas, err := cbclient.Get(key, valuePtr)

    cas, err := cbclient.Insert(key, value, expiry)

    cas, err := cbclient.Remove(key, cas)

    cas, err := cbpool.Replace(key, value, cas, expiry

    cas, err := cbclient.Touch(key, cas, expiry)

    cas, err := cbclient.Upsert(key, value, expiry)

### Connect manually, the above method functions will attempt to connect to configured CB bucket, and will remain connected after completing operation

    err := cbclient.Connect()

### Close connection, this will set the conn in the client to nil, so that it can be reconnected regardless of err

    err := cbclient.Close()

### Test a connection is working (This will first disconnect an activate connection.  It will leave the connection disconnected)

    if testBool := cbclient.Test(); testBool == false {
        fmt.Println("connection failed")
    }
