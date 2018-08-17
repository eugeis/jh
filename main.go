package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

const flagDebug = "debug"
const flagNop = "nop"
const flagUrl = "url"
const flagUser = "user"
const flagPassword = "pwd"
const flagPrefix = "prefix"

func main() {
	app := cli.NewApp()
	app.Usage = "Jenkins Helper"
	app.Version = "1.0"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  flagNop,
			Usage: "Enable non operational mode",
		}, cli.BoolFlag{
			Name:  flagDebug,
			Usage: "Enable debug log level",
		}, cli.StringFlag{
			Name:  flagUrl,
			Usage: "Url of jenkins",
		}, cli.StringFlag{
			Name:  flagUser,
			Usage: "User/ID for authentication",
		}, cli.StringFlag{
			Name:  flagPassword,
			Usage: "Password/Token for authentication",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "nodes",
			Usage: "Show nodes",
			Flags: []cli.Flag{},
			Action: func(c *cli.Context) (err error) {
				logrus.Infof("execute %v", c.Command.Name)
				var jh *Jh
				if jh, err = buildJh(c); err == nil {
					err = jh.Nodes()
				}
				logrus.Info("done")
				return
			},
		}, {
			Name:  "jobs",
			Usage: "Show jobs",
			Flags: []cli.Flag{},
			Action: func(c *cli.Context) (err error) {
				logrus.Infof("execute %v", c.Command.Name)
				var jh *Jh
				if jh, err = buildJh(c); err == nil {
					err = jh.Jobs()
				}
				logrus.Info("done")
				return
			},
		}, {
			Name:  "deleteNodes",
			Usage: "Delete nodes",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  flagPrefix,
					Usage: "Prefix for node names, e.g. docker",
				},
			},
			Action: func(c *cli.Context) (err error) {
				logrus.Infof("execute %v", c.Command.Name)
				var jh *Jh
				if jh, err = buildJh(c); err == nil {
					err = jh.DeleteNodesByPrefix(c.String(flagPrefix))
				}
				logrus.Info("done")
				return
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Warn("exit because of error.")
	}
}

func buildJh(c *cli.Context) (*Jh, error) {
	if debug := c.GlobalBool(flagDebug); debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	return NewJh(c.GlobalString(flagUrl), c.GlobalString(flagUser), c.GlobalString(flagPassword),
		c.GlobalBool(flagNop))
}
