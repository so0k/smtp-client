package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"

	"github.com/urfave/cli"
)

type (
	// Mail holds structured email data
	Mail struct {
		senderID string
		toIds    []string
		subject  string
		body     string
	}

	// SMTPServer holds host and port
	SMTPServer struct {
		Host string `json:"host"`
		Port string `json:"port"`
	}

	// Config holds server address and authentication details
	Config struct {
		AuthUser     string     `json:"auth_user"`
		AuthPassword string     `json:"auth_password"`
		SMTPServer   SMTPServer `json:"smtp_server"`
	}
)

// ServerName returns host port concatenated string
func (s *SMTPServer) ServerName() string {
	return s.Host + ":" + s.Port
}

// BuildMessage returns SMTP Message Body
func (mail *Mail) BuildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", mail.senderID)
	if len(mail.toIds) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.toIds, ";"))
	}

	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
	message += "\r\n" + mail.body

	return message
}

func main() {
	app := cli.NewApp()
	app.Name = "smtp-client"
	app.Version = "0.1"
	app.Usage = "Send test messages through SMTP"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "host,H",
			Usage:  "`HOSTNAME` for the SMTP server",
			EnvVar: "SMTP_HOST",
		},
		cli.StringFlag{
			Name:   "port,p",
			Usage:  "`PORT` for the SMTP server",
			EnvVar: "SMTP_PORT",
			Value:  "587",
		},
		cli.StringFlag{
			Name:   "user,u",
			Usage:  "`USERNAME` for authentication",
			EnvVar: "SMTP_USERNAME",
		},
		cli.StringFlag{
			Name:   "password",
			Usage:  "`PASSWORD` for authentication",
			EnvVar: "SMTP_PASSWORD",
		},
		cli.BoolFlag{
			Name:  "ssl",
			Usage: "use SSL/TLS (default: STARTTLS)",
		},
		cli.StringSliceFlag{
			Name:  "recipients,t",
			Usage: "`LIST` of recipients",
		},
		cli.StringFlag{
			Name:  "sender,f",
			Usage: "`EMAIL` address of sender",
		},
		cli.StringFlag{
			Name:  "subject,s",
			Usage: "`SUBJECT` for email message",
		},
		cli.StringFlag{
			Name:  "body,b",
			Usage: "`BODY` for email message",
		},
	}

	app.Action = run

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initConfig(c *cli.Context) (*Config, *Mail, error) {
	cfg := &Config{
		AuthUser:     c.String("user"),
		AuthPassword: c.String("password"),
		SMTPServer: SMTPServer{
			Host: c.String("host"),
			Port: c.String("port"),
		},
	}

	if len(cfg.AuthUser) == 0 || len(cfg.SMTPServer.Host) == 0 || len(cfg.SMTPServer.Host) == 0 {
		return nil, nil, fmt.Errorf("SMTP Server host, port and username required")
	}

	mail := &Mail{
		toIds:    c.StringSlice("recipients"),
		senderID: c.String("sender"),
		subject:  c.String("subject"),
		body:     c.String("body"),
	}

	if len(mail.toIds) < 1 || len(mail.subject) == 0 || len(mail.body) == 0 {
		return nil, nil, fmt.Errorf("At least one recipient is required and subject/body can't be empty")
	}

	if len(mail.senderID) == 0 {
		mail.senderID = cfg.AuthUser
	}

	return cfg, mail, nil
}

func run(c *cli.Context) error {
	cfg, mail, err := initConfig(c)
	if err != nil {
		cli.ShowAppHelp(c)
		return err
	}
	tlsConfig := &tls.Config{ServerName: cfg.SMTPServer.Host}

	//build an auth
	auth := smtp.PlainAuth("", cfg.AuthUser, cfg.AuthPassword, cfg.SMTPServer.Host)

	// Connect to SMTP server
	var client *smtp.Client
	if c.Bool("ssl") {
		fmt.Printf("Using SSL/TLS on port %q\n", cfg.SMTPServer.Port)
		conn, err := tls.Dial("tcp", cfg.SMTPServer.ServerName(), tlsConfig)
		if err != nil {
			return err
		}
		client, err = smtp.NewClient(conn, cfg.SMTPServer.Host)
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("Using STARTTLS on port %q\n", cfg.SMTPServer.Port)
		var err error
		client, err = smtp.Dial(cfg.SMTPServer.ServerName())
		if err != nil {
			return err
		}
		client.StartTLS(tlsConfig)
	}

	// step 1: Use Auth
	if err := client.Auth(auth); err != nil {
		return err
	}

	// step 2: add all from and to
	if err := client.Mail(mail.senderID); err != nil {
		return err
	}
	for _, k := range mail.toIds {
		if err := client.Rcpt(k); err != nil {
			return err
		}
	}

	// step 3: Send data
	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(mail.BuildMessage()))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	client.Quit()

	log.Println("Mail sent successfully")
	return nil
}
