package facts

import (
	"goansible/pkg/executor"
)

type Collector struct {
	conn *executor.SSHConnection
}

func NewCollector(conn *executor.SSHConnection) *Collector {
	return &Collector{conn: conn}
}

func (c *Collector) GatherAll() (map[string]interface{}, error) {
	facts := make(map[string]interface{})

	systemFacts, _ := c.gatherSystem()
	for k, v := range systemFacts {
		facts[k] = v
	}

	networkFacts, _ := c.gatherNetwork()
	for k, v := range networkFacts {
		facts[k] = v
	}

	return facts, nil
}

func (c *Collector) gatherSystem() (map[string]interface{}, error) {
	facts := make(map[string]interface{})

	if out, err := c.conn.Execute("hostname"); err == nil {
		facts["hostname"] = out
	}

	if out, err := c.conn.Execute("uname -s"); err == nil {
		facts["os"] = out
	}

	if out, err := c.conn.Execute("uname -r"); err == nil {
		facts["kernel"] = out
	}

	if out, err := c.conn.Execute("cat /etc/os-release"); err == nil {
		facts["os_release"] = out
	}

	return facts, nil
}

func (c *Collector) gatherNetwork() (map[string]interface{}, error) {
	facts := make(map[string]interface{})

	if out, err := c.conn.Execute("hostname -I"); err == nil {
		facts["ip_addresses"] = out
	}

	if out, err := c.conn.Execute("ip route | grep default"); err == nil {
		facts["default_gateway"] = out
	}

	return facts, nil
}
