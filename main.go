package main

import (
	"bytes"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gcash/bchd/chaincfg/chainhash"
	"github.com/gcash/bchd/rpcclient"
	"github.com/gcash/bchd/wire"
	"github.com/gcash/bchutil"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

var bchdClient *rpcclient.Client
var sqliteDB *sql.DB
var regexRequest *regexp.Regexp

func main() {
	regexRequest = regexp.MustCompile(`^/CashAccount/(\d+)/([a-zA-Z0-9_]{1,99})/?$`)
	http.HandleFunc("/CashAccount/", handler)

	// create default config
	cashaccHomeDir := bchutil.AppDataDir("cashAccount", false)
	cfgFile := filepath.Join(cashaccHomeDir, "cashAccountAPId.conf")
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		err = ioutil.WriteFile(cfgFile, []byte(defaultConfig), 0600)
		if err != nil {
			log.Fatal(err)
		}
	}

	// read config
	cfg, err := ini.Load(cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	// basic config check
	if cfg.Section("cashAccountAPId").Key("rpchost").String() == "" ||
		cfg.Section("cashAccountAPId").Key("rpcuser").String() == "" ||
		cfg.Section("cashAccountAPId").Key("rpcpass").String() == "" ||
		cfg.Section("cashAccountAPId").Key("rpcendpoint").String() == "" ||
		cfg.Section("cashAccountAPId").Key("webserverbindaddr").String() == "" {
		log.Fatal("Configuration not set up. Please set up the configuration! File:" + cfgFile)
	}

	// sqlite
	cashaccDbFileName := filepath.Join(cashaccHomeDir, "db.sqlite")
	sqliteDB, err = sql.Open("sqlite3", cashaccDbFileName)
	if err != nil {
		log.Fatal(err)
	}

	// bchd configuration
	log.Println("Connecting to bchd...")
	bchdHomeDir := bchutil.AppDataDir("bchd", false)
	certs, err := ioutil.ReadFile(filepath.Join(bchdHomeDir, "rpc.cert"))
	if err != nil {
		log.Fatal(err)
	}
	connCfg := &rpcclient.ConnConfig{
		Host:         cfg.Section("cashAccountAPId").Key("rpchost").String(),
		Endpoint:     cfg.Section("cashAccountAPId").Key("rpcendpoint").String(),
		User:         cfg.Section("cashAccountAPId").Key("rpcuser").String(),
		Pass:         cfg.Section("cashAccountAPId").Key("rpcpass").String(),
		Certificates: certs,
	}
	bchdClient, err = rpcclient.New(connCfg, &rpcclient.NotificationHandlers{})
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(cfg.Section("cashAccountAPId").Key("webserverbindaddr").String(), nil))
}

//
func handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")
	re := response{}

	m := regexRequest.FindStringSubmatch(r.URL.Path)
	if len(m) == 0 {
		re.Error = "Invalid request"
		_, _ = fmt.Fprint(w, re.Format())
		return
	}

	re.Block, _ = strconv.ParseInt(m[1], 10, 64)
	re.Block += 563620
	re.Name = m[2]
	re.Results = []txAndProofPair{}

	rr, err := sqliteDB.Query(`SELECT txid FROM nameindex WHERE name like ? AND block == ?;`, re.Name, re.Block)
	if err != nil {
		re.Error = "Database error"
		_, _ = fmt.Fprint(w, re.Format())
		return
	}

	for rr.Next() {
		txidRaw := []byte("")

		err = rr.Scan(&txidRaw)
		if err != nil {
			re.Error = "Database error"
			_, _ = fmt.Fprint(w, re.Format())
			return
		}

		txid := hex.EncodeToString(txidRaw)
		var txidChainhash chainhash.Hash
		_ = chainhash.Decode(&txidChainhash, txid)

		tx, err := bchdClient.GetRawTransaction(&txidChainhash)
		if err != nil {
			re.Error = "Database error (at GetTransaction)" + " >" + err.Error()
			_, _ = fmt.Fprint(w, re.Format())
			return
		}

		blockHash, err := bchdClient.GetBlockHash(re.Block)
		if err != nil {
			re.Error = "Database error (at GetBlockHash)" + " >" + err.Error()
			_, _ = fmt.Fprint(w, re.Format())
			return
		}

		shit := blockHash.String()
		txOutProof, err := bchdClient.GetTxOutProof([]string{txid}, &shit)
		if err != nil {
			re.Error = "Database error (at GetTxOutProof)" + " >" + err.Error()
			_, _ = fmt.Fprint(w, re.Format())
			return
		}

		var transactionRaw bytes.Buffer
		err = tx.MsgTx().BchEncode(&transactionRaw, 0, wire.BaseEncoding)
		if err != nil {
			re.Error = "Error encoding transaction"
			_, _ = fmt.Fprint(w, re.Format())
			return
		}

		tapp := txAndProofPair{}
		tapp.Tx = hex.EncodeToString(transactionRaw.Bytes())
		tapp.Proof = txOutProof

		re.Results = append(re.Results, tapp)

	}
	_ = rr.Close()

	_, _ = fmt.Fprint(w, re.Format())
}

const defaultConfig string = "[cashAccountAPId]\nrpchost=127.0.0.1:8334\nrpcendpoint=ws\nrpcuser=\nrpcpass=\nwebserverbindaddr=127.0.0.1:8080"

type Config struct {
	cashAccountDBd struct {
		rpchost           string
		rpcendpoint       string
		rpcuser           string
		rpcpass           string
		webserverbindaddr string
	}
}

type response struct {
	Error   string           `json:"error"`
	Name    string           `json:"name"`
	Block   int64            `json:"block"`
	Results []txAndProofPair `json:"results"`
}

type txAndProofPair struct {
	Tx    string `json:"transaction"`
	Proof string `json:"inclusion_proof"`
}

func (r response) Format() string {
	tmp, _ := json.MarshalIndent(r, " ", "    ")
	return string(tmp)
}
