package runner

import (
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/thomas-fossati/trafic/config"
)

type Runner struct {
	Role    Role
	Command *exec.Cmd
	At      time.Duration
	Label   string
	Logger  *log.Logger
}

func NewRunner(role Role, log *log.Logger, at time.Duration, label string, cfg config.Configurer) (*Runner, error) {
	args, err := cfg.ToArgs()
	if err != nil {
		return nil, err
	}

	return &Runner{
		Role:    role,
		Command: exec.Command("iperf3", args...),
		At:      at,
		Label:   label,
		Logger:  log,
	}, nil
}

func (r *Runner) Start() error {
	r.Logger.Printf("Starting %s %s\n",
		r.Command.Path, strings.Join(r.Command.Args[1:], " "))

	return r.Command.Start()
}

func (r *Runner) Wait() error {
	r.Logger.Printf("Waiting for %v (PID=%v) to complete\n",
		r.Command.Path, r.Command.Process.Pid)

	return r.Command.Wait()
}

func (r *Runner) Kill() error {
	r.Logger.Printf("Killing %v (PID=%v)\n",
		r.Command.Path, r.Command.Process.Pid)

	return r.Command.Process.Kill()
}
