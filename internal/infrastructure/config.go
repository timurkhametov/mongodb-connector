package infrastructure

type DBConfig struct {
	Host     string
	Port     uint16
	Database string
	Username string
	Password string
}

type MigrationConfig struct {
	DBConfig
	CollectionName     string
	UseTransactionMode bool
	AdvisorLock        *AdvisorLockConfig
}

type AdvisorLockConfig struct {
	CollectionName string
	Timeout        int
	Interval       int
}
