package gamedata

import "embed"

//go:embed resourcedata/*
//go:embed categorydata/*
var Files embed.FS
