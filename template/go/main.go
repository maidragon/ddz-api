package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"./ddz"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

/**
 *
 */
type GrpcFormatData struct {
	LordCards          []byte `json:"lordCards"`
	Farmer1Cards       []byte `json:"farmer1Cards"`
	Farmer2Cards       []byte `json:"farmer2Cards"`
	LastPlayerCards    []byte `json:"lastPlayerCards"`
	PlayerIdentity     int32  `json:"playerIdentity"`
	LastPlayerIdentity int32  `json:"lastPlayerIdentity"`
}

const ip string = IP_ADDRESS //"0.0.0.0:50001"
const apiAddress = "https://localhost:3005"

/**
 *
 */
func GrpcClientRobot(data GrpcFormatData) {
	conn, err := grpc.Dial(ip, grpc.WithInsecure())
	if err != nil {
		fmt.Println("did not connect :", err.Error())
	} else {
		fmt.Println("connect succ:" + conn.Target())
	}
	defer conn.Close()
	c := ddz.NewRobotServiceClient(conn)

	r, err := c.Play(context.Background(), &ddz.RobotRequest{
		Playeridentity:  data.PlayerIdentity,
		LordHandcard:    data.LordCards,
		Farmer1Handcard: data.Farmer1Cards,
		Farmer2Handcard: data.Farmer2Cards,
		LastIdentity:    data.LastPlayerIdentity,
		LastPlaycard:    data.LastPlayerCards,
	})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("%+x\n", r.Handcard)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Methods", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type,Access-Token")
	(*w).Header().Set("Access-Control-Expose-Headers", "*")
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method == "OPTIONS" {
		// handle preflight requests
	} else {
		decoder := json.NewDecoder(r.Body)
		var data GrpcFormatData
		err := decoder.Decode(&data)
		if err != nil {
			panic(err)
		}
		log.Println(data)

		GrpcClientRobot(data)
	}

	// GrpcClientRobot()
	w.WriteHeader(http.StatusCreated)
	return
}

func main() {

	port := flag.Int("port", 3005, "server listening port")
	flag.Parse()

	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		requestHandler(w, r)
	})
	fmt.Printf("the server is listening on %d\n", *port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
