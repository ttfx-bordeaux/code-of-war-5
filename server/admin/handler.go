package admin

import (
	"bufio"
	"log"

	"github.com/ttfx-bordeaux/code-of-war-5/server/io"
)

func NewHandler(actions map[string]func()) func(io.Accepter) {
	return func(server io.Accepter) {
		commands := make(chan io.Command)
		go read(server, commands)
		go execute(actions, commands)
	}
}

func read(server io.Accepter, cmds chan io.Command) {
	for {
		conn, err := server.Accept()
		if err != nil {
			continue
		}
		defer conn.Close()
		r := bufio.NewReader(conn)
		req := io.Request{}
		req.Decode(r)
		cmd := io.Command{}
		cmd.Decode(&req)
		cmds <- cmd
	}
}

func execute(actions map[string]func(), cmds chan io.Command) {
	for {
		c := <-cmds
		log.Printf("Command received %v", c)
		if fc, exist := actions[c.Value]; exist {
			fc()
		}
	}
}
