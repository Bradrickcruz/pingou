package main

import (
	// embute a base de fusos no binário (~450 KB): o fuso das Configurações
	// funciona mesmo sem /usr/share/zoneinfo na máquina ou na imagem
	_ "time/tzdata"

	"github.com/Bradrickcruz/pingou/cmd/pingou/commands"
)

func main() {
	commands.Execute()
}
