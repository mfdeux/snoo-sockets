package server

import (
	"encoding/json"
	"fmt"

	"github.com/mfdeux/pushshift"
)

func newFirehose() chan interface{} {
	client := pushshift.NewClient("testClient/0.1.0")
	return client.StreamFirehose(nil)
}

func broadCastMessages(things chan interface{}, hub *Hub) {
	for thing := range things {
		switch t := thing.(type) {
		case *pushshift.Comment:
			b, err := json.Marshal(t)
			if err != nil {
				fmt.Println("error converting comment")
				break
			}
			hub.broadcast <- b
			break
		case *pushshift.Submission:
			b, err := json.Marshal(t)
			if err != nil {
				fmt.Println("error converting submission")
				break
			}
			hub.broadcast <- b
			break
		default:
			fmt.Printf("I don't know about type %T!\n", t)
		}

	}
}
