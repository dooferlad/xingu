# xingu
An alternative AWS command line tool

# Why?
## The name
[A tributory of the Amazon](https://en.wikipedia.org/wiki/Xingu_River). It is written in Go so I could
have gone with xingo, and I do love a pun, but it turns out that Xingo is a character from Ben 10.

## The reason
Initially I wanted to download RDS logs and the way to do this is supposed to be
```bash
$ aws rds download-db-log-file-portion --db-instance-identifier <db ID> --log-file-name <filename> --output text --starting-token 0 > /tmp/psql.log
```
...but this didn't work. It would have bits of file missing. The reliable way seems to be using the downloadCompleteLogFile API endpoint. This tool does that.

# Usage
## Configuration
Primarily xingu reads your AWS configuration from your environment and your home directory, just like the `aws` tool does.
In addition it has its own config file at $HOME/.xingu.yaml that contains some optional configuration:

```yaml
# Xingu can act as an OTP token provider so you don't have to enter the
# token every time. It uses the same configuration as
# https://github.com/pcarrier/gauth
otp: <otp token>

# Set per-configuration SSH config. This is useful if you have lots of
# SSH keys and don't want to offer all of them to the remote server, since
# most will give up after a few and tell you that your authentication
# has failed. It is also useful because you can have individually
# configured jump hosts that match on 10.*.*.*. If you have many VPCs
# in the 10.*.*.* range with overlapping subnets, this allows you to pick
# the right jump host for the right VPC.
prod:
  ssh:
    config: /home/dooferlad/.ssh/prod-config
```

## Commands
```bash
$ xingu rds logs
Interact with RDS logs

Usage:
  xingu rds logs [command]

Available Commands:
  download    Download an RDS log file
  list        List RDS logs
  list        List RDS logs

Flags:
  -h, --help   help for logs

Global Flags:
      --config string   config file (default is $HOME/.xingu.yaml)

Use "xingu rds logs [command] --help" for more information about a command.
```

### List log files
```bash
$ xingu rds logs list -d <database ID>
```

### Download a single log file
```
$ xingu rds logs download -d <database ID> -f <filename>
```

### Download several days of logs
```bash
xingu rds logs download -d <database ID> --days 2
```

### List ec2 instances
```bash
xingu ec2 list

# or just one
xingu ec2 list -n <name filter>

# is equivalent to
xingu ec2 list --filters "tag:Name=<name filter>"
````

### SSH into an ec2 instance
Takes the same filters as list...

```bash
xingu ec2 ssh # the first one in the list

# you probably want to be more specific
xingu ec2 ssh -n <name filter>
````

# TODO
 - [ ] Download files if missing or smaller than on server
 - [ ] Gzip files before writing to disk
