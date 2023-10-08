package times

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func MustSmartParseTime(str string) (tm time.Time) {
	var err error
	tm, err = smartParseTime(str)
	if err != nil {
		hintInternal(err, "smartParseTime failed")
	}
	return
}

func SmartParseTime(str string) (tm time.Time, err error) {
	return smartParseTime(str)
}

var knownDateTimeFormats []string
var onceFormats sync.Once

func smartParseTime(str string) (tm time.Time, err error) {
	onceFormats.Do(func() {
		knownDateTimeFormats = []string{
			"2006-01-02 15:04:05.999999999 -0700",
			"2006-01-02 15:04:05.999999999Z07:00",
			"2006-01-02 15:04:05.999999999",
			"2006-01-02 15:04:05.999",
			"2006-01-02 15:04:05",
			"2006-01-02",
			"2006/01/02",
			"01/02/2006",
			"01-02",

			"2006-1-2 15:4:5.999999999 -0700",
			"2006-1-2 15:4:5.999999999Z07:00",
			"2006-1-2 15:4:5.999999999",
			"2006-1-2 15:4:5.999",
			"2006-1-2 15:4:5",
			"2006-1-2",
			"2006/1/2",
			"1/2/2006",
			"1-2",

			"15:04:05.999999999 -0700",
			"15:04:05.999999999Z0700",
			"15:04:05.999999999",
			"15:04:05.999",
			"15:04:05",
			"15:04",

			"15:4:5.999999999 -0700",
			"15:4:5.999999999Z0700",
			"15:4:5.999999999",
			"15:4:5.999",
			"15:4:5",
			"15:4",

			time.RFC3339Nano,
			time.RFC3339,
			time.RFC1123Z,
			time.RFC1123,
			time.RFC850,
			time.RFC822Z,
			time.RFC822,
			time.RubyDate,
			time.UnixDate,
			time.ANSIC,
		}
	})

	for _, layout := range knownDateTimeFormats {
		if tm, err = time.Parse(layout, str); err == nil {
			break
		}
	}
	return
}

func hintInternal(err error, msg string) {
	fmt.Fprint(os.Stderr, msg, " | Error: \n", err)
}
