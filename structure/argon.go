package structure

type ArgonConfig struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

var EncryptConfig = configSet()

func configSet() (EncryptConfig ArgonConfig) {
	EncryptConfig.Memory = 64 * 1024
	EncryptConfig.Iterations = 3
	EncryptConfig.Parallelism = 2
	EncryptConfig.SaltLength = 16
	EncryptConfig.KeyLength = 32
	return EncryptConfig
}
