package lint

import (
	"golang.stackrox.io/kube-linter/internal/flagutil"
)

const (
	jsonOutputFormat  = "json"
	plainOutputFormat = "plain"
)

var (
	formatValueFactory = flagutil.NewEnumValueFactory("Output format", []string{jsonOutputFormat, plainOutputFormat})
)
