# 🔒 sshbrute
* sshbrute is a simple ssh dictionary attack tool written in go

# 👷 installation && usage

### installation
`go get -u github.com/vilhelmbergsoe/sshbrute && go install github.com/vilhelmbergsoe/sshbrute`

### usage
`sshbrute --help`
output:
```
Usage of sshbrute:
  -h, --host string       The host you want to attack. (host:port)
  -u, --user string       The user you would like to attack
  -w, --wordlist string   The wordlist to use for dictionary attack
```
example usage: `sshbrute -h localhost:22 -u root -w wordlist.txt`

# 🥅 goals
* [ ] give it some pizzaz/color
