package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode/utf16"

	"github.com/marcosdy/windowsjobs/jobobject"
	"golang.org/x/sys/windows"
)

const (
	_READ_CONTROL         uint32 = 0x00020000
	JOB_OBJECT_ALL_ACCESS        = 0x1F001F
)

var (
	cIDFlag = flag.String("cid", "", "container ID")
	pIDFlag = flag.String("pid", "", "process ID")
)

func main() {
	flag.Parse()

	pID, err := strconv.ParseInt(*pIDFlag, 10, 0)
	if err != nil {
		log.Fatalf("failed to parse process ID:%v", err)
	}

	pHandle, err := windows.OpenProcess(JOB_OBJECT_ALL_ACCESS, false, uint32(pID))
	if err != nil {
		log.Fatalf("failed to open process: %v", err)
	}
	defer windows.Close(pHandle)

	log.Printf("process handle: %v\n", pHandle)

	// cID := `\Container_` + *cIDFlag
	cID := strings.TrimSpace(*cIDFlag)
	buf, err := windows.UTF16FromString(cID)
	if err != nil {
		log.Fatalf("failed to parse to UTF 16: %v", err)
	}
	fmt.Printf("buf size: %v\n", len(buf))
	fmt.Printf("%v\n", buf)
	fmt.Println(cID)
	ints := utf16.Decode(buf)
	fmt.Println(string(ints))

	jobHandle, err := jobobject.OpenJobObjectW(_READ_CONTROL, false, &buf[0])
	if err != nil {
		log.Fatalf("failed to open job: %v", err)
	}
	defer windows.Close(jobHandle)
	log.Printf("job handle: %v\n", jobHandle)

	var result bool
	if err := jobobject.IsProcessInJob(pHandle, jobHandle, &result); err != nil {
		log.Printf("failed to verify if process is in job: %v\n", jobHandle)
	}

	log.Printf("process is in job: %v", result)
}
