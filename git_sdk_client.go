package main

import (
	sdk "git-platform-sdk"
)

type SdkClient struct {
	lb sdk.LabelClient
	pr sdk.PRClient
}

var sc = SdkClient{}
