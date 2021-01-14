package cmd

import (
	"fmt"
	"strings"
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
			return val.(string)
		}
		attributes = string(a)

		return attributes
	}
}

func timeInMSTransformer() text.Transformer {
	return func(val interface{}) string {
		t, ok := val.(int64)
		if !ok {
			return val.(string)
		}

		return fmt.Sprintf("%s", time.Unix(0, t*int64(time.Millisecond)))
	}
}

func durationTransformer() text.Transformer {
	return func(val interface{}) string {

		t, ok := val.(int64)
		if !ok {
			return val.(string)
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

func taskDetailTransformer() text.Transformer {
	return func(val interface{}) string {
		s := val.(string)
		// Don't wrap if we only have one (the most common case)
		if strings.Count(s, ",") <= 1 {
			return s
		}
		return strings.Replace(s, ",", ",\n ", -1)
	}
}
