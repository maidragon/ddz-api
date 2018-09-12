package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"os"
	"io/ioutil"

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

type Deck struct {
	LordCards          []int `json:"l"`
	Farmer1Cards       []int `json:"f1"`
	Farmer2Cards       []int `json:"f2"`
	LastPlayerCards    []int `json:"lastCard"`
	PlayerIdentity     int32  `json:"Cur_Identity"`
	LastPlayerIdentity int32  `json:"Last_Identity"`
}

type AllDecksResponse struct {
	Decks map[string]([]Deck) `json:"decks"`
}

type GameTableResponse struct {
	Snapshots map[string](ddz.GameTableElement) `json:"snapshots"`
}

type DdzReponse struct {
	Status bool `json:"status"`
	Handcard []int `json:"handcard"`
}

type SnapshotsReponse struct {
	Status bool `json:"status"`
	Snapshots []Snapshot `json:"snapshots"`
}

type Snapshot struct {
	Playeridentity       int32    `json:"player_identity"`
	LordHandcard         []int   `json:"lord_handcard"`
	Farmer1Handcard      []int   `json:"farmer1_handcard"`
	Farmer2Handcard      []int   `json:"farmer2_handcard"`
	LastIdentity         int32    `json:"last_identity"`
	LastPlaycard         []int   `json:"last_playcard"`
	Result               []int   `json:"result"`
}

const ip string = IP_ADDRESS //"0.0.0.0:50001"
const apiAddress = "https://localhost:3005"

/**
 *
 */
func GrpcClientRobot(data GrpcFormatData) (err error, handcard []byte) {
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
		fmt.Println(err)
		return err, []byte{}
	} else {
		return nil, r.Handcard
	}
}

func GRPCGameTableRobot(data GrpcFormatData) (err error, snapshots []ddz.GameTableElement) {
	conn, err := grpc.Dial(ip, grpc.WithInsecure())
	if err != nil {
		fmt.Println("did not connect :", err.Error())
	} else {
		fmt.Println("connect succ:" + conn.Target())
	}
	defer conn.Close()
	c := ddz.NewGameTableServiceClient(conn)

	r, err := c.Play(context.Background(), &ddz.RobotRequest{
		Playeridentity:  data.PlayerIdentity,
		LordHandcard:    data.LordCards,
		Farmer1Handcard: data.Farmer1Cards,
		Farmer2Handcard: data.Farmer2Cards,
		LastIdentity:    data.LastPlayerIdentity,
		LastPlaycard:    data.LastPlayerCards,
	})
	
	if err != nil {
		fmt.Println(err)
		return err, []ddz.GameTableElement{}
	} else {

		snapshots := []ddz.GameTableElement{}

		for _, snapshot := range r.Element {
			snapshots = append(snapshots, *snapshot)
		}
		
		// AllSnapshotsResponse := make(map[string]([]ddz.GameTableElement))
		// AllSnapshotsResponse["snapshots"] = snapshots
	
		// w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode(AllSnapshotsResponse)
		// w.WriteHeader(http.StatusOK)
		// fmt.Println(r.Element)
		return nil, snapshots
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
		w.WriteHeader(http.StatusOK)
	} else {
		decoder := json.NewDecoder(r.Body)
		var data GrpcFormatData
		err := decoder.Decode(&data)
		if err != nil {
			panic(err)
		}
		log.Println(data)

		err, handcard := GrpcClientRobot(data)

		errResponse := DdzReponse{
			Handcard: []int{},
			Status: false,
		}
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(errResponse)
			return;
		}


		handcardIntArray := []int{};

		for _, v := range handcard {
			handcardIntArray = append(handcardIntArray, int(v))
		}
		responseData := DdzReponse{
			Handcard: handcardIntArray,
			Status: true,
		}
		fmt.Println(responseData)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseData)
	}

	// GrpcClientRobot()
	
	return
}

func gameTableRequestHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method == "OPTIONS" {
		// handle preflight requests
		w.WriteHeader(http.StatusOK)
	} else {
		decoder := json.NewDecoder(r.Body)
		var data GrpcFormatData
		err := decoder.Decode(&data)
		if err != nil {
			panic(err)
		}
		log.Println(data)

		err, snapshots := GRPCGameTableRobot(data)

		errResponse := SnapshotsReponse{
			Snapshots: []Snapshot{},
			Status: false,
		}
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(errResponse)
			return;
		}

		formatedSnapshots := []Snapshot{}

		for _, s := range snapshots {
			formatedSnapshot := Snapshot{
				Farmer1Handcard: ConvertByteArrayToIntArray(s.Farmer1Handcard),
				Farmer2Handcard: ConvertByteArrayToIntArray(s.Farmer2Handcard),
				LordHandcard: ConvertByteArrayToIntArray(s.LordHandcard),
				Playeridentity: s.Playeridentity,
				LastIdentity: s.LastIdentity,
				LastPlaycard: ConvertByteArrayToIntArray(s.LastPlaycard),
				Result: ConvertByteArrayToIntArray(s.Result),
			}
			formatedSnapshots = append(formatedSnapshots, formatedSnapshot)
		}
		responseData := SnapshotsReponse {
			Status: true,
			Snapshots: formatedSnapshots,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseData)
	}	
	return
}

func getDeckInfo(jsonName string) Deck {
	jsonFile, err := os.Open("../../decks/" + jsonName)
	if err != nil {
		fmt.Println(err)
	}


	fmt.Println("Successfully Opened " + jsonName)
	defer jsonFile.Close()


	byteValue, _ := ioutil.ReadAll(jsonFile)

	var deckInfo Deck
	
	json.Unmarshal(byteValue, &deckInfo)
	return deckInfo;
}

func readJSONFileHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	
	decks := []Deck{}
	files, err := ioutil.ReadDir("../../decks/")
	if err != nil {
		log.Fatal(err)
	}
		
	for _, f := range files {
		decks = append(decks, getDeckInfo(f.Name()))
	}
	
	AllDecksResponse := make(map[string]([]Deck))
	AllDecksResponse["decks"] = decks

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AllDecksResponse)
	w.WriteHeader(http.StatusOK)
	return
}

func main() {

	port := flag.Int("port", 3005, "server listening port")
	flag.Parse()

	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		requestHandler(w, r)
	})

	http.HandleFunc("/gametable", func(w http.ResponseWriter, r *http.Request) {
		gameTableRequestHandler(w, r)
	})

	http.HandleFunc("/decks", func(w http.ResponseWriter, r *http.Request) {
		readJSONFileHandler(w, r)
	})
	fmt.Printf("the server is listening on %d\n", *port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
