package inventory

type Host struct {
	Name     string
	Address  string
	Port     int
	User     string
	KeyFile  string
	Password string
	Vars     map[string]interface{}
	Groups   []string
}

func NewHost(name string) *Host {
	return &Host{
		Name:    name,
		Address: name,
		Port:    22,
		User:    "root",
		Vars:    make(map[string]interface{}),
		Groups:  []string{},
	}
}

func (h *Host) GetVar(key string) (interface{}, bool) {
	val, ok := h.Vars[key]
	return val, ok
}

func (h *Host) SetVar(key string, value interface{}) {
	h.Vars[key] = value
}
