package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/bndr/gojenkins"
	"net/http"
	"strings"
)

type Jh struct {
	Url     string
	Nop     bool
	jenkins *gojenkins.Jenkins
}

func NewJh(url string, user string, pwd string, nop bool) (ret *Jh, err error) {
	var jenkins *gojenkins.Jenkins
	client := http.DefaultClient
	if jenkins, err = gojenkins.CreateJenkins(client, url, user, pwd).Init(); err == nil {
		ret = &Jh{Url: url, Nop: nop, jenkins: jenkins}
	}
	return
}

func (o *Jh) Nodes() (err error) {
	var items []*gojenkins.Node
	if items, err = o.jenkins.GetAllNodes(); err != nil {
		return
	}

	logrus.Infof("%v nodes in %v", len(items), o.Url)
	for _, item := range items {
		// Fetch Node Data
		item.Poll()
		logrus.Infof("  %v", item.GetName())
	}
	return
}

func (o *Jh) Jobs() (err error) {
	var items []gojenkins.InnerJob
	if items, err = o.jenkins.GetAllJobNames(); err != nil {
		return
	}

	logrus.Infof("%v jobs in %v", len(items), o.Url)
	for _, item := range items {
		logrus.Infof("  %v", item.Name)
	}
	return
}

func (o *Jh) DeleteNodesByPrefix(prefix string) (err error) {
	var items []*gojenkins.Node
	if items, err = o.jenkins.GetAllNodes(); err != nil {
		return
	}

	logrus.Infof("%v nodes in %v", len(items), o.Url)
	for _, item := range items {
		if strings.HasPrefix(item.GetName(), prefix) {
			if o.Nop {
				logrus.Infof("  %v: delete simulated, nop=true", item.GetName())
			} else {
				ok, actionErr := item.Delete()
				if ok {
					logrus.Warnf("  %v: delete successfully", item.GetName())
				} else {
					logrus.Warnf("  %v: delete failed because of %v", item.GetName(), actionErr)
				}
			}
		} else {
			logrus.Debugf("  %v: delete node skipped, because the name does not match prefix", item.GetName())
		}
	}
	return
}
