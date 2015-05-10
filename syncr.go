package syncr

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
)

type Syncr struct {
	// Protects all public fields
	Lock sync.Mutex

	Description string
	Status      [2]string
	cmd         *exec.Cmd
	Dead        bool
	Error       error
}

func NewSyncr(src, dst string) (*Syncr, error) {
	if strings.HasPrefix(src, "~/") {
		usr, err := user.Current()
		if err != nil {
			return nil, err
		}
		src = fmt.Sprintf("%s%c%s", usr.HomeDir, filepath.Separator, src[2:])
	}

	var syncr Syncr

	syncr.Description = fmt.Sprintf("%s -> %s", src, dst)
	syncr.cmd = exec.Command("rsync", "-HP", "-vax", "--delete", src, dst)
	pr, err := syncr.cmd.StdoutPipe()
	syncr.cmd.Stderr = syncr.cmd.Stdout
	if err != nil {
		syncr.Dead = true
		syncr.Error = err
		return &syncr, err
	}

	if err = syncr.cmd.Start(); err != nil {
		syncr.Dead = true
		syncr.Error = err
		return &syncr, err
	}

	go syncr.watch(pr)

	return &syncr, nil
}

func (s *Syncr) watch(r io.Reader) {
	buf := bufio.NewReader(r)

	for {
		c, err := buf.ReadByte()
		s.Lock.Lock()
		if err == io.EOF {
			s.Dead = true
			s.Lock.Unlock()
			break
		}
		if c == '\r' {
			s.Status[1] = ""
		} else if c == '\n' {
			if len(s.Status[1]) > 0 {
				s.Status[0] = s.Status[1]
				s.Status[1] = ""
			}
		} else {
			if len(s.Status[1]) > 79 {
				s.Status[1] = s.Status[1][0:30] + "..." + s.Status[1][35:]
			}
			s.Status[1] += string(c)
		}
		s.Lock.Unlock()
	}

	err := s.cmd.Wait()
	if err != nil && s.Error == nil {
		s.Lock.Lock()
		s.Error = err
		s.Lock.Unlock()
	}
}
