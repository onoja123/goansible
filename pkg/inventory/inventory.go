package inventory

type Inventory struct {
	Hosts     map[string]*Host
	Groups    map[string]*Group
	GroupVars map[string]map[string]interface{}
	HostVars  map[string]map[string]interface{}
}

type Group struct {
	Name     string
	Hosts    []*Host
	Children []*Group
	Vars     map[string]interface{}
}

func NewInventory() *Inventory {
	return &Inventory{
		Hosts:     make(map[string]*Host),
		Groups:    make(map[string]*Group),
		GroupVars: make(map[string]map[string]interface{}),
		HostVars:  make(map[string]map[string]interface{}),
	}
}

func (inv *Inventory) AddHost(host *Host) {
	inv.Hosts[host.Name] = host
}

func (inv *Inventory) AddGroup(group *Group) {
	inv.Groups[group.Name] = group
}

func (inv *Inventory) GetHosts(pattern string) []*Host {
	if pattern == "all" {
		return inv.AllHosts()
	}

	if group, ok := inv.Groups[pattern]; ok {
		return group.Hosts
	}

	if host, ok := inv.Hosts[pattern]; ok {
		return []*Host{host}
	}

	return []*Host{}
}

func (inv *Inventory) AllHosts() []*Host {
	hosts := make([]*Host, 0, len(inv.Hosts))
	for _, h := range inv.Hosts {
		hosts = append(hosts, h)
	}
	return hosts
}
