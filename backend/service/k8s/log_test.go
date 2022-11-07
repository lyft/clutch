package k8s

import (
	k8sv1 "github.com/lyft/clutch/backend/api/k8s/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOptsConversion(t *testing.T) {
	{
		o, err := protoOptsToK8sOpts(nil)
		assert.NoError(t, err)
		assert.NotNil(t, o)
		assert.True(t, o.Timestamps)
	}
	{
		o, err := protoOptsToK8sOpts(&k8sv1.GetPodLogsOptions{
			ContainerName: "foo",
			Previous:      true,
			SinceTs:       "2022-11-07T19:30:38.974187286Z",
			TailNumLines:  25,
		})
		assert.NoError(t, err)
		assert.Equal(t, "foo", o.Container)
		assert.True(t, o.Previous)
		assert.EqualValues(t, 25, *o.TailLines)
		assert.Equal(t, "2022-11-07T19:30:38.974187286Z", (*o.SinceTime).Format(rfc3339NanoFixed))
	}
	{
		o, err := protoOptsToK8sOpts(&k8sv1.GetPodLogsOptions{})

		assert.NoError(t, err)
		assert.Equal(t, "", o.Container)
		assert.False(t, o.Previous)
		assert.Nil(t, o.TailLines)
		assert.Nil(t, o.SinceTime)
	}
}

func TestBytesToLogLine(t *testing.T) {
	{
		ll := bytesToLogLine([]byte("2022-11-07T19:30:38.974187286Z Hello world this is it!\n"))
		assert.Equal(t, "2022-11-07T19:30:38.974187286Z", ll.Ts)
		assert.Equal(t, "Hello world this is it!", ll.S)
	}
	{
		ll := bytesToLogLine([]byte("2022-11-07T19:30:38.974187287Z Hello world this is it!"))
		assert.Equal(t, "2022-11-07T19:30:38.974187287Z", ll.Ts)
		assert.Equal(t, "Hello world this is it!", ll.S)
	}
	{
		ll := bytesToLogLine([]byte("2022 Hello world"))
		assert.Equal(t, "", ll.Ts)
		assert.Equal(t, "2022 Hello world", ll.S)
	}
}
