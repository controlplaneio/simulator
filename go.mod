module github.com/kubernetes-simulator/simulator

require (
	github.com/aws/aws-sdk-go v1.21.7
	github.com/fatih/structs v1.1.0
	github.com/glendc/go-external-ip v0.0.0-20170425150139-139229dcdddd
	github.com/go-test/deep v1.0.5
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/olekukonko/tablewriter v0.0.4
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.3.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/text v0.3.8 // indirect
	gopkg.in/yaml.v2 v2.2.2
)

replace github.com/go-critic/go-critic@v0.0.0-20181204210945-1df300866540 => github.com/go-critic/go-critic v0.3.5-0.20190526074819-1df300866540

go 1.13
