package main

import (
    "log"
    "os"

    "github.com/ahmetkarakayaoffical/scnorionplus-console/internal/commands"
    "github.com/urfave/cli/v2"
)

func main() {
    app := &cli.App{
        Name:      "scnorionplus-console",
        Commands:  getCommands(),
        Usage:     "SCNOrion Plus - Kurumsal Uç Nokta Yönetim Konsolu",
        Authors:   []*cli.Author{{Name: "SCN GLOBAL", Email: "info@scn-global.com"}},
        Version:   "1.0.0",
        Copyright: "© 2024-2025 SCN GLOBAL - Tüm hakları saklıdır. <https://www.scn-global.com>",
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}

func getCommands() []*cli.Command {
    return []*cli.Command{
        commands.StartConsole(),
        commands.StopConsole(),
    }
}