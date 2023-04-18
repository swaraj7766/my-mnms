package mnms

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/qeof/q"
)

// direct commands
type directCommands struct {
	GenerateRSAKeyPair bool   // generate rsa key pair
	Decrypt            bool   // decrypt something
	Encrypt            bool   // encrypt something
	Export             bool   // export mnms config
	Import             bool   // import mnms config
	PublickeyPath      string // public key path
	PrivatekeyPath     string // private key path
	ConfigFile         bool   //config file
	MnmsPubkey         bool   // mnms public key
	In                 string // input file
	Out                string // output file
	Name               string // name
}

// getFileWithDefault get io.Writer from filename, if filename is empty, return os.Stdout
func getFileWithDefault(filename string, def *os.File) (*os.File, error) {
	if filename == "" {

		return def, nil
	}

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		q.Q(err)
		return nil, err
	}
	return f, nil
}

// ProcessDirectCommands checks if there is any direct commands
/*
Utility command.

Usage: util -mnmspubkey -out [out_file]
	Get default mnms public key. If [out_file] is empty, output to stdout.
	public key file use for decrypt mnms encrypted config file.

	[out_file]    : output public file name (optional)

Example :
	util -mnmspubkey -out mnms.pub
	util -mnmspubkey > mnms.pub


Usage: util -genrsa
	Generate rsa key pair. default output is
	$HOME/.mnms/id_rsa and $HOME/.mnms/id_rsa.pub

Example :
	util -genrsa


Usage: util -genrsa -name [file_prefix]
	Generate rsa key pair to [file_prefix].pub and [file_prefix]
	[file_prefix]  : output file prefix (optional)

Example :
	util -genrsa -name ~/mnmskey

Usage: util -encrypt -pubkey [pubkey_file] -in [in_file] -out [out_file]
	Encrypt [in_file] with [pubkey_file] and output to [out_file].
	If [out_file] is empty, output to stdout. if [in_file] is empty,
	input from stdin.

	[pubkey_file]  : public key file
	[in_file]      : input file (optional)
	[out_file]     : output file (optional)

Example :
	util -encrypt -pubkey mnms.pub -in mnms.conf -out mnms.conf.enc


Usage: util -decrypt -privkey [privkey_file] -in [in_file] -out [out_file]
	Decrypt [in_file] with [privkey_file] and output to [out_file].
	If [out_file] is empty, output to stdout. if [in_file] is empty,
	input from stdin.

	[privkey_file] : private key file
	[in_file]      : input file (optional)
	[out_file]     : output file (optional)

Example :
	util -decrypt -privkey mnms.key -in mnms.conf.enc -out mnms.conf


Usage: util -export -configfile -pubkey [pubkey_file] -privkey [privkey_file] -out [out_file]
	Export config file to [out_file]. If [out_file] is empty, output to stdout.
	if [privkey_file] is empty, use default private key file to decrypt config.json.

	[pubkey_file]  : public key file
	[privkey_file] : private key file (optional)
	[out_file]     : output file (optional)

Example :
	util -export -configfile -pubkey mnms.pub -out mnms.conf
	util -export -configfile -pubkey mnms.pub -privkey mnms.key -out mnms.conf


Usage: util -import -configfile -in [in_file]
	Import config file from [in_file]. If [in_file] is empty,
	input from stdin. input file must be encrypted by pair public key.

	[in_file]      : input file (optional)

Example :
	util -import -configfile -in mnms.conf
*/
func ProcessDirectCommands() error {

	dc := directCommands{
		//default values
		GenerateRSAKeyPair: false,
		Decrypt:            false,
		Export:             false,
		Import:             false,
	}

	subsetFlags := flag.NewFlagSet("subset", flag.ExitOnError)
	subsetFlags.BoolVar(&dc.GenerateRSAKeyPair, "genrsa", false, "generate rsa key pair")
	subsetFlags.BoolVar(&dc.Decrypt, "decrypt", false, "decrypt something")
	subsetFlags.BoolVar(&dc.Encrypt, "encrypt", false, "encrypt something")
	subsetFlags.BoolVar(&dc.Export, "export", false, "export mnms config")
	subsetFlags.BoolVar(&dc.Import, "import", false, "import mnms config")
	subsetFlags.BoolVar(&dc.ConfigFile, "configfile", false, "config file")
	subsetFlags.BoolVar(&dc.MnmsPubkey, "mnmspubkey", false, "mnms public key")
	subsetFlags.StringVar(&dc.PublickeyPath, "pubkey", "", "public key path")
	subsetFlags.StringVar(&dc.PrivatekeyPath, "privkey", "", "private key path")
	subsetFlags.StringVar(&dc.In, "in", "", "input file")
	subsetFlags.StringVar(&dc.Out, "out", "", "output file")
	subsetFlags.StringVar(&dc.Name, "name", "", "name")
	help := subsetFlags.Bool("help", false, "help")
	if err := subsetFlags.Parse(os.Args[2:]); err != nil {
		return err
	}

	// dump usage
	if *help {
		subsetFlags.Usage()
		return nil
	}

	mnmsFolder, err := CheckMNMSFolder()
	if err != nil {
		return err
	}

	// mnmsctl util -genrsa -pubkey {public_key_file} -privkey {private_key_file}
	if dc.GenerateRSAKeyPair {
		fmt.Fprintln(os.Stderr, "generating rsa key pair...")
		prikeyPath := path.Join(mnmsFolder, "id_rsa")
		pubkeyPath := path.Join(mnmsFolder, "id_rsa.pub")

		if len(dc.Name) > 0 {

			prikeyPath = dc.Name
			pubkeyPath = dc.Name + ".pub"
		}

		fmt.Fprintln(os.Stderr, "output private key to", prikeyPath)
		fmt.Fprintln(os.Stderr, "output public key to", pubkeyPath)

		// generate rsa key pair
		prikey, err := GenerateRSAKeyPair(4096)
		if err != nil {
			return err
		}
		// generate private and public key bytes
		prikeyBytes, err := EndcodePrivateKeyToPEM(prikey)
		if err != nil {
			return err
		}
		// write private key to prikeyPath
		err = ioutil.WriteFile(prikeyPath, prikeyBytes, 0600)
		if err != nil {
			return err
		}

		pubkeyBytes, err := GenerateRSAPublickey(prikey)
		if err != nil {
			return err
		}

		// write public key to pubkeyPath
		err = ioutil.WriteFile(pubkeyPath, pubkeyBytes, 0644)
		if err != nil {
			return err
		}

		return nil
	}

	// mnms util -mnmspubkey
	if dc.MnmsPubkey {
		outputFile, err := getFileWithDefault(dc.Out, os.Stdout)
		if err != nil {
			return err
		}
		if outputFile != os.Stdout {
			defer outputFile.Close()
		}
		pubkeyBytes, err := GenerateOwnPublickey()
		if err != nil {
			return err
		}
		_, err = outputFile.Write(pubkeyBytes)
		if err != nil {
			return err
		}
		return nil
	}

	// mnmsctl util -encrypt -in {plain_file} -out {cipher_file} -pubkey {public_key_file}
	if dc.Encrypt {
		inputFile, err := getFileWithDefault(dc.In, os.Stdin)
		if err != nil {
			return err
		}
		if inputFile != os.Stdin {
			defer inputFile.Close()
		}

		outputFile, err := getFileWithDefault(dc.Out, os.Stdout)
		if err != nil {
			return err
		}
		if outputFile != os.Stdout {
			defer outputFile.Close()
		}

		if dc.PublickeyPath == "" {
			return fmt.Errorf("public key is required, specify public key with -pubkey flag")
		}
		// read public key
		pubkeyBytes, err := ioutil.ReadFile(dc.PublickeyPath)
		if err != nil {
			return err
		}
		// read plain text
		plainBytes, err := ioutil.ReadAll(inputFile)
		if err != nil {
			return err
		}
		// encrypt plain text
		cipherBytes, err := EncryptWithPublicKey(plainBytes, pubkeyBytes)
		if err != nil {
			return err
		}
		// write cipher text to output file
		_, err = outputFile.Write(cipherBytes)
		if err != nil {
			return err
		}
		return nil
	}

	// mnmsctl util -decrypt -in {cipher_file} -out {plain_file} -privkey {private_key_file}
	if dc.Decrypt {
		inputFile, err := getFileWithDefault(dc.In, os.Stdin)
		if err != nil {
			return err
		}
		if inputFile != os.Stdin {
			defer inputFile.Close()
		}

		output, err := getFileWithDefault(dc.Out, os.Stdout)
		if err != nil {
			return err
		}
		if output != os.Stdout {
			defer output.Close()
		}

		if dc.PrivatekeyPath == "" {
			return fmt.Errorf("private key is required")

		}
		// read private key
		prikeyBytes, err := ioutil.ReadFile(dc.PrivatekeyPath)
		if err != nil {
			return err
		}
		// read encrypted pass
		cipherBytes, err := ioutil.ReadFile(dc.In)
		if err != nil {
			return err
		}
		// decrypt pass
		decryptedPass, err := DecryptWithPrivateKeyPEM(cipherBytes, prikeyBytes)
		if err != nil {
			return err
		}
		// write decrypted data to output file
		_, err = output.Write(decryptedPass)
		if err != nil {
			return err
		}

		return nil
	}

	// mnmsctl util -export -adminpass -out {output_file} -pubkey {public_key_file}
	// mnmsctl util -export -configfile -out {output_file} -pubkey {public_key_file}

	if dc.Export {
		outputFile, err := getFileWithDefault(dc.Out, os.Stdout)
		if err != nil {
			return err
		}
		if outputFile != os.Stdout {
			defer outputFile.Close()
		}

		// check private key
		if dc.PrivatekeyPath != "" {
			// read private key
			prikeyBytes, err := ioutil.ReadFile(dc.PrivatekeyPath)
			if err != nil {
				return err
			}
			// set private key
			SetPrivateKey(string(prikeyBytes))

		}

		c, err := GetMNMSConfig()
		if err != nil {
			return err
		}
		if dc.PublickeyPath == "" {
			return fmt.Errorf("public key is required, specify public key with -pubkey flag")
		}
		// read public key
		pubkeyBytes, err := ioutil.ReadFile(dc.PublickeyPath)
		if err != nil {
			return err
		}

		// export mnms config
		if dc.ConfigFile {
			// export mnms config
			configJSON, err := json.Marshal(c)
			if err != nil {
				return err
			}

			// encrypt with public key
			encryptedConfig, err := EncryptWithPublicKey(configJSON, pubkeyBytes)
			if err != nil {
				return err
			}
			// write encrypted config to output file
			_, err = outputFile.Write(encryptedConfig)
			if err != nil {
				return err
			}

			return nil
		}
	}

	// mnmsctl cmd -import -configfile -in {config_file}
	if dc.Import {
		inputFile, err := getFileWithDefault(dc.In, os.Stdin)
		if err != nil {
			return err
		}
		if inputFile != os.Stdin {
			defer inputFile.Close()
		}

		if dc.ConfigFile {
			// read config file
			configEncryptedBytes, err := ioutil.ReadFile(dc.In)
			if err != nil {
				return err
			}
			configBytes, err := DecryptWithOwnPrivateKey(configEncryptedBytes)
			if err != nil {
				return err
			}

			c := MNMSConfig{}
			err = json.Unmarshal(configBytes, &c)
			if err != nil {
				return err
			}
			// save config
			err = WriteMNMSConfig(&c)
			if err != nil {
				return err
			}
			return nil

		}
		return nil

	}
	return nil
}
