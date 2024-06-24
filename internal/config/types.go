package config

type (
	Config struct {
		Server      ServerConfig      `yaml:"server"`
		Database    DatabaseConfig    `yaml:"database"`
		FeatureFlag FeatureFlagConfig `yaml:"featureFlag"`
		Jwt         JwtConfig         `yaml:"jwt"`
	}

	ServerConfig struct {
		Port               string `yaml:"port"`
		ShutdownTimeMillis int64  `yaml:"shutdownTimeMillis"`
		SecretKey          string `yaml:"secretKey"`
		Env                string `yaml:"env"`
	}

	DatabaseConfig struct {
		Postgres PostgresConfig `yaml:"postgres"`
	}

	PostgresConfig struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
		Ssl      string `yaml:"ssl"`
	}

	FeatureFlagConfig struct {
		EnableMigrations bool `yaml:"enableMigrations"`
	}

	JwtConfig struct {
		AccessTokenExpiryHours  int64 `yaml:"accessTokenExpiryHours"`
		RefreshTokenExpiryHours int64 `yaml:"refreshTokenExpiryHours"`
	}
)
