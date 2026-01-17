package inventory

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func LoadInventory(filename string) (*Inventory, error) {
	inv := NewInventory()

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentGroup *Group

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			groupName := strings.Trim(line, "[]")

			if strings.Contains(groupName, ":") {
				parts := strings.Split(groupName, ":")
				if parts[1] == "vars" {
					currentGroup = nil
					continue
				}
			}

			currentGroup = &Group{
				Name:  groupName,
				Hosts: []*Host{},
				Vars:  make(map[string]interface{}),
			}
			inv.AddGroup(currentGroup)
			continue
		}

		if currentGroup != nil {
			host := parseHostLine(line)
			if host != nil {
				host.Groups = append(host.Groups, currentGroup.Name)
				currentGroup.Hosts = append(currentGroup.Hosts, host)
				inv.AddHost(host)
			}
		}
	}

	return inv, scanner.Err()
}

func parseHostLine(line string) *Host {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return nil
	}

	host := NewHost(parts[0])

	for _, part := range parts[1:] {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}

		key, value := kv[0], kv[1]
		switch key {
		case "ansible_host":
			host.Address = value
		case "ansible_port":
			port, _ := strconv.Atoi(value)
			host.Port = port
		case "ansible_user":
			host.User = value
		case "ansible_ssh_private_key_file":
			host.KeyFile = value
		case "ansible_password":
			host.Password = value
		default:
			host.SetVar(key, value)
		}
	}

	return host
}
