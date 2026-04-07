package cmd

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/user/gitmap/constants"
)

// runHelpDashboard serves the docs site locally.
func runHelpDashboard(args []string) {
	checkHelp("help-dashboard", args)

	port := parseHelpDashboardFlags(args)
	binaryDir := resolveBinaryDir()
	docsDir := filepath.Join(binaryDir, constants.HDDocsDir)

	if _, err := os.Stat(docsDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, constants.ErrHDNoDocsDir, docsDir)
		os.Exit(1)
	}

	distDir := filepath.Join(docsDir, constants.HDDistDir)

	if info, err := os.Stat(distDir); err == nil && info.IsDir() {
		serveStatic(distDir, port)
	} else {
		fmt.Println(constants.MsgHDNoDistFallback)
		serveDev(docsDir, port)
	}
}

// parseHelpDashboardFlags parses the --port flag.
func parseHelpDashboardFlags(args []string) int {
	fs := flag.NewFlagSet(constants.CmdHelpDashboard, flag.ExitOnError)
	port := fs.Int("port", constants.HDDefaultPort, constants.FlagDescHDPort)
	fs.Parse(args)

	return *port
}

// resolveBinaryDir returns the directory containing the gitmap binary.
func resolveBinaryDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}

	resolved, err := filepath.EvalSymlinks(exe)
	if err != nil {
		return filepath.Dir(exe)
	}

	return filepath.Dir(resolved)
}

// serveStatic serves pre-built dist/ files over HTTP.
func serveStatic(distDir string, port int) {
	fmt.Printf(constants.MsgHDServingStatic, distDir, port)
	openBrowser(port)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.FileServer(http.Dir(distDir)),
	}

	go handleShutdown(server)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, constants.ErrHDServe, err)
		os.Exit(1)
	}

	fmt.Print(constants.MsgHDStopped)
}

// serveDev runs npm install + npm run dev as a fallback.
func serveDev(docsDir string, port int) {
	npmPath, err := exec.LookPath("npm")
	if err != nil {
		fmt.Fprint(os.Stderr, constants.ErrHDNPMNotFound)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgHDRunningNPM)

	install := exec.Command(npmPath, "install")
	install.Dir = docsDir
	install.Stdout = os.Stdout
	install.Stderr = os.Stderr

	if err := install.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrHDNPMInstall, err)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgHDStartingDev, docsDir)

	dev := exec.Command(npmPath, "run", "dev", "--", "--port", fmt.Sprintf("%d", port))
	dev.Dir = docsDir
	dev.Stdout = os.Stdout
	dev.Stderr = os.Stderr

	if err := dev.Start(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrHDDevServer, err)
		os.Exit(1)
	}

	openBrowser(port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	dev.Process.Kill()
	fmt.Print(constants.MsgHDStopped)
}

// openBrowser opens the URL in the default browser.
func openBrowser(port int) {
	url := fmt.Sprintf("http://localhost:%d", port)
	fmt.Printf(constants.MsgHDOpening, port)

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case constants.OSWindows:
		cmd = exec.Command(constants.CmdWindowsShell, constants.CmdArgSlashC, constants.CmdArgStart, url)
	case constants.OSDarwin:
		cmd = exec.Command(constants.CmdOpen, url)
	default:
		cmd = exec.Command(constants.CmdXdgOpen, url)
	}

	cmd.Start()
}

// handleShutdown gracefully stops the static server on Ctrl+C.
func handleShutdown(server *http.Server) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	server.Close()
}
