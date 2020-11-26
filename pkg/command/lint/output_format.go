package lint

import (
	"golang.stackrox.io/kube-linter/internal/flagutil"
)

var (
	jsonOutputFormat  = "json"
	plainOutputFormat = "plain"

	formatValueFactory = flagutil.NewEnumValueFactory("output format", []string{jsonOutputFormat, plainOutputFormat})
)
