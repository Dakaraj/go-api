package main

import (
	jsonparse "encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

// Args holds arguments passed to JSON RPC service
type Args struct {
	ID string
}

// Book struct holds Book JSON structure
type Book struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Author string `json:"author,omitempty"`
}

// JSONServer is a struct that handles an RPC server
type JSONServer struct{}

// GiveBookDetail responds with a details of a book
func (t *JSONServer) GiveBookDetail(r *http.Request, args *Args, reply *Book) error {
	var books []Book
	// cmd := exec.Command("C:\\Windows\\System32\\cmd.exe", "cd")
	// var stdout bytes.Buffer
	// cmd.Stdout = &stdout
	// cmdErr := cmd.Run()
	// if cmdErr != nil {
	// 	log.Println("CMD error:", cmdErr)
	// }
	// log.Println("cur dur:", stdout.String())

	// Read JSON file and load data
	raw, readerr := ioutil.ReadFile("./jsonRPCServer/books.json")
	if readerr != nil {
		log.Println("error:", readerr)
	}

	// Unmarshall JSON raw data into books array
	marshalerr := jsonparse.Unmarshal(raw, &books)
	if marshalerr != nil {
		log.Println("error:", marshalerr)
		os.Exit(1)
	}

	// Iterate over each book to find the given book
	log.Println(args.ID)
	for _, book := range books {
		if book.ID == args.ID {
			// If book found, fill reply with it
			*reply = book
			break
		}
	}
	return nil
}

func main() {
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(new(JSONServer), "")
	r := mux.NewRouter()
	r.Handle("/rpc", s)
	http.ListenAndServe(":80", r)
}
