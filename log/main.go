package main

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/ssh"
)

// const cmd = "pm2 restart slot-fafafa slot-xiyouji slot-fafafa2 slot-light slot-jelly slot-south slot-graph slot-panda slot-alice slot-station"
var cmd = "tail -n 300 ~/.pm2/logs/pb-panda-out.log"

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

var mapidcommand = map[int]string{
	1:  "tail -n 300 ~/.pm2/logs/pb-fafafa-out.log",
	2:  "tail -n 300 ~/.pm2/logs/pb-panda-out.log",
	3:  "tail -n 300 ~/.pm2/logs/pb-dragon-out.log",
	4:  "tail -n 300 ~/.pm2/logs/pb-xiyou-out.log",
	5:  "tail -n 300 ~/.pm2/logs/pb-faplus-out.log",
	6:  "tail -n 300 ~/.pm2/logs/pb-light-out.log",
	7:  "tail -n 300 ~/.pm2/logs/pb-south-out.log",
	8:  "tail -n 300 ~/.pm2/logs/pb-jelly-out.log",
	9:  "tail -n 300 ~/.pm2/logs/pb-graph-out.log",
	10: "tail -n 300 ~/.pm2/logs/pb-station-out.log",
	11: "tail -n 300 ~/.pm2/logs/pb-alice-out.log",
	12: "tail -n 300 ~/.pm2/logs/pb-penguin-out.log",
	13: "tail -n 300 ~/.pm2/logs/pb-volcano-out.log",
}

func main() {
	round()
}

func round() {
	for {
		fmt.Println("选择游戏：")
		fmt.Println("0：退出 1：发发发 2：熊猫 3：金龙 4：西游 5：发发发2 6：光浪 7：东南亚 8：水果 9：几何 10：车站 11：爱丽丝 12：企鹅 13：火山")
		var i int = 1
		n, err := fmt.Scanln(&i)
		if err != nil && n != 1 {
			continue
		}
		if i == 0 {
			break
		}
		_, ok := mapidcommand[i]
		if ok {
			for {
				fmt.Println("刷新日志?y/n")
				ss := "y"
				n, err = fmt.Scanln(&ss)
				if err == nil && n == 1 && ss == "n" {
					break
				}
				cmd = mapidcommand[i]
				do()
			}
		}
	}
}

func do() {
	sshHost := "47.115.43.246"
	sshUser := "root"
	sshPassword := ""
	sshType := "key"                                               //password 或者 key
	sshKeyPath := "D:\\servercode\\go\\ckgame\\restart\\unity.ppk" //ssh id_rsa.id 路径"
	sshPort := 3800

	//创建sshp登陆配置
	config := &ssh.ClientConfig{
		Timeout:         time.Second, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            sshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
		// HostKeyCallback: hostKeyCallBackFunc(h.Host),
	}
	if sshType == "password" {
		config.Auth = []ssh.AuthMethod{ssh.Password(sshPassword)}
	} else {
		config.Auth = []ssh.AuthMethod{publicKeyAuthFunc(sshKeyPath)}
	}

	//dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", sshHost, sshPort)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Fatal("创建ssh client 失败", err)
	}
	defer sshClient.Close()

	//创建ssh-session
	session, err := sshClient.NewSession()
	if err != nil {
		log.Fatal("创建ssh session 失败", err)
	}
	defer session.Close()
	//执行远程命令
	combo, err := session.CombinedOutput(cmd)
	if err != nil {
		log.Fatal("远程执行cmd 失败", err)
	}
	log.Println("命令输出:", string(combo))
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
