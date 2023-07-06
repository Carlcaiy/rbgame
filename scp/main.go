package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/povsister/scp"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v3"
)

// const cmd = "pm2 restart slot-fafafa slot-xiyouji slot-fafafa2 slot-light slot-jelly slot-south slot-graph slot-panda slot-alice slot-station"

const key = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAw8QszMDxSHwBHc4wvikFk4PjnViGyC1HgsC6UZ2lM92+W9hM
LG7Trb/smLPRcRhOTn/TI2LxBIKkbGpqpDpMWct0HTl5SpilvIjlWk+L3QKUOVaf
84cBEpK9hAL6qG66d7bX89t7rUh6LLdcKvQbZn88CHeLrON9HGZVR3TMLkfhe3Rr
bosS2fTNvpQNQHnxQsBngoqgNn9b0lXbzz1vsW8iPHouc+/D7VyMCrQad6mr/COY
olJThVr0mygBUfyi99GsWpqs62zI014Dsfgs2qV7N1F+iZRw3PoflVmw6GrWEroe
tciJXduxGRYQM35JmlpT66Du3aS2gkDyjKJEPQIDAQABAoIBAQC4rNeWUG1Nc62Z
8GkK+qfIASM3y8taC1zqe+VIGO8/fm+VNPam8+W8gtEvPHLXvZYhd3Q2bZ/wIU36
+GihhF2CV+uxpgZF2LqAoKO8Dk5ir4wkixNZKIJagxNE9YqAWfSN+m6+HM8PKNAY
XuETpDQ0/NIlKqEY94GOyPqp7gSvonY5mykmJxFTIEltTXWdaNyJV7M7vtmiosRu
Y9kWmtja+SJig7JBdWIlfWPp6P5uFTSA1mQjC0kPKgOKwZRMosNGh+EUcpsuwnZn
0gf4bsJ3hr9mfUKlYqJB0lsv9gLXS3vq5JaNpqlRqZv0MxVO04J7EWWPCBUoq6XD
rIccvbABAoGBAPwTjc9k/OUPA7150Ts8r+QGqiC3k1jKQBucHRQrERQdwmFy4fOT
tbYqfAYcB9yiVbWsb0OdcTI5n6OuEgaf08VYZDyMxC0dSBXxV/QiXSFDyx/T1YRa
xG16k/65aHKaxqJcmww+IvKEgqzpxv09mkZGhfwVq+/0Yb/GoGsaVB3pAoGBAMbQ
Pi8c1YCAANiPZi/qfLGZJ4Qx2PqyOYYafnJUuFoHaOveo0Y//xLzPkxDKF1GJB34
vGTyZbe2UcO424xkjHGXsizndl5Gma3wCMcreFmb2Zl6jdPgvsGCxwCErsfJ3zhB
pbEAI4GuEpmG9w0fVApsIMvGqZU8NNvFW1EUvJs1AoGAJjWYk01NgDMMcBYc/wut
5bSU6SyqaxUNLqI1Ti9RAWmZY0gUs+U58Aj0j8CD+I+qykT/AJLG68QMTLVACj0D
zrAdUYhM3EHKAXl5yYnD8Bzkl0h23v8bMzUXZc8Y4/ZOEaJT5kEs7vHjFO7CuPQz
hF+AQ7vNOiwi/PzQqCtvgLECgYABQxEh+2zfg/B9b2uWokZoWjCGBkr6TGdREpn6
387Lw4BG/wLKT12vIRwkH8kBzpAsIIgRm/hJUj7ynxnFql252tymrFF0B76x+/QS
T5RT/UlEUCLWVXbgg5P/zNfPLNjd4ozKstWG3TQBBXpL+wbtigIrSPeKmvc67eXG
ffs0gQKBgQCAczH3w/bPipXB6coRP9fNesvk5vXaqfyOK58IYpHXJW2HV7zrg1QY
zKmjPxP4nHFPgIuWSYz37LMzKqFEkjlwMYGfer+hc1cigNb+LPIBMcgExbuZmLKg
13G56yzo2hvjvREfFf68FAs69zyqxjH4EY5ujxKIrDpYokCmv0Uqmg==
-----END RSA PRIVATE KEY-----`

var (
	addr      = "47.115.43.246:3800"
	mapidname = map[string]bool{
		"fafafa":  true,
		"panda":   true,
		"dragon":  true,
		"xiyou":   true,
		"faplus":  true,
		"light":   true,
		"south":   true,
		"jelly":   true,
		"graph":   true,
		"station": true,
		"alice":   true,
		"penguin": true,
		"volcano": true,
	}
)

type Config struct {
	Download []string `yaml:"download"`
	Upload   []string `yaml:"upload"`
}

func main() {

	config := new(Config)
	file, err := ioutil.ReadFile("config.yaml")
	if err == nil {
		err = yaml.Unmarshal(file, config)
		if err != nil {
			return
		}
	}
	// Build a SSH config from username/password
	// sshConf := scp.NewSSHConfigFromPassword("root", "password")

	// Build a SSH config from private key
	// privPEM, err := ioutil.ReadFile("/path/to/privateKey")
	// without passphrase
	sshConf, _ := scp.NewSSHConfigFromPrivateKey("root", []byte(key))
	// with passphrase
	// sshConf, err := scp.NewSSHConfigFromPrivateKey("username", privPEM, passphrase)

	// Dial SSH to "my.server.com:22".
	// If your SSH server does not listen on 22, simply suffix the address with port.
	// e.g: "my.server.com:1234"
	scpClient, err := scp.NewClient(addr, sshConf, &scp.ClientOption{})
	checkErr(err)

	// Build a SCP client based on existing "golang.org/x/crypto/ssh.Client"
	// scpClient, err := scp.NewClientFromExistingSSH(existingSSHClient, &scp.ClientOption{})

	defer scpClient.Close()

	// Do the file copy with timeout, context and file properties preserved.
	// Note that the context and timeout will both take effect.
	fo := &scp.FileTransferOption{
		Context:      context.Background(),
		Timeout:      30 * time.Second,
		PreserveProp: true,
	}

	for _, file := range config.Upload {
		if mapidname[file] {
			if _, err := os.Stat(file + ".yaml"); err == nil {
				err = scpClient.CopyFileToRemote(file+".yaml", fmt.Sprintf("/data/gameserver/gopb/%s/slots.yaml", file), fo)
				checkErr(err)
				fmt.Printf("Upload %s Success\n", file)
				time.Sleep(time.Second)
			}
		}
	}

	for _, file := range config.Download {
		if mapidname[file] {
			if _, err := os.Stat(file + ".yaml"); os.IsNotExist(err) {
				err = scpClient.CopyFileFromRemote(fmt.Sprintf("/data/gameserver/gopb/%s/slots.yaml", file), file+".yaml", fo)
				checkErr(err)
				fmt.Printf("Download %s Success\n", file)
				time.Sleep(time.Second)
			}
		}
	}

	if len(config.Download) == 0 && len(config.Upload) == 0 {
		for file := range mapidname {
			if _, err := os.Stat(file + ".yaml"); err == nil {
				err = scpClient.CopyFileToRemote(file+".yaml", fmt.Sprintf("/data/gameserver/gopb/%s/slots.yaml", file), fo)
				checkErr(err)
				fmt.Printf("Upload %s Success\n", file)
				time.Sleep(time.Second)
				return
			}
		}
		file := "panda"
		err = scpClient.CopyFileFromRemote(fmt.Sprintf("/data/gameserver/gopb/%s/slots.yaml", file), file+".yaml", fo)
		checkErr(err)
		fmt.Printf("Download %s Success\n", file)
		time.Sleep(time.Second)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func publicKeyAuthFunc(kPath string) ssh.AuthMethod {
	// fmt.Println(kPath)
	// keyPath, err := homedir.Expand(kPath)
	// fmt.Println(keyPath)
	// if err != nil {
	// 	log.Fatal("find key's home dir failed", err)
	// }
	// key, err := ioutil.ReadFile("id_rsa")
	// if err != nil {
	// 	log.Fatal("ssh key file read failed", err)
	// }
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey([]byte(key))
	if err != nil {
		log.Fatal("ssh key signer failed", err)
	}
	return ssh.PublicKeys(signer)
}
