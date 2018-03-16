package lib

import (
	"strconv"
	"strings"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/zk101/nixie/lib/network"
)

// regConsul sets up a consul service config from Config and returns it
func (c *Core) regConsul(serviceID, addrFilter string) (*consulapi.AgentServiceRegistration, error) {
	address, err := network.FindInetAddress(addrFilter)
	if err != nil {
		return nil, err
	}

	schema := "http://"
	if c.Config.HTTPD.SSL == true {
		schema = "https://"
	}

	serviceTags := strings.Split(c.Config.Controls.ServiceTags, ",")
	serviceTags = append(serviceTags, time.Now().Format("Jan 02 15:04:05.000 MST"))

	checkData := consulapi.AgentServiceCheck{
		HTTP:          schema + address + ":" + strconv.Itoa(c.Config.HTTPD.Port) + "/auth/health",
		TLSSkipVerify: true,
		Interval:      "15s",
		Timeout:       "5s",
		DeregisterCriticalServiceAfter: "60s",
	}

	serviceData := consulapi.AgentServiceRegistration{
		ID:      serviceID,
		Name:    c.Config.Controls.ServiceName,
		Tags:    serviceTags,
		Port:    c.Config.HTTPD.Port,
		Address: address,
		Check:   &checkData,
	}

	return &serviceData, nil
}

// doConsulReg holds the logic for registering the consul service
func (c *Core) doConsulReg() error {
	if c.Config.Controls.ServiceConsul == true {
		service, err := c.regConsul(c.Clients.ServiceID, c.Config.Controls.ServiceAddrFilter)
		if err != nil {
			return err
		}

		if err := c.Clients.Consul.RegisterService(service); err != nil {
			return err
		}
	}

	return nil
}

// doConsulUnreg holds logic to unregister from consul
func (c *Core) doConsulUnreg() error {
	if c.Config.Controls.ServiceConsul == true {
		if err := c.Clients.Consul.UnregisterService(c.Clients.ServiceID); err != nil {
			return err
		}
	}

	return nil
}

// EOF
