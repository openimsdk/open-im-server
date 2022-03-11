package main

import (
	"bufio"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

// docker rm -f registry-srv  && rm -rf /root/docker/dockerRegistry && mkdir -p dockerRegistry && docker run -d -p 80:5000 --restart=always --name registry-srv -v /root/docker/dockerRegistry:/var/lib/registry registry
func main() {
	var paths = map[string]string{}
	pwd, _ := os.Getwd()
	readDir, _ := os.ReadDir(".")
	for _, dir := range readDir {
		if dir.IsDir() {
			if strings.HasPrefix(dir.Name(), "open_im") {
				paths[path.Join(pwd, dir.Name())] = ""
			}
		}
	}
	for dir := range paths {
		//execCommand("docker", "rmi",  "$(docker images -a | awk '/<none>/{print $3}')")
		/*go */func(dir string) {
			fmt.Printf("\n----------------------------------------%s构建开始----------------------------------------\n", dir)
			execCommand("bash", "-c", fmt.Sprintf(""+
				"cd %s"+
				" && "+
				"bash build.sh", dir))
			fmt.Printf("\n----------------------------------------%s构建完成----------------------------------------\n", dir)
			delete(paths, dir)
		}(dir)
	}
	for {
		if len(paths) == 0 {
			fmt.Println("----------------------------------------全部构建成功----------------------------------------")
			os.Exit(0)
		}
	}
}

//封装一个函数来执行命令
func execCommand(commandName string, params ...string) bool {

	//执行命令
	cmd := exec.Command(commandName, params...)

	//显示运行的命令
	fmt.Println(cmd.Args)

	stdout, err := cmd.StdoutPipe()
	errReader, errr := cmd.StderrPipe()

	if errr != nil {
		fmt.Println("err:" + errr.Error())
	}

	//开启错误处理
	go handlerErr(errReader)

	if err != nil {
		fmt.Println(err)
		return false
	}

	cmd.Start()
	in := bufio.NewScanner(stdout)
	for in.Scan() {
		cmdRe := ConvertByte2String(in.Bytes(), "UTF-8")
		fmt.Println(cmdRe)
	}

	cmd.Wait()
	cmd.Wait()
	return true
}

//开启一个协程来错误
func handlerErr(errReader io.ReadCloser) {
	in := bufio.NewScanner(errReader)
	for in.Scan() {
		cmdRe := ConvertByte2String(in.Bytes(), "UTF-8")
		fmt.Errorf(cmdRe)
	}
}

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

//对字符进行转码
func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}
