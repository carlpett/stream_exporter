package input

import (
	"testing"
	"time"

	"github.com/valyala/fasttemplate"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

var testMessages = []map[string]interface{}{
	{
		"string1": "hello",
		"int":     13,
		"string2": "longlonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglong",
		"time":    time.Now(),
	},
	{
		"string1": "both-longlonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglong",
		"int":     61,
		"string2": "longlonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglong",
		"time":    time.Now(),
	},
	{
		"string1": "only-string1",
		"int":     123718,
		"time":    time.Now(),
	},
	{
		"int":     0,
		"string2": "only-string2",
		"time":    time.Now(),
	},
	{
		"string1": "hello",
		"int":     -1541,
		"string2": "world",
	},
}

func BenchmarkMessageHandler(b *testing.B) {
	t := fasttemplate.New("[string1][int][string2][time]", "[", "]")
	ch := make(chan string, b.N)
	logPartCh := make(chan format.LogParts, b.N)

	for n := 0; n < b.N; n++ {
		logPartCh <- testMessages[n%len(testMessages)]
	}
	close(logPartCh)

	for n := 0; n < b.N; n++ {
		messageHandler(t, ch, logPartCh)
	}
}
