package simulator

import (
	"github.com/kubernetes-simulator/simulator/pkg/progress"
	"github.com/kubernetes-simulator/simulator/pkg/ssh"
	"github.com/kubernetes-simulator/simulator/pkg/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"os"
)

// Simulator represents a session with simulator and holds all the configuration
// necessary to run simulator
type Simulator struct {
	// Logger is the logger the simulator will use
	Logger *logrus.Logger
	// TfDir is the path to the terraform code used to standup the simulator cluster
	TfDir string
	// BucketName is the remote state bucket to use for terraform
	BucketName string
	// AttackTag is the docker tag for the attack container that terraform will use
	// when creating the infrastructure: e.g. latest
	AttackTag string
	// AttackRepo is the docker repo for the attack container that terraform will use
	// when creating the infrastructure: e.g. controlplane/simulator-attack
	AttackRepo string
	// scenarioID is the unique identifier of the scenario used for the launch function
	ScenarioID string
	// TfVarsDir is the location to store the terraform variables file that are detected
	// automatically for use when creating the infrastructure
	TfVarsDir string
	// ScenariosDir is the location of the scenarios for perturb to use when perturbing
	// the cluster
	ScenariosDir string
	// disableIPDetection enables IP checks used for cidr access. Enabled by default.
	DisableIPDetection bool
	// Extra CIDRs to be added to the bastion security group to allow SSH from arbitrary
	// locations
	ExtraCIDRs string
	// SSHStateProvider manages retrieving and persisting SSH state
	SSHStateProvider ssh.StateProvider
	// ProgressStateProvider manages retrieving and persisting progress state
	ProgressStateProvider progress.StateProvider
	// sshLogger is the file logger for the SSH related components
	SSHLogger *logrus.Logger
}

// Option is a type used to configure a `Simulator` instance
type Option func(*Simulator)

// NewSimulator constructs a new instance of `Simulator`
func NewSimulator(options ...Option) *Simulator {
	simulator := Simulator{}

	for _, option := range options {
		option(&simulator)
	}

	if simulator.SSHStateProvider == nil {
		simulator.SSHStateProvider = ssh.LocalStateProvider{}
	}

	if simulator.SSHLogger == nil {
		logpath := util.MustExpandTilde("~/.kubesim/ssh.log")
		file, err := os.OpenFile(logpath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			panic(errors.Wrapf(err, "Error opening %v for SSH logging", logpath))
		}

		simulator.SSHLogger = logrus.New()
		simulator.SSHLogger.SetOutput(file)

	}

	if simulator.ProgressStateProvider == nil {
		simulator.ProgressStateProvider = progress.NewLocalStateProvider(
			simulator.SSHLogger)
	}

	return &simulator
}

// WithLogger returns a configurer for creating a `Simulator` instance with
// `NewSimulator`
func WithLogger(logger *logrus.Logger) Option {
	return func(s *Simulator) {
		s.Logger = logger
	}
}

// WithSSHLogger returns a configurer for creating a `Simulator` instance with
// `NewSimulator`
func WithSSHLogger(logger *logrus.Logger) Option {
	return func(s *Simulator) {
		s.SSHLogger = logger
	}
}

// WithAttackTag returns a configurer for creating a `Simulator` instance with
// `NewSimulator`
func WithAttackTag(attackTag string) Option {
	return func(s *Simulator) {
		s.AttackTag = attackTag
	}
}

// WithAttackRepo returns a configurer for creating a `Simulator` instance with
// `NewSimulator`
func WithAttackRepo(attackRepo string) Option {
	return func(s *Simulator) {
		s.AttackRepo = attackRepo
	}
}

// WithTfDir returns a configurer for creating a `Simulator` instance with
// `NewSimulator`
func WithTfDir(tfDir string) Option {
	return func(s *Simulator) {
		s.TfDir = tfDir
	}
}

// WithTfVarsDir returns a configurer for creating a `Simulator` instance with
// `NewSimulator`
func WithTfVarsDir(tfVarsDir string) Option {
	return func(s *Simulator) {
		s.TfVarsDir = tfVarsDir
	}
}

// WithScenarioID returns a configurer for creating a `Simulator` instance with
// `NewSimulator`
func WithScenarioID(scenarioID string) Option {
	return func(s *Simulator) {
		s.ScenarioID = scenarioID
	}
}

// WithScenariosDir returns a configurer for creating a `Simulator` instance with
// `NewSimulator`
func WithScenariosDir(scenariosDir string) Option {
	return func(s *Simulator) {
		s.ScenariosDir = scenariosDir
	}
}

// WithBucketName returns a configurer for creating a `Simulator` instance with
// `NewSimulator`
func WithBucketName(bucketName string) Option {
	return func(s *Simulator) {
		s.BucketName = bucketName
	}
}

// WithoutIPDetection returns a configurer for creating a `Simulator` instance with
// `NewSimulator`
func WithoutIPDetection(disableIPDetection bool) Option {
	return func(s *Simulator) {
		s.DisableIPDetection = disableIPDetection
	}
}

// WithExtraCIDRs returns a configurer for creating a `Simulator` instance with
// `NewSimulator`
func WithExtraCIDRs(extraCIDRs string) Option {
	return func(s *Simulator) {
		s.ExtraCIDRs = extraCIDRs
	}
}

// WithSSHStateProvider returns a configurer for creating a `Simulator` instance with
// `NewSimulator`
func WithSSHStateProvider(stateProvider ssh.StateProvider) Option {
	return func(s *Simulator) {
		s.SSHStateProvider = stateProvider
	}
}

// WithProgressStateProvider returns a configurer for creating a `Simulator` instance with
// `NewSimulator`
func WithProgressStateProvider(stateProvider progress.StateProvider) Option {
	return func(s *Simulator) {
		s.ProgressStateProvider = stateProvider
	}
}
