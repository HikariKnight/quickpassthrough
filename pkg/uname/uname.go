package uname

import "syscall"

type Uname struct {
	Sysname    string
	Nodename   string
	Hostname   string
	Release    string
	Kernel     string
	Version    string
	Machine    string
	Arch       string
	Domainname string
}

// A utility to convert int8 values to proper strings.
func int8ToStr(arr []int8) string {
	b := make([]byte, 0, len(arr))
	for _, v := range arr {
		if v == 0x00 {
			break
		}
		b = append(b, byte(v))
	}
	return string(b)
}

func New() *Uname {
	var system syscall.Utsname
	uname := &Uname{}
	if err := syscall.Uname(&system); err == nil {
		// extract members:
		// type Utsname struct {
		//  Sysname    [65]int8
		//  Nodename   [65]int8
		//  Release    [65]int8
		//  Version    [65]int8
		//  Machine    [65]int8
		//  Domainname [65]int8
		// }

		// Add to the uname struct for humans
		uname.Sysname = int8ToStr(system.Sysname[:])
		uname.Nodename = int8ToStr(system.Nodename[:])
		uname.Hostname = uname.Nodename
		uname.Release = int8ToStr(system.Release[:])
		uname.Kernel = uname.Release
		uname.Version = int8ToStr(system.Version[:])
		uname.Machine = int8ToStr(system.Machine[:])
		uname.Arch = uname.Machine
		uname.Domainname = int8ToStr(system.Domainname[:])
	}

	return uname
}
