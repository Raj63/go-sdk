package config

import (
	"embed"
	"errors"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v6"
	"go.mozilla.org/sops/v3/decrypt"
)

// ParseAppConfig parses the environment file variables into the interface.
func ParseAppConfig(c interface{}) error {
	return parseEnv(c)
}

// LoadDotenv is a helper function to load the dotenv file into environment variables
func LoadDotenv(embedFS embed.FS, resourcePaths map[string]string) error {
	configFilePath := resourcePaths["configs"] + "/" + os.Getenv("APP_ENV") + ".env"
	envs, err := embedFS.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	if os.Getenv("APP_ENV") != "development" {
		f, err := embedFS.Open(resourcePaths["sops"])
		if err == nil {
			if err := os.Setenv("AWS_PROFILE", os.Getenv("APP_ENV")); err != nil {
				return err
			}

			encryptedEnvs := strings.Trim(string(envs), "\n")
			encryptedEnvs = strings.Trim(encryptedEnvs, " ")
			envs, err = decrypt.Data([]byte(encryptedEnvs), "dotenv")
			if err != nil {
				return errors.New("unable to decrypt '" + configFilePath + "' with the specified AWS KMS key, please ensure that your AWS credential is configured properly with non-expired session")
			}
		}

		if f != nil {
			defer f.Close()
		}
	}

	if len(envs) == 0 {
		return nil
	}

	lines := strings.Split(string(envs), "\n")
	for _, line := range lines {
		if line != "" {
			splits := strings.SplitN(line, "=", 2)
			if os.Getenv(splits[0]) == "" {
				if err := os.Setenv(splits[0], splits[1]); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func parseEnv(c interface{}) error {
	if err := env.ParseWithFuncs(c, map[reflect.Type]env.ParserFunc{
		reflect.TypeOf([]byte{}):            parseByteArray,
		reflect.TypeOf([][]byte{}):          parseByte2DArray,
		reflect.TypeOf(map[string]int{}):    parseMapStrInt,
		reflect.TypeOf(map[string]string{}): parseMapStrStr,
		reflect.TypeOf(http.SameSite(1)):    parseHTTPSameSite,
	}); err != nil {
		return err
	}

	return nil
}

func parseByteArray(v string) (interface{}, error) {
	return []byte(v), nil
}

func parseByte2DArray(v string) (interface{}, error) {
	newBytes := [][]byte{}
	bytes := strings.Split(v, ",")
	for _, b := range bytes {
		newBytes = append(newBytes, []byte(b))
	}

	return newBytes, nil
}

func parseHTTPSameSite(v string) (interface{}, error) {
	ss, err := strconv.Atoi(v)
	if err != nil {
		return nil, err
	}

	return http.SameSite(ss), nil
}

func parseMapStrInt(v string) (interface{}, error) {
	newMaps := map[string]int{}
	maps := strings.Split(v, ",")
	for _, m := range maps {
		splits := strings.Split(m, ":")
		if len(splits) != 2 {
			continue
		}

		val, err := strconv.Atoi(splits[1])
		if err != nil {
			return nil, err
		}

		newMaps[splits[0]] = val
	}

	return newMaps, nil
}

func parseMapStrStr(v string) (interface{}, error) {
	newMaps := map[string]string{}
	maps := strings.Split(v, ",")
	for _, m := range maps {
		splits := strings.Split(m, ":")
		if len(splits) != 2 {
			continue
		}

		newMaps[splits[0]] = splits[1]
	}

	return newMaps, nil
}
