package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/prometheus/common/log"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

var (
	installService   = flag.Bool("service.install", false, "Install as Windows service and exit")
	uninstallService = flag.Bool("service.uninstall", false, "Uninstall Windows service and exit")
	serviceName      = flag.String("service.name", "stream_exporter", "Name of Windows service to run, install or uninstall")
)

func init() {
	startupTasks["processServiceInstallationFlags"] = processServiceInstallationFlags
	startupTasks["runAsService"] = runAsServiceIfNoninteractive
}

func processServiceInstallationFlags() error {
	if *installService {
		if *serviceName == "" {
			return fmt.Errorf("Cannot install service, name not set")
		}
		err := install(*serviceName)
		if err != nil {
			return err
		} else {
			log.Info("Successfully installed service")
			os.Exit(0)
		}
	} else if *uninstallService {
		if *serviceName == "" {
			return fmt.Errorf("Cannot uninstall service, name not set")
		}
		err := uninstall(*serviceName)
		if err != nil {
			return err
		} else {
			log.Info("Successfully uninstalled service")
			os.Exit(0)
		}
	}

	return nil
}

func runAsServiceIfNoninteractive() error {
	if *installService || *uninstallService {
		return nil
	}

	isInteractive, err := svc.IsAnInteractiveSession()
	if err != nil {
		return err
	}

	if !isInteractive {
		if *serviceName == "" {
			return fmt.Errorf("Cannot start service, name not set")
		}
		log.Debugf("Starting service %s", *serviceName)
		go svc.Run(*serviceName, &streamExporterService{stopCh: quitSig})
	}

	return nil
}

type streamExporterService struct {
	stopCh chan<- os.Signal
}

func (s *streamExporterService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				s.stopCh <- os.Interrupt
				break loop
			default:
				log.Error(fmt.Sprintf("Unexpected control request #%d", c))
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

// ServiceManager integration adapted from https://github.com/martinlindhe/winservice

func install(name string) error {
	exepath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("Could not determine own path, %s", err)
	}

	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(name)
	if err == nil {
		s.Close()
		return fmt.Errorf("Service %s already exists", name)
	}
	config := mgr.Config{
		DisplayName: name,
		StartType:   mgr.StartAutomatic,
		Description: "Prometheus stream_exporter",
	}
	args := getServiceArgs(os.Args)
	s, err = m.CreateService(name, exepath, config, args...)
	if err != nil {
		return err
	}
	defer s.Close()

	err = eventlog.InstallAsEventCreate(name, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		s.Delete()
		return fmt.Errorf("Install() failed: %s", err)
	}
	return nil
}

// Remove first arg (program name), and -service.install flag
// Ensure -service.name is included
func getServiceArgs(inputArgs []string) []string {
	args := make([]string, len(inputArgs)-1)
	copy(args, inputArgs[1:])
	foundServiceName := false

	for idx, arg := range args {
		switch arg {
		case "-service.install":
			args = append(args[:idx], args[idx+1:]...)
		case "-service.name":
			foundServiceName = true
		}
	}

	if !foundServiceName {
		args = append(args, "-service.name")
		args = append(args, "stream_exporter")
	}

	args = append(args, "-log.format")
	args = append(args, "logger:eventlog?name=stream_exporter")

	return args
}

func uninstall(name string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("Service %q is not installed", name)
	}
	defer s.Close()

	status, err := s.Query()
	if err != nil {
		return fmt.Errorf("Could not query for service %q state: %v", name, err)
	}

	if status.State != svc.Stopped {
		_, err := s.Control(svc.Stop)
		if err != nil {
			return fmt.Errorf("Could not send stop signal to service %q", err)
		}
	}

	err = s.Delete()
	if err != nil {
		return err
	}
	err = eventlog.Remove(name)
	if err != nil {
		return fmt.Errorf("Remove() failed: %s", err)
	}
	return nil
}
