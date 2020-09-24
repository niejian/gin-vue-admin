package utils

import (
	"fmt"
	"gin-vue-admin/global"
	"github.com/pkg/sftp"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"log"
	"mime/multipart"
	"net"
	"os"
	"path"
	"time"
)

// 创建ssh连接，并返回当前会话
func DoSshConnect(username, password, host string, port int) (*ssh.Session, error) {
	if 0 == port {
		// 使用默认连接端口
		port = 22
	}

	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)

	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:    username,
		Auth:    auth,
		Timeout: 60 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// ssh 连接
	addr = fmt.Sprintf("%s:%d", host, port)
	client, err = ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		global.GVA_LOG.Error("ssh连接失败:", zap.Any("err", err))
		return nil, err
	}
	// 创建session
	session, err = client.NewSession()
	if nil != err {
		global.GVA_LOG.Error("创建session失败:", zap.Any("err", err))
		return nil, err
	}
	return session, nil
}

func doSftpConnect(username, password, host string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)

	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))
	clientConfig = &ssh.ClientConfig{
		User:    username,
		Auth:    auth,
		Timeout: 60 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	addr = fmt.Sprintf("%s:%d", host, port)
	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}

	return sftpClient, nil
}

// 直接文件上传
func DoFileScp(username, password, host, remoteDir, fileName string,
	port int, localFile *multipart.FileHeader, restart bool) error {

	var (
		sftpClient *sftp.Client
		err        error
	)

	// ssh建立连接
	sftpClient, err = doSftpConnect(username, password, host, port)
	if err != nil {
		log.Printf("对目标主机 %v建立连接失败\n", host)
		return err
	}

	defer sftpClient.Close()

	// 远程创建文件夹

	session, err := DoSshConnect(username, password, host, port)
	if err != nil {
		log.Printf("ssh: %v失败 ", host)
		return err
	}

	session.Run("mkdir -p " + remoteDir)

	defer session.Close()
	sftpClient.Remove(path.Join(remoteDir, fileName))
	dest, err := sftpClient.Create(path.Join(remoteDir, fileName))
	if err != nil {
		log.Printf("远程文件夹：%v 创建失败", dest)
		return err
	}

	defer dest.Close()
	buf := make([]byte, 1024)
	mutiFile, err := localFile.Open()
	if nil != err {
		log.Printf("localFile.Open fail: %v", err)
	}
	for {
		read, _ := mutiFile.Read(buf)
		if read == 0 {
			break
		}
		dest.Write(buf[0:read])
	}
	// 修改文件权限
	dest.Chmod(os.FileMode(0755))
	defer mutiFile.Close()

	// 是否重启watchDog
	if restart {

		go func() error {
			session, err := DoSshConnect(username, password, host, port)
			if err != nil {
				log.Printf("ssh: %v失败 ", host)
				return err
			}

			//var startCMD = " nohup "+remoteDir+"/watch-dog 2>&1 &"
			//output, err := session.CombinedOutput("cd " + remoteDir + " && ./restart.sh watch-dog " + " && cd " + remoteDir + " && " + startCMD)
			//output, err := session.Output("cd " + remoteDir + " && ./restart.sh watch-dog " + " && cd " + remoteDir + " && " + startCMD)
			//output, err := session.Output("cd " + remoteDir + " && ./restart.sh watch-dog > /dev/null 2>&1 &")
			//stopCmd := fmt.Sprintf("%s%s%s", "/bin/bash ", remoteDir, "/stop.sh watch-dog > /dev/null 2>&1 \n")
			stopCmd := fmt.Sprintf("%s%s", remoteDir, "/stop.sh watch-dog ")
			//output, err := session.Output("cd " + remoteDir + " && ./stop.sh watch-dog > /dev/null 2>&1 &")
			output, err := session.Output(stopCmd)
			fmt.Printf("stop 命令 %v \n", stopCmd)

			global.GVA_LOG.Info("重启命令返回结果：", zap.Any("重启命令返回结果", string(output)))
			defer session.Close()
			return nil
		}()

	}

	return nil
}

// 单个文件复制路径
func DoPathscp(username, password, host,
	localFilePath, remoteDir string, port int) error {
	var (
		sftpClient *sftp.Client
		err        error
	)

	// ssh建立连接
	sftpClient, err = doSftpConnect(username, password, host, int(port))
	if err != nil {
		log.Printf("对目标主机 %v建立连接失败\n", host)
		return err
	}

	defer sftpClient.Close()
	src, err := os.Open(localFilePath)
	if err != nil {
		log.Printf("打开本地文件: %v失败 ", localFilePath)
		return err
	}
	defer src.Close()
	// 远程创建文件夹

	session, err := DoSshConnect(username, password, host, int(port))
	if err != nil {
		log.Printf("ssh: %v失败 ", localFilePath)
		return err
	}

	session.Run("mkdir -p " + remoteDir)

	defer session.Close()

	remoteFileName := path.Base(localFilePath)
	dest, err := sftpClient.Create(path.Join(remoteDir, remoteFileName))
	if err != nil {
		log.Printf("远程文件夹：%v 创建失败", dest)
		return err
	}

	defer dest.Close()
	buf := make([]byte, 1024)
	for {
		read, _ := src.Read(buf)
		if read == 0 {
			break
		}
		dest.Write(buf[0:read])
	}
	dest.Chmod(os.FileMode(0755))
	return nil
}

// 获取远程环境的指定目录下的文件信息
func GetRemoteFiles(username, password, host, remoteDir string, port int) (bool, error) {
	session, err := DoSshConnect(username, password, host, port)
	if err != nil {
		log.Fatalf("ssh连接失败，%v", err)
		return false, err
	}
	defer session.Close()
	output, err := session.CombinedOutput("ls -l " + remoteDir)
	if err != nil {
		log.Fatalf("查看文件列表失败，%v", err)
		return false, err
	}
	log.Printf("查看文件信息结果：%v", string(output))
	return true, nil

}
