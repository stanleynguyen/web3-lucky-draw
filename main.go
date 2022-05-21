package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/pat"
	"github.com/joho/godotenv"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	_ "embed"
)

type entrySubmission struct {
	Name      string `json:"name"`
	Signature string `json:"signature"`
}

type entry struct {
	Name          string `json:"name"`
	WalletAddress string `json:"wallet_address"`
}

var drawPot []entry = []entry{}

//go:embed index.html
var indexHTML []byte

func main() {
	if err := godotenv.Load(); err != nil && os.Getenv("GO_ENV") != "production" {
		log.Fatal(err)
	}

	router := pat.New()

	router.Post("/enter", func(w http.ResponseWriter, r *http.Request) {
		var e entrySubmission
		err := json.NewDecoder(r.Body).Decode(&e)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if e.Signature == "" || e.Name == "" {
			http.Error(w, "Missing input", http.StatusBadRequest)
			return
		}

		prefixedHash := crypto.Keccak256Hash(
			[]byte(
				"\x19Ethereum Signed Message:\n" +
					strconv.Itoa(len("WAGMI")) +
					"WAGMI",
			),
		)

		sig := hexutil.MustDecode(e.Signature)
		if sig[64] != 27 && sig[64] != 28 {
			http.Error(w, "Invalid signature", http.StatusBadRequest)
			return
		}
		sig[64] -= 27

		pubKey, err := crypto.SigToPub(prefixedHash.Bytes(), sig)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		recoveredAddr := crypto.PubkeyToAddress(*pubKey)
		drawPot = append(drawPot, entry{e.Name, recoveredAddr.String()})
		w.WriteHeader(200)
	})
	router.Get("/reset", func(w http.ResponseWriter, _ *http.Request) {
		drawPot = []entry{}
		w.WriteHeader(200)
	})
	router.Get("/pot", func(w http.ResponseWriter, _ *http.Request) {
		res := ""
		for _, e := range drawPot {
			name := e.Name
			if len(e.Name) > 7 {
				name = e.Name[:5] + ".."
			}
			res += fmt.Sprintf(
				"%s (%s..%s)\n",
				name,
				e.WalletAddress[:4],
				e.WalletAddress[len(e.WalletAddress)-2:])
		}
		w.Write([]byte(res))
	})
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(indexHTML)
	})

	var addr string
	if os.Getenv("PORT") != "" {
		addr = fmt.Sprintf(":%s", os.Getenv("PORT"))
	} else {
		addr = "localhost:8080"
	}
	http.ListenAndServe(addr, router)
}
