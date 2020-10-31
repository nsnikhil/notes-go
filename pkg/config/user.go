package config

type passwordConfig struct {
	saltLength, iterations, keyLength int
}

func newPasswordConfig() passwordConfig {
	return passwordConfig{
		saltLength: getInt("USER_PASSWORD_HASH_SALT_LENGTH"),
		iterations: getInt("USER_PASSWORD_HASH_ITERATIONS"),
		keyLength:  getInt("USER_PASSWORD_HASH_KEY_LENGTH"),
	}
}

type authenticationConfig struct {
	pemString string
}

func newAuthenticationConfig() authenticationConfig {
	return authenticationConfig{
		pemString: getString("USER_TOKEN_SIGNING_KEY"),
	}
}

type UserConfig struct {
	passwordConfig
	authenticationConfig
}

func newUserConfig() UserConfig {
	return UserConfig{
		newPasswordConfig(),
		newAuthenticationConfig(),
	}
}

func (uc UserConfig) SaltLength() int {
	return uc.saltLength
}

func (uc UserConfig) Iterations() int {
	return uc.iterations
}

func (uc UserConfig) KeyLength() int {
	return uc.keyLength
}

func (uc UserConfig) PemString() string {
	return uc.pemString
}
