package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/viper"
	"github.com/thomas-fossati/trafic/config"
	"github.com/thomas-fossati/trafic/runner"
)

type StatusReport struct {
	Label string
	Role  runner.Role
	Error error
}

type RunnersMap map[string]runner.Runner

var R RunnersMap

func NewLogger(tag string) (*log.Logger, error) {
	return log.New(os.Stderr, tag, log.LstdFlags|log.LUTC|log.Lshortfile), nil
}

func loadFlows(dir string) ([]config.FlowConfig, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Printf("cannot read dir: %s: %v", dir, err)
		return nil, err
	}

	var flows []config.FlowConfig

	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		fqn := filepath.Join(dir, file.Name())

		flow, err := config.NewFlowConfigFromFile(fqn)
		if err != nil {
			log.Printf("cannot parse %s: %v", fqn, err)
			return nil, err
		}

		flows = append(flows, *flow)
	}

	return flows, nil
}

func sortFlowsByDeadline(flows []config.FlowConfig) {
	sort.Slice(flows, func(i, j int) bool {
		return flows[i].When[0] < flows[j].When[0]
	})
}

func configurerForRole(role runner.Role, flow config.FlowConfig) config.Configurer {
	if role == runner.RoleClient {
		return &flow.ClientCfg
	}

	return &flow.ServerCfg
}

func runners(role runner.Role) {
	log, err := NewLogger(fmt.Sprintf("[%s] ", viper.GetString("log.tag")))
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}

	flows, err := loadFlows(viper.GetString("flows.dir"))
	if err != nil {
		log.Fatalf("cannot load flows: %v", err)
	}

	sortFlowsByDeadline(flows)

	done := make(chan StatusReport)

	go sched(flows, log, done, role)

	if err = wait(done); err != nil {
		log.Printf("waiting: %v", err)
	}
}

func sched(flows []config.FlowConfig, log *log.Logger, done chan StatusReport,
	role runner.Role) {
	R = make(RunnersMap)

	for _, flow := range flows {
		cfg := configurerForRole(role, flow)

		r, err := runner.NewRunner(role, log, cfg)
		if err != nil {
			log.Fatalf("cannot create %s %s: %v", role, flow.Label, err)
		}

		err = r.Start()
		if err != nil {
			log.Fatalf("cannot start %s %s: %v", role, flow.Label, err)
		}

		R[flow.Label] = *r

		// Start watchdog for this iperf3 instance
		go watchdog(*r, flow.Label, done)
	}
}

func watchdog(r runner.Runner, label string, done chan StatusReport) {
	err := r.Wait()
	if err != nil {
		log.Printf("cannot reap %s %s: %v", r.Role, label, err)
	}

	delete(R, label)

	done <- StatusReport{label, r.Role, err}
}

func wait(done chan StatusReport) error {
	// XXX if server and not one-off, there's no way this loop can break
	// XXX we need another (interactive - keyboard, signals) source of
	// XXX events
	for {
		select {
		case s := <-done:
			// If one has failed, flag the test as invalid and bail out
			if s.Error != nil {
				tearDownRunners()
				return fmt.Errorf("%v %s failed: %v", s.Role, s.Label, s.Error)
			}

			log.Printf("%v %s finished ok", s.Role, s.Label)

			if len(R) == 0 {
				log.Printf("all %v(s) finished ok", s.Role)
				return nil
			}

			log.Printf("%d %v(s) to go", len(R), s.Role)
		}
	}
}

func tearDownRunners() {
	log.Printf("TODO(tho) tearDownRunners")
}
