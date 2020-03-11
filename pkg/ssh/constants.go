package ssh

// KnownHostsPath is the path to write the simulator SSH config file
const KnownHostsPath = "~/.ssh/cp_simulator_known_hosts"

// ConfigPath is the path to write the simulator SSH config file
const ConfigPath = "~/.ssh/cp_simulator_config"

// PrivateKeyPath is the path to the key to be generated and used by simulator
const PrivateKeyPath = "~/.ssh/cp_simulator_rsa"

// PublicKeyPath is the path to the key to be generated and used by simulator
const PublicKeyPath = PrivateKeyPath + ".pub"
