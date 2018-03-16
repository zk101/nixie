package userchat

import (
	"log"
	"os"
	"testing"

	"github.com/zk101/nixie/lib/config"
	"github.com/zk101/nixie/lib/ldap"
	"github.com/zk101/nixie/lib/testutil"
	"github.com/zk101/nixie/models/ldap/auth/usernew"
)

// Global vars used for testing
var (
	conf     *testutil.Config
	connLDAP *ldap.Client
	user     *usernew.Model
)

// defaultUser returns a pre-populated Model
func defaultUser() *usernew.Model {
	user := usernew.New()

	user.DN = "uid=testuser,ou=users," + conf.LDAP.Base
	user.Base = conf.LDAP.Base
	user.Name = "Test User"
	user.Username = "testuser"
	user.Password = "password"

	return user
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

	user = defaultUser()
	if err := user.Create(connLDAP); err != nil {
		log.Fatal(err.Error())
	}

	retCode := m.Run()

	if err := user.Remove(connLDAP); err != nil {
		log.Fatal(err.Error())
	}

	os.Exit(retCode)
}

// TestChat
func TestChat(t *testing.T) {
	userchat := New()
	userchat.DN = user.DN

	//if err := userchat.Create(connLDAP); err != nil {
	//	t.Error(err.Error())
	//}

	if err := userchat.Fetch(connLDAP); err != nil {
		t.Error(err.Error())
	}

	//for key, value := range userchat.Friends {
	//	fmt.Println(key)
	//	fmt.Printf("%t\n", value)
	//}

	//t.Error("Force Error")
}

// EOF
