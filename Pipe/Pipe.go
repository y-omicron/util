package Pipe

import (
	"fmt"
	"github.com/y-omicron/util/DataList"
	"github.com/y-omicron/util/Util"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

func New(SrcAddr, DestAddr net.Addr) *RunInfo {
	return &RunInfo{
		Start:    time.Now(),
		Recv:     0,
		Send:     0,
		SrcAddr:  SrcAddr,
		DestAddr: DestAddr,
		TimeCost: 0,
		Status:   "",
	}
}

var bLi = DataList.New()

type RunInfo struct {
	Start    time.Time
	Recv     int64
	Send     int64
	SrcAddr  net.Addr
	DestAddr net.Addr
	TimeCost time.Duration
	Status   string
}

func (pi RunInfo) String() string {
	var buf strings.Builder
	// [Status] [ XXX.XXX.XXX.XXX:XXXIX - XXX ] Read: [ XXX  B ] Write: [ XXX  B ] Last Time: XXXIX.XX.XX XX:XX:XX [Time: XXX  s ]
	if pi.Status == "Running" {
		buf.WriteString(fmt.Sprintf("[ \x1B[38;2;186;207;101m%-7s\x1B[0m ] [ \x1B[1m%21s - %-21s\x1B[0m ]", pi.Status, pi.SrcAddr.String(), pi.DestAddr.String()))
	} else {
		buf.WriteString(fmt.Sprintf("[ \x1B[38;2;241;151;144m%-7s\x1B[0m ] [ \x1B[38;2;105;105;105m%21s - %-21s\x1B[0m ]", pi.Status, pi.SrcAddr.String(), pi.DestAddr.String()))
	}

	buf.WriteString(fmt.Sprintf(" Send: [ %s ] Recv: [ %s ] ", Util.FormatSize(pi.Send), Util.FormatSize(pi.Recv)))
	buf.WriteString(fmt.Sprintf("(Last Time: %s [Time: %.2f s])", pi.Start.Format("2006-01-02 15:04:05"), pi.TimeCost.Seconds()))
	return buf.String()
}

type InfoData struct {
	Index *DataList.DataNode
	Data  *RunInfo
}

func Pipe(src net.Conn, dest net.Conn, pInfo InfoData, timeOut time.Duration) {
	var err error
	var n int64
	var wg = sync.WaitGroup{}
	bLi.Add(pInfo.Index)
	pInfo.Data.Start = time.Now()
	wg.Add(1)
	pInfo.Data.Status = "Running"
	pInfo.Index.Set(pInfo.Data.String())
	bLi.Update()
	go func() {
		defer wg.Done()
		var n, err = io.Copy(dest, src)
		if err != nil {
			return
		}
		err = dest.SetReadDeadline(time.Now().Add(timeOut))
		if err != nil {
			return
		}
		pInfo.Data.Send += n
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		n, err = io.Copy(src, dest)
		if err != nil {
			return
		}
		err = src.SetReadDeadline(time.Now().Add(timeOut))
		if err != nil {
			return
		}
		pInfo.Data.Recv += n
	}()
	wg.Wait()
	pInfo.Data.Status = "Closed"
	pInfo.Data.TimeCost = time.Now().Sub(pInfo.Data.Start)
	pInfo.Index.Set(pInfo.Data.String())
	bLi.Update()
}
