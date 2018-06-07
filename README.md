## SMTP Client

Supports SSL/TLS and STARTTLS with plain auth only

```bash
NAME:
   smtp-client - Send test messages through SMTP

USAGE:
   smtp-client [global options] command [command options] [arguments...]

VERSION:
   0.1

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --host HOSTNAME, -H HOSTNAME   HOSTNAME for the SMTP server [$SMTP_HOST]
   --port PORT, -p PORT           PORT for the SMTP server (default: "587") [$SMTP_PORT]
   --user USERNAME, -u USERNAME   USERNAME for authentication [$SMTP_USERNAME]
   --password PASSWORD            PASSWORD for authentication [$SMTP_PASSWORD]
   --ssl                          use SSL/TLS (default: STARTTLS)
   --recipients LIST, -t LIST     LIST of recipients
   --sender EMAIL, -f EMAIL       EMAIL address of sender
   --subject SUBJECT, -s SUBJECT  SUBJECT for email message
   --body BODY, -b BODY           BODY for email message
   --help, -h                     show help
   --version, -v                  print the version
```
