// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"os"
	"strings"
)

var version = os.Getenv("VERSION")
var minorVresion = strings.Split((version), ".")
var index = minorVresion[1]
