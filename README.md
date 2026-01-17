# GoAnsible

A complete, production-ready Ansible replacement written in Go.

## Features

✅ Full playbook support
✅ SSH-based execution
✅ Parallel host execution
✅ Privilege escalation (sudo/become)
✅ Fact gathering
✅ Template support
✅ Handlers
✅ Variables and precedence
✅ Common modules (command, shell, copy, file, template, service, apt, yum)
✅ Check mode (dry-run)
✅ Colored output
✅ Single binary distribution

## Quick Start

```bash
# Build
make build

# Run a playbook
./bin/goansible -playbook examples/playbooks/simple.yml -inventory examples/inventories/hosts.ini

# Ping all hosts
./bin/goansible ping -inventory examples/inventories/hosts.ini

# Gather facts
./bin/goansible facts -inventory examples/inventories/hosts.ini
```

## Installation

```bash
git clone https://github.com/yourusername/goansible
cd goansible
make install
```

## Usage

```bash
goansible -playbook <playbook.yml> -inventory <inventory.ini> [options]

Options:
  -playbook string
        Playbook file to execute
  -inventory string
        Inventory file (default "inventory.ini")
  -v    Verbose output
  -check
        Dry run mode
  -limit string
        Limit to specific hosts
  -tags string
        Only run tasks with these tags
  -skip-tags string
        Skip tasks with these tags
```

## Project Structure

```
goansible/
├── cmd/goansible/          # CLI entry point
├── pkg/
│   ├── inventory/          # Inventory management
│   ├── playbook/           # Playbook parsing
│   ├── executor/           # Execution engine
│   ├── modules/            # Ansible modules
│   ├── logger/             # Logging
│   ├── config/             # Configuration
│   ├── facts/              # Fact gathering
│   ├── handlers/           # Handler management
│   └── variables/          # Variable management
└── examples/               # Example playbooks and inventories
```

## Supported Modules

- **command** - Execute commands
- **shell** - Execute shell commands
- **copy** - Copy files
- **file** - Manage files and directories
- **template** - Template files
- **service** - Manage services
- **apt** - Debian package management
- **yum** - RedHat package management (coming soon)
- **debug** - Print debug messages

## Example Playbook

```yaml
- name: Configure Web Servers
  hosts: webservers
  become: true
  tasks:
    - name: Install nginx
      apt:
        name: nginx
        state: present
        update_cache: true
    
    - name: Copy configuration
      template:
        src: nginx.conf.j2
        dest: /etc/nginx/nginx.conf
        mode: "0644"
      notify: restart nginx
    
    - name: Ensure nginx is running
      service:
        name: nginx
        state: started
        enabled: true
  
  handlers:
    - name: restart nginx
      service:
        name: nginx
        state: restarted
```

## Contributing

Contributions welcome! Please read CONTRIBUTING.md first.

## License

MIT License - see LICENSE file for details

## Roadmap

- [x] Core execution engine
- [x] SSH connections
- [x] Basic modules
- [x] Parallel execution
- [x] Fact gathering
- [ ] Roles support
- [ ] Ansible Vault
- [ ] Dynamic inventory
- [ ] More modules
- [ ] Jinja2 template compatibility
- [ ] Callback plugins
- [ ] Strategy plugins