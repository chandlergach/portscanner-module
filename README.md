# Portscanner-Module
A lightweight port scanner module written in Go.

# Setup
domains.txt is direct output from the tool Amass.
> amass enum -ipv4 -d XXX.com > domains.txt


> go run portscanner.go

# Settings
Currently scans top 20 most common TCP ports:


  
	21/22:   "ftp/ssh",
	23:   "telnet",
	25:   "smtp",
	53:   "domain",
	80:   "http", 
	110:  "pop3", 
	111:  "rpcbind", 
	135:  "msrpc",
	139:  "netbios-ssn",
	143:  "imap",
	443:  "https",
	445:  "microsoft-ds",
	993:  "imaps",
	995:  "pop3s",
	1723: "pptp",
	3306: "mysql",
	3389: "ms-wbt-server",
	5900: "vnc",
	8080: "http-proxy",
