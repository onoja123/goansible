package playbook

type Task struct {
	Name   string                 `yaml:"name"`
	Module string                 `yaml:"module,omitempty"`
	Args   map[string]interface{} `yaml:"args,omitempty"`

	// Module shortcuts
	Command  string        `yaml:"command,omitempty"`
	Shell    string        `yaml:"shell,omitempty"`
	Copy     *CopyArgs     `yaml:"copy,omitempty"`
	File     *FileArgs     `yaml:"file,omitempty"`
	Template *TemplateArgs `yaml:"template,omitempty"`
	Service  *ServiceArgs  `yaml:"service,omitempty"`
	Apt      *PackageArgs  `yaml:"apt,omitempty"`
	Yum      *PackageArgs  `yaml:"yum,omitempty"`
	Git      *GitArgs      `yaml:"git,omitempty"`
	User     *UserArgs     `yaml:"user,omitempty"`
	Debug    *DebugArgs    `yaml:"debug,omitempty"`

	// Control flow
	When      string   `yaml:"when,omitempty"`
	Loop      []string `yaml:"loop,omitempty"`
	WithItems []string `yaml:"with_items,omitempty"`
	Register  string   `yaml:"register,omitempty"`
	Notify    []string `yaml:"notify,omitempty"`
	Tags      []string `yaml:"tags,omitempty"`
	Changed   bool     `yaml:"changed_when,omitempty"`
	Failed    string   `yaml:"failed_when,omitempty"`
	Ignore    bool     `yaml:"ignore_errors,omitempty"`
}

type CopyArgs struct {
	Src    string `yaml:"src"`
	Dest   string `yaml:"dest"`
	Mode   string `yaml:"mode,omitempty"`
	Owner  string `yaml:"owner,omitempty"`
	Group  string `yaml:"group,omitempty"`
	Backup bool   `yaml:"backup,omitempty"`
}

type FileArgs struct {
	Path    string `yaml:"path"`
	State   string `yaml:"state"`
	Mode    string `yaml:"mode,omitempty"`
	Owner   string `yaml:"owner,omitempty"`
	Group   string `yaml:"group,omitempty"`
	Recurse bool   `yaml:"recurse,omitempty"`
}

type TemplateArgs struct {
	Src  string `yaml:"src"`
	Dest string `yaml:"dest"`
	Mode string `yaml:"mode,omitempty"`
}

type ServiceArgs struct {
	Name    string `yaml:"name"`
	State   string `yaml:"state"`
	Enabled *bool  `yaml:"enabled,omitempty"`
}

type PackageArgs struct {
	Name   string `yaml:"name"`
	State  string `yaml:"state"`
	Update bool   `yaml:"update_cache,omitempty"`
}

type GitArgs struct {
	Repo    string `yaml:"repo"`
	Dest    string `yaml:"dest"`
	Version string `yaml:"version,omitempty"`
}

type UserArgs struct {
	Name   string `yaml:"name"`
	State  string `yaml:"state,omitempty"`
	Groups string `yaml:"groups,omitempty"`
	Shell  string `yaml:"shell,omitempty"`
}

type DebugArgs struct {
	Msg string `yaml:"msg,omitempty"`
	Var string `yaml:"var,omitempty"`
}
