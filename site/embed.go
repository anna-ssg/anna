package site

import "embed"

//go:embed content/** layout/** public/** static/**
var FS embed.FS
