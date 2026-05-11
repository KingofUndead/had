package core

import (
	"bytes"
	"encoding/pem"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/elazarl/goproxy"
)

func InstallCertificate() error {
	ca := goproxy.GoproxyCa
	certBuffer := new(bytes.Buffer)
	pem.Encode(certBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: ca.Certificate[0]})

	certPath := "Fad.crt"
	if err := os.WriteFile(certPath, certBuffer.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to save certificate: %v", err)
	}

	fmt.Printf("✅ Certificate saved to: %s\n", certPath)

	switch runtime.GOOS {
	case "windows":
		return installWindowsCert(certPath)
	case "darwin":
		return installMacCert(certPath)
	case "linux":
		return installLinuxCert(certPath)
	default:
		return fmt.Errorf("unsupported OS for auto-installation: %s", runtime.GOOS)
	}
}

func installWindowsCert(certPath string) error {
	cmd := exec.Command("certutil", "-addstore", "-user", "Root", certPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		cmd = exec.Command("powershell", "-Command", 
			fmt.Sprintf("Import-Certificate -FilePath \"%s\" -CertStoreLocation Cert:\\CurrentUser\\Root", certPath))
		output, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to install certificate on Windows: %v\n%s", err, output)
		}
	}
	fmt.Println("✅ Certificate installed to Windows Trust Store")
	return nil
}

func installMacCert(certPath string) error {
	cmd := exec.Command("security", "add-trusted-cert", "-d", "-r", "trustRoot", "-k", "/Library/Keychains/System.keychain", certPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		cmd = exec.Command("security", "add-trusted-cert", "-d", "-r", "trustRoot", "-k", "~/Library/Keychains/login.keychain-db", certPath)
		output, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to install certificate on macOS: %v\n%s", err, output)
		}
	}
	fmt.Println("✅ Certificate installed to macOS Keychain")
	return nil
}

func installLinuxCert(certPath string) error {
	cmd := exec.Command("sudo", "cp", certPath, "/usr/local/share/ca-certificates/Fad.crt")
	output, err := cmd.CombinedOutput()
	if err != nil {
		homeDir, _ := os.UserHomeDir()
		certDir := homeDir + "/.local/share/ca-certificates"
		os.MkdirAll(certDir, 0755)
		cmd = exec.Command("cp", certPath, certDir+"/Fad.crt")
		output, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to copy certificate: %v\n%s", err, output)
		}
		
		cmd = exec.Command("update-ca-certificates", "--fresh")
		output, err = cmd.CombinedOutput()
		if err != nil {
			cmd = exec.Command("sudo", "update-ca-certificates")
			output, err = cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("⚠️  Certificate copied but may need manual trust: %v\n", err)
			}
		}
	} else {
		cmd = exec.Command("sudo", "update-ca-certificates")
		output, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to update CA certificates: %v\n%s", err, output)
		}
	}
	
	fmt.Println("✅ Certificate installed to Linux CA store")
	return nil
}

func ShowManualInstructions() {
	certPath := "Fad.crt"
	
	fmt.Printf("\n%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", "\033[33m", "\033[0m")
	fmt.Printf("%s                    MANUAL CERTIFICATE INSTALLATION%s\n", "\033[1m", "\033[0m")
	fmt.Printf("%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", "\033[33m", "\033[0m")
	fmt.Printf("Certificate file: %s%s%s\n", "\033[32m", certPath, "\033[0m")
	fmt.Println()
	
	switch runtime.GOOS {
	case "windows":
		fmt.Println("Windows Installation:")
		fmt.Println("  1. Double-click Fad.crt")
		fmt.Println("  2. Click 'Install Certificate'")
		fmt.Println("  3. Select 'Current User' or 'Local Machine'")
		fmt.Println("  4. Select 'Place all certificates in the following store'")
		fmt.Println("  5. Click 'Browse' and select 'Trusted Root Certification Authorities'")
		fmt.Println("  6. Click 'OK' then 'Next' then 'Finish'")
		
	case "darwin":
		fmt.Println("macOS Installation:")
		fmt.Println("  1. Double-click Fad.crt")
		fmt.Println("  2. Keychain Access will open")
		fmt.Println("  3. Find 'Fad' certificate")
		fmt.Println("  4. Double-click it")
		fmt.Println("  5. Expand 'Trust' section")
		fmt.Println("  6. Set 'SSL' to 'Always Trust'")
		fmt.Println("  7. Close the window (password required)")
		
	case "linux":
		fmt.Println("Linux Installation:")
		fmt.Println("  Option 1 (Ubuntu/Debian):")
		fmt.Println("    sudo cp Fad.crt /usr/local/share/ca-certificates/")
		fmt.Println("    sudo update-ca-certificates")
		fmt.Println()
		fmt.Println("  Option 2 (Fedora/RHEL):")
		fmt.Println("    sudo cp Fad.crt /etc/pki/ca-trust/source/anchors/")
		fmt.Println("    sudo update-ca-trust")
		fmt.Println()
		fmt.Println("  Option 3 (Arch Linux):")
		fmt.Println("    sudo cp Fad.crt /etc/ca-certificates/trust-source/anchors/")
		fmt.Println("    sudo trust extract-compat")
	}
	
	fmt.Printf("%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", "\033[33m", "\033[0m")
}
