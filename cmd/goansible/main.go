package main

import (
	"flag"
	"fmt"
	"os"

	"goansible/pkg/config"
	"goansible/pkg/executor"
	"goansible/pkg/inventory"
	"goansible/pkg/logger"
	"goansible/pkg/playbook"
)

func main() {
	var (
		playbookFile  = flag.String("playbook", "", "Playbook file to execute")
		inventoryFile = flag.String("inventory", "inventory.ini", "Inventory file")
		verbose       = flag.Bool("v", false, "Verbose output")
		checkMode     = flag.Bool("check", false, "Dry run mode")
		limitHosts    = flag.String("limit", "", "Limit to specific hosts")
		tags          = flag.String("tags", "", "Only run tasks with these tags")
		skipTags      = flag.String("skip-tags", "", "Skip tasks with these tags")
	)

	flag.Parse()

	// Initialize logger
	log := logger.NewConsoleLogger(*verbose)

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config: %v", err)
	}

	// Handle subcommands
	if len(flag.Args()) > 0 {
		switch flag.Args()[0] {
		case "ping":
			runPing(*inventoryFile, log)
			return
		case "facts":
			runFacts(*inventoryFile, log)
			return
		case "version":
			fmt.Println("GoAnsible v1.0.0")
			return
		}
	}

	if *playbookFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Load inventory
	inv, err := inventory.LoadInventory(*inventoryFile)
	if err != nil {
		log.Fatal("Failed to load inventory: %v", err)
	}

	// Load playbook
	pb, err := playbook.LoadPlaybook(*playbookFile)
	if err != nil {
		log.Fatal("Failed to load playbook: %v", err)
	}

	// Create executor
	exec := executor.NewExecutor(cfg, log)
	exec.SetCheckMode(*checkMode)

	// Execute playbook
	if err := exec.ExecutePlaybook(pb, inv); err != nil {
		log.Fatal("Playbook execution failed: %v", err)
	}
}

func runPing(inventoryFile string, log logger.Logger) {
	inv, err := inventory.LoadInventory(inventoryFile)
	if err != nil {
		log.Fatal("Failed to load inventory: %v", err)
	}

	log.Info("Pinging all hosts...")
	for _, host := range inv.AllHosts() {
		if err := executor.Ping(host); err != nil {
			log.Error("%s | UNREACHABLE", host.Name)
		} else {
			log.Success("%s | SUCCESS", host.Name)
		}
	}
}

func runFacts(inventoryFile string, log logger.Logger) {
	inv, err := inventory.LoadInventory(inventoryFile)
	if err != nil {
		log.Fatal("Failed to load inventory: %v", err)
	}

	log.Info("Gathering facts from all hosts...")
	for _, host := range inv.AllHosts() {
		facts, err := executor.GatherFacts(host)
		if err != nil {
			log.Error("%s | FAILED", host.Name)
		} else {
			log.Success("%s | Facts gathered: %d items", host.Name, len(facts))
		}
	}
}
