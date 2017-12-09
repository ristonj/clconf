package clconf

import (
	"fmt"

	"github.com/urfave/cli"
)

func Cgetv(c *cli.Context) error {
	path := c.Args().First()
	value, ok := GetValue(path, load(c))
	if !ok {
		return cli.NewExitError(
			fmt.Sprintf("[%v] does not exist", path), 1)
	}
	secretAgent, err := newSecretAgentFromCli(c)
	if err != nil {
		return err
	}
	decrypted, err := secretAgent.Decrypt(value.(string))
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return dump(decrypted)
}

func dump(conf interface{}) cli.ExitCoder {
	yaml, err := MarshalYaml(conf)
	if err != nil {
	    return cli.NewExitError(
			fmt.Sprintf("Unable to dump value: %v", err), 1)
	}
	fmt.Println(yaml)
	return nil
}

func Getv(c *cli.Context) error {
	path := c.Args().First()
	if value, ok := GetValue(path, load(c)); ok {
		return dump(value)
	}
	return cli.NewExitError(
		fmt.Sprintf("[%v] does not exist", path), 1)
}

func load(c *cli.Context) map[interface{}]interface{} {
	return LoadConfFromEnvironment(
		c.GlobalStringSlice("yaml-file"),
		c.GlobalStringSlice("yaml-vars"))
}

func newSecretAgentFromCli(c *cli.Context) (*SecretAgent, error) {
	if secretKeysBase64 := c.GlobalString("secret-keys"); secretKeysBase64 != "" {
        return NewSecretAgentFromBase64(secretKeysBase64);
	}
	if secretKeysFile := c.GlobalString("secret-keys-file"); secretKeysFile != "" {
		return NewSecretAgentFromFile(secretKeysFile);
	}
	return nil, cli.NewExitError("--secret-keys or --secret-keys-file required", 1)
}