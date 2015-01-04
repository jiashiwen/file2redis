package main

import (
	"./utils"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	redis "github.com/garyburd/redigo/redis"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	errlog     ErrLog
	cfg        Config
	app        *cli.App
	intypelist bool = false
	typelist        = []string{"STRINGS", "SETS", "LISTS", "HASHES", "SORTEDSETS"}
)

type Config struct {
	ServerIp     string `json:"serverip"`
	ServerPort   string `json:"serverport"`
	Separator    string `json:"separator"`
	KeyIndex     string `json:"keyindex"`
	RedisType    string `json:"redistype"`
	RedisKeyName string `json:"rediskeyname"`
	Format       string `json:"format"`
	FileName     string `json:"filename"`
	LogFile      string `json:"logfile"`
	Silent       bool   `json:"silent"`
}

type ErrLog struct {
	logfile *os.File
	prefix  string
	silent  bool
	content interface{}
}

func main() {
	initapp()
	app.Action = func(c *cli.Context) {
		doaction(c)
	}
	app.Run(os.Args)

}

func initapp() {
	app = cli.NewApp()
	app.Version = "0.0.1"
	app.Name = "file2redis"
	app.Usage = "usage of deal file to insert line to redis"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "separator, s",
			Value: " ",
			Usage: "Separator of line",
		},
		cli.StringFlag{
			Name:  "keyindex, k",
			Value: "0",
			Usage: "field key index",
		},
		cli.StringFlag{
			Name:  "redistype, t",
			Usage: "inserted redis type",
		},
		cli.StringFlag{
			Name:  "rediskeyname, n",
			Usage: "insered redis key name",
		},
		cli.StringFlag{
			Name:  "format, f",
			Usage: "fomat the key string",
		},
		cli.StringFlag{
			Name:  "serverip, i",
			Value: "127.0.0.1",
			Usage: "redis server ip address",
		},
		cli.StringFlag{
			Name:  "serverport,p",
			Value: "6379",
			Usage: "redis server port",
		},
		cli.StringFlag{
			Name:  "parafile",
			Usage: "parameter file name",
		},
		cli.StringFlag{
			Name:  "toparafile",
			Usage: "parameter file name to write",
		},
		cli.StringFlag{
			Name:  "logfile, l",
			Usage: "log file name",
		},
		cli.BoolFlag{
			Name:  "silent",
			Usage: "if slient set 'true' do not print logs to console",
		},
	}

}

func doaction(c *cli.Context) {

	if c.String("parafile") != "" {
		FileToPara(c.String("parafile"), &cfg)
	} else {
		//如果输入参数为0，显示程序help
		if len(c.Args()) != 1 {
			cli.ShowAppHelp(c)
			fmt.Println("input filename must be set!!")
			return
		}
		cfg.FileName = c.Args()[0]
		cfg.Format = c.String("format")
		cfg.KeyIndex = c.String("keyindex")
		cfg.RedisKeyName = c.String("rediskeyname")
		cfg.RedisType = c.String("redistype")
		cfg.Separator = c.String("separator")
		cfg.ServerIp = c.String("serverip")
		cfg.ServerPort = c.String("serverport")
		cfg.LogFile = c.String("logfile")
		cfg.Silent = c.Bool("silent")
	}

	if cfg.FileName == "" {
		fmt.Println("input filename must be set!!")
		return
	}
	if cfg.LogFile == "" {
		cfg.LogFile = cfg.FileName + "_" + time.Now().Format("20060102150405") + ".log"
	}

	errlog.logfile, _ = os.OpenFile(cfg.LogFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0660)
	errlog.silent = cfg.Silent

	if cfg.RedisKeyName == "" {
		cli.ShowAppHelp(c)
		errlog.content = "must get the redis key name"
		errlog.DealLog()
		return
	}

	if !IsRedisType(cfg.RedisType) {
		errlog.content = "Redis Type must be 'STRINGS, SETS, LISTS, HASHES, SORTEDSETS'"
		errlog.DealLog()
		return
	}
	if c.String("toparafile") != "" && c.String("parafile") == "" {
		ParaToFile(c.String("toparafile"), cfg)
		fmt.Println(c.String("toparafile"), "saved!")
		return
	} else if c.String("toparafile") != "" && c.String("parafile") != "" {
		cli.ShowAppHelp(c)
		fmt.Println("parafile and toparafile can not be set together!")
		return
	}

	pool := utils.GetPool(cfg.ServerIp + ":" + cfg.ServerPort)
	conn := utils.GetConnection(pool)
	// DealFileToHashes(GetIntarray(cfg.KeyIndex), conn, cfg.FileName, cfg.Separator, cfg.RedisKeyName, cfg.Format)
	DealFileToSets(conn, cfg.FileName, cfg.RedisKeyName)
	pool.Close()

}

//此函数用于处理文件，将文件每行处理为key:value并输出
func DealFileToHashes(intarr []int, conn redis.Conn, filename, separator, rediskeyname, format string) {

	var keystr string

	flushcount := 0

	fin, openerr := os.Open(filename)
	defer fin.Close()
	if openerr != nil {
		errlog.content = openerr
		errlog.DealLog()
		return
	}

	scanner := bufio.NewScanner(fin)
	for scanner.Scan() {
		key, value, verr := Linetokeyvalue(scanner.Text(), separator, intarr)
		if verr != nil {
			fmt.Println(verr)
			continue
		}
		if format != "" {
			keystr, _ = Format(key, format)
		} else {
			for i := range key {
				if i == 0 {
					keystr = keystr + key[i]
				} else {
					keystr = keystr + separator + key[i]
				}
			}
		}
		conn.Send("hset", rediskeyname, keystr, value)
		flushcount = flushcount + 1
		if flushcount >= 1000 {
			conn.Flush()
			_, receiveerr := conn.Receive()
			if receiveerr != nil {
				fmt.Println(receiveerr)
			}
			flushcount = 0
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	conn.Flush()
	_, receiveerr := conn.Receive()
	if receiveerr != nil {
		fmt.Println(receiveerr)
	}
	fin.Close()
	rep, err := conn.Do("hgetall", rediskeyname)
	str, _ := redis.Strings(rep, err)
	fmt.Println(str)
	conn.Close()
}

func DealFileToSets(conn redis.Conn, filename, rediskeyname string) {
	flushcount := 0

	fin, openerr := os.Open(filename)
	defer fin.Close()
	if openerr != nil {
		errlog.content = openerr
		errlog.DealLog()
		return
	}

	scanner := bufio.NewScanner(fin)
	for scanner.Scan() {
		conn.Send("SADD", rediskeyname, scanner.Text())
		flushcount = flushcount + 1
		if flushcount >= 1000 {
			conn.Flush()
			_, receiveerr := conn.Receive()
			if receiveerr != nil {
				fmt.Println(receiveerr)
			}
			flushcount = 0
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	conn.Flush()
	_, receiveerr := conn.Receive()
	if receiveerr != nil {
		fmt.Println(receiveerr)
	}
	fin.Close()
	conn.Close()
}

//此函数用于输入以“,”作为分隔符的一组数字字符串，返回非零int型数组
func GetIntarray(s string) (intarr []int) {
	//判断字符串元素是否为非零整数
	reg := regexp.MustCompile(`^[0-9]\d*$`)
	strarr := strings.Split(s, ",")
	for index := range strarr {
		if reg.MatchString(strarr[index]) {
			i, _ := strconv.Atoi(strarr[index])
			intarr = append(intarr, i)
		}
	}
	return
}

// 此函数用于将字符串根据分隔符转换为key[]:value
func Linetokeyvalue(text string, separator string, keyindex []int) (key []string, value string, err error) {
	strarr := strings.Split(text, separator)

	if len(strarr) < len(keyindex) {
		err = errors.New("error:too much index field!")
		return
	}
	for index := range keyindex {
		key = append(key, strings.Trim(strarr[keyindex[index]], " "))
	}
	value = text
	return
}

//该函数用于格式化字符串，用输入的字符串数组顺序替换格式数据中的“||”占位符
func Format(str []string, format string) (string, error) {

	if len(str) != strings.Count(format, "||") {
		return "", errors.New("error:input string can not fit the format string")
	}
	for i := range str {
		fmt.Println(str[i])
		format = strings.Replace(format, "||", str[i], 1)
	}
	return format, nil
}

func ParaToFile(filename string, cfg Config) {

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0660)
	if err != nil {
		fmt.Println(err)
	}
	by, marshalerr := json.MarshalIndent(cfg, "", "    ")
	if marshalerr != nil {
		fmt.Println(marshalerr)
	}

	f.Write(by)
	f.Close()
}

func FileToPara(filename string, cfg *Config) {
	var b []byte
	f, err := os.Open(filename)

	defer f.Close()
	if err != nil {
		fmt.Println(err)
	} else {
		b, _ = ioutil.ReadAll(f)
		json.Unmarshal(b, cfg)
	}
}

func (errlog *ErrLog) DealLog() {
	log.SetPrefix(errlog.prefix)
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	if !errlog.silent {
		log.SetOutput(os.Stdout)
		log.Println(errlog.content)
		log.SetOutput(errlog.logfile)
		log.Fatalln(errlog.content)
	} else {
		log.SetOutput(errlog.logfile)
		log.Fatalln(errlog.content)
	}
}

func IsRedisType(input string) bool {
	input = strings.ToUpper(input)
	result := false
	for i := range typelist {
		if input == typelist[i] {
			result = true
		}
	}
	return result
}
