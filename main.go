package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
)

// ANSI Color Codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
)

type Profile struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	SSHKey     string `json:"ssh_key"`     // Optional
	SigningKey string `json:"signing_key"` // Optional (GPG/SSH Signing)
}

type Config map[string]Profile

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	config := loadConfig()

	switch command {
	case "list":
		listProfiles(config)
	case "status":
		showStatus(config)
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: git-swap add <profile-name>")
			os.Exit(1)
		}
		addProfile(os.Args[2], config)
	case "edit":
		if len(os.Args) < 3 {
			fmt.Println("Usage: git-swap edit <profile-name>")
			os.Exit(1)
		}
		editProfile(os.Args[2], config)
	case "remove", "rm":
		if len(os.Args) < 3 {
			fmt.Println("Usage: git-swap remove <profile-name>")
			os.Exit(1)
		}
		removeProfile(os.Args[2], config)
	case "help":
		printUsage()
	default:
		// Unknown command is treated as a profile swap request
		swapProfile(command, config)
	}
}

func getConfigPath() string {
	usr, _ := user.Current()
	return filepath.Join(usr.HomeDir, ".git-swap-config.json")
}

func loadConfig() Config {
	configPath := getConfigPath()
	configFile, err := os.ReadFile(configPath)

	// If file doesn't exist, return empty config
	if os.IsNotExist(err) {
		return make(Config)
	}

	if err != nil {
		fmt.Printf("%sError reading config file: %s%s\n", ColorRed, err, ColorReset)
		os.Exit(1)
	}

	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		fmt.Printf("%sJSON Format Error: %s%s\n", ColorRed, err, ColorReset)
		os.Exit(1)
	}
	return config
}

func saveConfig(config Config) {
	configPath := getConfigPath()
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Printf("%sError processing data: %s%s\n", ColorRed, err, ColorReset)
		os.Exit(1)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		fmt.Printf("%sError saving file: %s%s\n", ColorRed, err, ColorReset)
		os.Exit(1)
	}
}

func addProfile(key string, config Config) {
	if _, exists := config[key]; exists {
		fmt.Printf("%sProfile '%s' already exists. Use 'edit' to modify it.%s\n", ColorRed, key, ColorReset)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Enter Name: ")
	name, _ := reader.ReadString('\n')

	fmt.Printf("Enter Email: ")
	email, _ := reader.ReadString('\n')

	fmt.Printf("Enter SSH Key Path (Optional, press Enter to skip): ")
	sshKey, _ := reader.ReadString('\n')

	fmt.Printf("Enter Signing Key (GPG Key ID or SSH Public Key) (Optional): ")
	signingKey, _ := reader.ReadString('\n')

	config[key] = Profile{
		Name:       strings.TrimSpace(name),
		Email:      strings.TrimSpace(email),
		SSHKey:     strings.TrimSpace(sshKey),
		SigningKey: strings.TrimSpace(signingKey),
	}

	saveConfig(config)
	fmt.Printf("%s✅ Profile '%s' added successfully!%s\n", ColorGreen, key, ColorReset)
}

func editProfile(key string, config Config) {
	profile, exists := config[key]
	if !exists {
		fmt.Printf("%sError: Profile '%s' does not exist.%s\n", ColorRed, key, ColorReset)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("💡 Tip: Press Enter to keep current value. Type '-' to clear.")

	fmt.Printf("Enter Name [%s]: ", profile.Name)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if name == "-" { profile.Name = "" } else if name != "" { profile.Name = name }

	fmt.Printf("Enter Email [%s]: ", profile.Email)
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)
	if email == "-" { profile.Email = "" } else if email != "" { profile.Email = email }

	fmt.Printf("Enter SSH Key Path [%s]: ", profile.SSHKey)
	sshKey, _ := reader.ReadString('\n')
	sshKey = strings.TrimSpace(sshKey)
	if sshKey == "-" { profile.SSHKey = "" } else if sshKey != "" { profile.SSHKey = sshKey }

	fmt.Printf("Enter Signing Key [%s]: ", profile.SigningKey)
	signingKey, _ := reader.ReadString('\n')
	signingKey = strings.TrimSpace(signingKey)
	if signingKey == "-" {
		profile.SigningKey = ""
	} else if signingKey != "" {
		profile.SigningKey = signingKey
	}

	config[key] = profile
	saveConfig(config)
	fmt.Printf("%s✅ Profile '%s' updated successfully!%s\n", ColorGreen, key, ColorReset)
}

func removeProfile(key string, config Config) {
	if _, exists := config[key]; !exists {
		fmt.Printf("%sError: Profile '%s' does not exist.%s\n", ColorRed, key, ColorReset)
		os.Exit(1)
	}

	delete(config, key)
	saveConfig(config)
	fmt.Printf("%s🗑️ Profile '%s' removed successfully.%s\n", ColorGreen, key, ColorReset)
}

func listProfiles(config Config) {
	if len(config) == 0 {
		fmt.Println("No profiles found. Use 'git-swap add <name>' to create one.")
		return
	}

	fmt.Println("Available Identities:")
	keys := make([]string, 0, len(config))
	for k := range config {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		p := config[k]
		fmt.Printf(" 🔄 %s%s%s (%s)\n", ColorCyan, k, ColorReset, p.Email)
	}
}

func showStatus(config Config) {
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		fmt.Printf("%sError: Not a git repository.%s\n", ColorRed, ColorReset)
		os.Exit(1)
	}

	cmdName := exec.Command("git", "config", "user.name")
	outName, _ := cmdName.Output()
	currentName := strings.TrimSpace(string(outName))

	cmdEmail := exec.Command("git", "config", "user.email")
	outEmail, _ := cmdEmail.Output()
	currentEmail := strings.TrimSpace(string(outEmail))

	fmt.Println("Current Git Config (Local):")
	fmt.Printf("  Name:  %s%s%s\n", ColorYellow, currentName, ColorReset)
	fmt.Printf("  Email: %s%s%s\n", ColorYellow, currentEmail, ColorReset)

	matchedProfile := ""
	for key, p := range config {
		if p.Name == currentName && p.Email == currentEmail {
			matchedProfile = key
			break
		}
	}

	if matchedProfile != "" {
		fmt.Printf("✅ This matches profile: %s%s%s\n", ColorGreen, matchedProfile, ColorReset)
	} else {
		fmt.Printf("%s⚠️  No matching profile found in git-swap config.%s\n", ColorYellow, ColorReset)
	}
}

func swapProfile(profileName string, config Config) {
	profile, ok := config[profileName]
	if !ok {
		fmt.Printf("%sError: Profile '%s' not found.%s\n", ColorRed, profileName, ColorReset)
		fmt.Println("To create it: git-swap add", profileName)
		os.Exit(1)
	}

	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		fmt.Printf("%sError: Not a git repository.%s\n", ColorRed, ColorReset)
		os.Exit(1)
	}

	setGitConfig("user.name", profile.Name)
	setGitConfig("user.email", profile.Email)

	if profile.SSHKey != "" {
		sshCmd := fmt.Sprintf("ssh -i %s -o IdentitiesOnly=yes -F /dev/null", profile.SSHKey)
		setGitConfig("core.sshCommand", sshCmd)
		fmt.Printf("🔑 SSH Key locked: %s\n", profile.SSHKey)
	} else {
		unsetGitConfig("core.sshCommand")
	}

	if profile.SigningKey != "" {
		setGitConfig("user.signingkey", profile.SigningKey)
		setGitConfig("commit.gpgsign", "true")

		if strings.HasPrefix(profile.SigningKey, "ssh-") {
			setGitConfig("gpg.format", "ssh")
		} else {
			unsetGitConfig("gpg.format")
		}
		fmt.Printf("🔏 Commit Signing: Enabled (%s)\n", profile.SigningKey)
	} else {
		unsetGitConfig("user.signingkey")
		setGitConfig("commit.gpgsign", "false") 
		unsetGitConfig("gpg.format")
	}

	fmt.Printf("%s✅ Swapped to: %s%s\n", ColorGreen, profileName, ColorReset)
	fmt.Printf("👤 %s <%s>\n", profile.Name, profile.Email)
}

func setGitConfig(key, value string) {
	cmd := exec.Command("git", "config", "--local", key, value)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to set git config (%s): %s\n", key, err)
	}
}

func unsetGitConfig(key string) {
	exec.Command("git", "config", "--local", "--unset", key).Run()
}

func printUsage() {
	fmt.Println("git-swap: Manage git identities directly from CLI.")
	fmt.Println("\nUsage:")
	fmt.Println("  git-swap list              -> List all profiles")
	fmt.Println("  git-swap status            -> Show current repo identity")
	fmt.Println("  git-swap add <name>        -> Add a new profile")
	fmt.Println("  git-swap edit <name>       -> Edit an existing profile")
	fmt.Println("  git-swap remove <name>     -> Delete a profile")
	fmt.Println("  git-swap <name>            -> Apply profile to current repo")
}