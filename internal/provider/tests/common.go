// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"os"
	"strings"
)

const (
	spaceName        string = "TorqueTerraformProvider"
	web_hook         string = "https://web_hook.com"
	notificationName string = "notification"
	webhook_token    string = "token"
)

var version = os.Getenv("VERSION")
var minorVresion = strings.Split((version), ".")
var index = minorVresion[1]
var fullSpaceName = spaceName + index
