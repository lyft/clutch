package k8s

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

var (
	tsDelimiter = []byte{' '}
)

const (
	limitBytes       = 1024 * 1024
	rfc3339NanoFixed = "2006-01-02T15:04:05.000000000Z07:00"
)

func (s *svc) GetPodLogs(ctx context.Context, clientset, cluster, namespace, name string, opts *k8sapiv1.PodLogsOptions) (*k8sapiv1.GetPodLogsResponse, error) {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
	if err != nil {
		return nil, err
	}

	// Verify pod's existence, since the pod logs command does not first verify the pod exists or return an error.
	_, err = cs.CoreV1().Pods(cs.Namespace()).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// Construct get logs request.
	k8sOpts, err := protoOptsToK8sOpts(opts)
	if err != nil {
		return nil, err
	}

	k8sOpts.LimitBytes = pointer.Int64(limitBytes)
	req := cs.CoreV1().Pods(cs.Namespace()).GetLogs(name, k8sOpts)
	if req == nil {
		return nil, fmt.Errorf("an unknown error occurred when constructing the GetLogs request")
	}

	// Stream logs into buffer and parse each line into the required struct for return.
	readCloser, err := req.Stream(ctx)
	if err != nil {
		return nil, err
	}
	defer readCloser.Close()

	return bufferToResponse(readCloser)
}

func bufferToResponse(r io.Reader) (*k8sapiv1.GetPodLogsResponse, error) {
	var logs []*k8sapiv1.PodLogLine
	var nbytes int

	buf := bufio.NewReader(r)
	for {
		b, err := buf.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}

		nbytes += len(b)

		if len(b) > 0 {
			ll := bytesToLogLine(b)
			logs = append(logs, ll)
		}

		if err == io.EOF {
			break
		}
	}

	// If we're within 1K of the byte limit, remove the last line since it may have be truncated (the byte limit is supposedly
	// not exact). Client will need to do a subsequent request for more data using the last timestamp.
	if nbytes >= (limitBytes - 1024) {
		logs = logs[:len(logs)-1]
	}

	ret := &k8sapiv1.GetPodLogsResponse{Logs: logs}

	// Find the last timestamp and return it so the client can use it in subsequent requests.
	for i := len(logs) - 1; i >= 0; i-- {
		if logs[i].Ts != "" {
			ret.LatestTs = logs[i].Ts
			break
		}
	}

	return ret, nil
}

func protoOptsToK8sOpts(in *k8sapiv1.PodLogsOptions) (*v1.PodLogOptions, error) {
	ret := &v1.PodLogOptions{
		Timestamps: true,
	}
	if in == nil {
		return ret, nil
	}

	ret.Previous = in.Previous

	if in.TailNumLines != 0 {
		ret.TailLines = pointer.Int64(in.TailNumLines)
	}
	if in.SinceTs != "" {
		ts, err := time.Parse(rfc3339NanoFixed, in.SinceTs)
		if err != nil {
			return nil, err
		}
		ret.SinceTime = &metav1.Time{Time: ts}
	}
	if in.ContainerName != "" {
		ret.Container = in.ContainerName
	}
	return ret, nil
}

func byteToStringTrimmed(b []byte) string {
	return strings.TrimSuffix(string(b), "\n")
}

func bytesToLogLine(b []byte) *k8sapiv1.PodLogLine {
	idx := bytes.Index(b, tsDelimiter)
	if idx < 0 {
		return &k8sapiv1.PodLogLine{S: byteToStringTrimmed(b)}
	}

	ts := string(b[:idx])
	if _, err := time.Parse(rfc3339NanoFixed, ts); err != nil {
		return &k8sapiv1.PodLogLine{S: byteToStringTrimmed(b)}
	}

	return &k8sapiv1.PodLogLine{
		Ts: ts,
		S:  byteToStringTrimmed(b[idx+1:]),
	}
}
