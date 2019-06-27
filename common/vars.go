package common

import (
	"github.com/bwmarrin/discordgo"
)

// ===== BOT =====

// Inline struct as it's not reused and just used to containerize parameters.
var BotVars = struct {
	Commands  chan int // TODO: Establish a enum of sorts
	Guilds    map[string]*discordgo.UserGuild
	HostGuild string
}{
	Commands: make(chan int, 50),
	Guilds:   make(map[string]*discordgo.UserGuild),
}

// ===== CLI =====

var CliVars = struct {
	Commands    chan int // TODO: Establish a enum of sorts
	Initialized bool
}{
	Commands:    make(chan int, 50),
	Initialized: false,
}
