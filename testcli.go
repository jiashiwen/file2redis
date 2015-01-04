package main

/*import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Name = "firstcli"
	app.Usage = "usage of firestcli"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "lang, l",
			Value:  "english",
			Usage:  "language for the greeting",
			EnvVar: "LEGACY_COMPAT_LANG,APP_LANG,LANG",
		},
		cli.StringFlag{
			Name:  "name, n",
			Value: "Tom",
			Usage: "name of person",
		},
	}

	app.Action = func(c *cli.Context) {
		name := "someone"
		println(c.String("name"))
		//如果输入参数为0，显示程序help
		if len(c.Args()) == 0 {
			cli.ShowAppHelp(c)
		}
		if c.String("lang") == "spanish" {
			println("Hola", name)
		} else {
			println("Hello", name)
		}
	}
	app.Commands = []cli.Command{
		{
			Name:      "add",
			ShortName: "a",
			Usage:     "add a task to the list",
			Action: func(c *cli.Context) {
				println("added task: ", c.Args().First())
			},
		},
		{
			Name:      "complete",
			ShortName: "c",
			Usage:     "complete a task on the list",
			Action: func(c *cli.Context) {
				println("completed task: ", c.Args().First())
			},
		},
		{
			Name:      "template",
			ShortName: "r",
			Usage:     "options for task templates",
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "add a new template",
					Action: func(c *cli.Context) {
						println("new task template: ", c.Args().First())
					},
				},
				{
					Name:  "remove",
					Usage: "remove an existing template",
					Action: func(c *cli.Context) {
						println("removed task template: ", c.Args().First())
					},
				},
			},
		},
	}
	app.Run(os.Args)
}
*/
import (
	"errors"
	"fmt"
	"strings"
)

func main() {
	s := []string{"fdsaf", "fewfe", "dsaf32wf"}
	str, err := Format(s, "1||2||3||4")

	fmt.Println(str, err)
	fmt.Println(append(s, "00000000"))
}

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
