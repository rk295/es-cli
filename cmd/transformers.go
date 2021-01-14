package cmd

import (
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
