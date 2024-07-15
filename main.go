package main

import (
	"community-robot-lib/config"
	"community-robot-lib/interrupts"
	liboptions "community-robot-lib/options"
	"community-robot-lib/secret"
	"flag"
	"fmt"
	sdk "git-platform-sdk"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/sys/windows"

	"github.com/sirupsen/logrus"

	_ "github.com/go-playground/validator/v10"
)

type options struct {
	service liboptions.ServiceOptions
	client  liboptions.ClientOptions
}

func (o *options) Validate() error {

	if err := o.service.Validate(); err != nil {
		return err
	}

	return o.client.Validate()
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options

	o.client.AddFlags(fs)
	o.service.AddFlags(fs)

	_ = fs.Parse(args)
	return o
}

func main() {
	//opt := gatherOptions(flag.NewFlagSet(os.Args[0], flag.ExitOnError), os.Args[1:]...)
	opt := options{
		service: liboptions.ServiceOptions{
			Port:         7102,
			ConfigFile:   "D:\\Project\\github\\ibfru\\robot-platform-cache\\local\\config.yaml",
			GracePeriod:  300 * time.Second,
			ReadTimeout:  120 * time.Second,
			WriteTimeout: 120 * time.Second,
			IdleTimeout:  30 * time.Minute,
		},
		client: liboptions.ClientOptions{
			TokenPath:   "D:\\Project\\github\\ibfru\\robot-platform-cache\\local\\token",
			HandlerPath: "/dd",
		},
	}
	if err := opt.Validate(); err != nil {
		logrus.WithError(err).Fatal("Configuration invalid: " + err.Error())
		return
	}

	secretAgent := new(secret.Agent)
	if err := secretAgent.Start([]string{opt.client.TokenPath}); err != nil {
		logrus.WithError(err).Fatal("Error starting secret agent.")
	}

	defer secretAgent.Stop()

	bot := newRobot(sdk.GetClientInstance(secretAgent.GetSecret(opt.client.TokenPath)))

	agent := config.NewConfigAgent(bot.NewConfig)
	if err := agent.Start(opt.service.ConfigFile); err != nil {
		logrus.WithError(err).Errorf("start config:%s", opt.service.ConfigFile)
		return
	}

	defer interrupts.WaitForGracefulShutdown()

	interrupts.OnInterrupt(func() {
		agent.Stop()
	})

	httpServer := &http.Server{
		Addr:         ":" + strconv.Itoa(opt.service.Port),
		Handler:      bot.setupRouter(),
		ReadTimeout:  opt.service.ReadTimeout,
		WriteTimeout: opt.service.WriteTimeout,
		IdleTimeout:  opt.service.IdleTimeout,
	}

	fmt.Printf("\u001B[0;31;6m %s%d \u001B[0;30;6m \n", "=========================",
		windows.GetCurrentProcessId())
	interrupts.ListenAndServe(httpServer, opt.service.GracePeriod)

}
