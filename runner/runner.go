package runner

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/thomas-fossati/trafic/config"
	"github.com/thomas-fossati/trafic/cruncher"
)

type Runner struct {
	Role     Role
	Command  *exec.Cmd
	At       time.Duration
	Label    string
	Logger   *log.Logger
	Stdout   *bytes.Buffer
	Stderr   *bytes.Buffer
	Cruncher cruncher.Cruncher
}

func NewRunner(role Role, log *log.Logger, at time.Duration, label string, cfg config.Configurer, crunch cruncher.Cruncher) (*Runner, error) {
	args, err := cfg.ToArgs()
	if err != nil {
		return nil, err
	}

	return &Runner{
		Role:     role,
		Command:  exec.Command("iperf3", args...),
		At:       at,
		Label:    label,
		Logger:   log,
		Stdout:   &bytes.Buffer{},
		Stderr:   &bytes.Buffer{},
		Cruncher: crunch,
	}, nil
}

func (r *Runner) Start() error {
	r.Logger.Printf("Starting %s %s\n",
		r.Command.Path, strings.Join(r.Command.Args[1:], " "))

	// we want to capture output and error
	r.Command.Stdout = r.Stdout
	r.Command.Stderr = r.Stderr

	return r.Command.Start()
}

func (r *Runner) Wait() ([]byte, error) {
	r.Logger.Printf("Waiting for %v (PID=%v) to complete\n",
		r.Command.Path, r.Command.Process.Pid)

	err := r.Command.Wait()
	if err != nil {
		return nil, err
	}

	// Run test data through a filter if we have been configured to do so
	if r.Cruncher != nil {
		return r.Cruncher.Crunch(r.Stdout.Bytes())
	}

	return r.Stdout.Bytes(), nil
}

func (r *Runner) Kill() error {
	r.Logger.Printf("Killing %v (PID=%v)\n",
		r.Command.Path, r.Command.Process.Pid)

	return r.Command.Process.Kill()
}
