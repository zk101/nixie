package usernew

import (
	"log"
	"os"
	"testing"

	"github.com/zk101/nixie/lib/config"
	"github.com/zk101/nixie/lib/ldap"
	"github.com/zk101/nixie/lib/testutil"
)

// Global vars used for testing
var (
	conf     *testutil.Config
	connLDAP *ldap.Client
)

// defaultUser returns a pre-populated Model
func defaultUser() *Model {
	user := New()

	user.DN = "uid=testuser,ou=users," + conf.LDAP.Base
	user.Base = conf.LDAP.Base
	user.Name = "Test User"
	user.Username = "testuser"
	user.Password = "password"

	return user
}

// attributes test data is a series of arrays broken into to sections, data that should pass and data that should not.
var attrNamePass = []string{
	"Jane Doe",
	" Jack Doe _",
	"Jane Doe 123 _",
	"-Jack 69 Doe_",
	"Hello 世界",
}

var attrNameFail = []string{
	"Jane\nDoe",
	"Jack$Doe",
	"Jack<?Doe?>",
	"",
}

// TestMain sets up LDAP and tears it down again
func TestMain(m *testing.M) {
	var err error

	conf, err = testutil.LoadConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	cacert, err := config.LoadCAcerts(conf.Controls.CAcertPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	connLDAP = ldap.NewClient(&conf.LDAP, cacert)
	defer connLDAP.Close()

	if err := connLDAP.Test(); err != nil {
		log.Fatal(err.Error())
	}

	os.Exit(m.Run())
}

// TestCreate
func TestCreate(t *testing.T) {
	user := defaultUser()

	if err := user.Create(connLDAP); err != nil {
		t.Fatalf("TestCreate Create failed: %s\n", err.Error())
	}

	if err := connLDAP.CheckBind(user.DN, user.Password); err != nil {
		t.Errorf("TestCreate CheckBind failed: %s\n", err.Error())
	}

	if err := user.Remove(connLDAP); err != nil {
		t.Fatalf("TestCreate remove failed: %s\n", err.Error())
	}
}

//TestAttrName
func TestAttrName(t *testing.T) {
	user := defaultUser()

	for _, value := range attrNamePass {
		user.Name = value
		_, err := user.userAttributes()
		if err != nil {
			t.Errorf("TestAttrName (%s) expected error to be nil, got: %s\n", value, err.Error())
		}
	}

	for _, value := range attrNameFail {
		user.Name = value
		_, err := user.userAttributes()
		if err == nil {
			t.Error("TestAttrName expected error, got nil")
		}
	}
}

// EOF
