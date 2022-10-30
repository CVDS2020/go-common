package svc

import (
	"gitee.com/sy_183/common/uns"
	"net"
	"os"
	"syscall"
)

//func SystemdPidNotifyWithFds(pid int, unsetEnvironment bool, state string, fds []int) error {
//	address := os.Getenv("NOTIFY_SOCKET")
//	if address == "" {
//		return nil
//	}
//
//	defer func() {
//		if unsetEnvironment {
//			assert.MustSuccess(os.Unsetenv("NOTIFY_SOCKET"))
//		}
//	}()
//
//	if state == "" {
//		return syscall.EINVAL
//	}
//	addr := assert.Must(net.ResolveUnixAddr("unixgram", address))
//	conn, err := net.DialUnix("unixgram", nil, addr)
//	if err != nil {
//		return err
//	}
//
//	curPid := os.Getpid()
//	uid := os.Getuid()
//	gid := os.Getgid()
//	sendUcred := (pid != 0 && pid != curPid) || uid != os.Geteuid() || gid != os.Getegid()
//
//	if len(fds) > 0 || sendUcred {
//		var fdsCmsgSpace, ucredCmsgSpace int
//		if len(fds) > 0 {
//			fdsCmsgSpace = syscall.CmsgSpace(int(unsafe.Sizeof(int32(0))) * len(fds))
//		}
//		if sendUcred {
//			ucredCmsgSpace = syscall.CmsgSpace(int(unsafe.Sizeof(syscall.Ucred{})))
//		}
//		msgControl := make([]byte, fdsCmsgSpace+ucredCmsgSpace)
//
//		p := msgControl
//		cmsg := (*syscall.Cmsghdr)(unsafe.Pointer(&p[0]))
//		if len(fds) > 0 {
//			cmsg.Level = syscall.SOL_SOCKET
//			cmsg.Type = syscall.SCM_RIGHTS
//			cmsg.Len = uint64(syscall.CmsgLen(int(unsafe.Sizeof(int32(0))) * len(fds)))
//			p = p[unsafe.Sizeof(syscall.Cmsghdr{}):]
//			for _, fd := range fds {
//				binary.LittleEndian.PutUint32(p, uint32(fd))
//				p = p[unsafe.Sizeof(uint32(0)):]
//			}
//		}
//		if sendUcred {
//			p = msgControl[fdsCmsgSpace:]
//			cmsg = (*syscall.Cmsghdr)(unsafe.Pointer(&p[0]))
//			cmsg.Level = syscall.SOL_SOCKET
//			cmsg.Type = syscall.SCM_CREDENTIALS
//			cmsg.Len = uint64(syscall.CmsgLen(int(unsafe.Sizeof(syscall.Ucred{}))))
//			p = p[unsafe.Sizeof(syscall.Cmsghdr{}):]
//			ucred := (*syscall.Ucred)(unsafe.Pointer(&p[0]))
//			ucred.Pid = int32(pid)
//			if pid == 0 {
//				ucred.Pid = int32(curPid)
//			}
//			ucred.Uid = uint32(uid)
//			ucred.Gid = uint32(gid)
//		}
//
//		_, _, err := conn.WriteMsgUnix(uns.StringToBytes(state), msgControl, addr)
//		if err != nil {
//			if sendUcred {
//				if len(fds) > 0 {
//					_, _, err := conn.WriteMsgUnix(uns.StringToBytes(state), msgControl[:fdsCmsgSpace], addr)
//					return err
//				}
//				_, _, err = conn.WriteMsgUnix(uns.StringToBytes(state), nil, addr)
//				return err
//			}
//			return err
//		}
//	}
//
//	_, _, err = conn.WriteMsgUnix(uns.StringToBytes(state), nil, addr)
//	return err
//}
//
//func SystemdPidNotify(pid int, unsetEnvironment bool, state string) error {
//	return SystemdPidNotifyWithFds(pid, unsetEnvironment, state, nil)
//}

// SystemdNotify function notify service manager about start-up completion and
// other service status changes
func SystemdNotify(state string) error {
	address := os.Getenv("NOTIFY_SOCKET")
	if address == "" {
		return nil
	}

	if state == "" {
		return syscall.EINVAL
	}
	addr, err := net.ResolveUnixAddr("unixgram", address)
	if err != nil {
		return err
	}
	conn, err := net.DialUnix("unixgram", nil, addr)
	if err != nil {
		return err
	}
	_, err = conn.Write(uns.StringToBytes(state))
	return err
}
