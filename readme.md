JumpZoneLookup
---
A bruteforce reverse DNS lookup tool for active directory red teaming. Drop the executable onto the victim/pivot machine, and run with a cidr notation. If PTR records are avaible, will resolve the to domain names. I needed something quick and portable for an engagement so I wrote this.

Can brute force a /16 with in a few mins depending on the workstation (i5, 4gb ram) and DNS server specs.

Uses the default resolver of the host. NOTE: On Windows, no matter what, using the `net.lookupaddr` with use the host default resolver. No getting around that unless you use a non standard library. Because of that, I won't be implementing a custom resolver feature as this was most made for windows.

---

Build
---
```
go build -o JumpZoneLookup.exe -ldflags '-s' main.go 
```
---

## Usage
```
JumpZoneLookup --t 192.168.0.0/16
```

Output
```
192.168.1.20 - pi.demo.local.
192.168.1.51 - camera-system.demo.local.
192.168.1.50 - homeassistant.demo.local.
192.168.1.52 - jellyfin.demo.local.
192.168.1.111 - Jason.demo.local.
192.168.1.128 - terrys-Air-5.demo.local.
192.168.1.159 - DESKTOP-4JIB526.demo.local.
192.168.1.152 - LGwebOSTV.demo.local.
192.168.1.177 - LAPTOP-demo.demo.local.
192.168.1.175 - JasonMacbookPro.demo.local.
192.168.1.168 - user-ThinkPad-X270.demo.local.
192.168.1.200 - DESKTOP-A8VETIH.demo.local.
```

---
### Warning
No rate limiting

Pretty Noisy

Uses the default resolver(s) of the victem/proxy host.

----

*Part of the Jump Collection*

