package variables

type VarManager struct {
	vars map[string]interface{}
}

func NewVarManager() *VarManager {
	return &VarManager{
		vars: make(map[string]interface{}),
	}
}

func (vm *VarManager) Set(key string, value interface{}) {
	vm.vars[key] = value
}

func (vm *VarManager) Get(key string) (interface{}, bool) {
	val, ok := vm.vars[key]
	return val, ok
}

func (vm *VarManager) Merge(other map[string]interface{}) {
	for k, v := range other {
		vm.vars[k] = v
	}
}

func (vm *VarManager) All() map[string]interface{} {
	return vm.vars
}
