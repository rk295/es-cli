package cmd

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/hokaccha/go-prettyjson"
	"github.com/jedib0t/go-pretty/v6/text"
)

func prettyJSONTransformer() text.Transformer {
	return func(val interface{}) string {

		var attributes string

		a, err := prettyjson.Marshal(val)
		if err != nil {
			// TODO: Unsure what to do here?
			return "error: failed to marshal"
		}
		attributes = string(a)

		return attributes
	}
}

func timeInMSTransformer() text.Transformer {
	return func(val interface{}) string {
		t, ok := val.(int64)
		if !ok {
			return "error: failed to parse the time"
		}

		return fmt.Sprintf("%s", time.Unix(0, t*int64(time.Millisecond)))
	}
}

func durationTransformer() text.Transformer {
	return func(val interface{}) string {

		t, ok := val.(int64)
		if !ok {
			return "error: failed to parse the time"
		}

		duration := time.Duration(t) * time.Nanosecond
		var prettyDuration string

		if duration > redDuration {
			prettyDuration = color.RedString(fmt.Sprintf("%v", duration))
		} else if duration > yellowDuration {
			prettyDuration = color.YellowString(fmt.Sprintf("%v", duration))
		} else {
			prettyDuration = color.GreenString(fmt.Sprintf("%v", duration))
		}

		return prettyDuration
	}
}
