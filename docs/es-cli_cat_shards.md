## es-cli cat shards

Lists the shards in the cluster.

### Synopsis

Lists the shards in the cluster. Supports sorting and changing the byte unit to use.

```
es-cli cat shards [--state=initializing|relocating|started|unassigned] [--unassigned] [flags]
```

### Options

```
  -b, --byte string        Byte unit to use. Valid values are: "b", "k", "kb", "m", "mb", "g", "gb", "t", "tb", "p" or "pb" (default "mb")
  -h, --help               help for shards
  -s, --sort string        Field to sort by, possible to list multiple comma separated See https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-shards.html for full list of fields (default "index")
      --state string       Filter shards by state, possible values are: initializing,relocating,started,unassigned
      --unassigned         Print just the unassigned shards
      --unassigned-count   Just print the number of unassigned shards
```

### Options inherited from parent commands

```
      --colour          Enable/Disable Colour. (default true)
      --es-url string   url for elasticsearch (default "http://localhost:9200/")
  -m, --markdown        Produce Markdown output
```

### SEE ALSO

* [es-cli cat](es-cli_cat.md)	 - cat api commands

###### Auto generated by spf13/cobra on 17-May-2021
