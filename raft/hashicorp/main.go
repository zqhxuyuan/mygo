package main

import (
	// "encoding/json"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

type Config struct {
	Bind    string `json:bind`
	DataDir string `json:data_dir`
}

type Word struct {
	words string
}

func (*Word) Apply(l *raft.Log) interface{} {
	return nil
}

func (*Word) Snapshot() (raft.FSMSnapshot, error) {
	return new(WordSnapshot), nil
}

func (*Word) Restore(snap io.ReadCloser) error {
	return nil
}

type WordSnapshot struct {
	words string
}

func (snap *WordSnapshot) Persist(sink raft.SnapshotSink) error {
	return nil
}

func (snap *WordSnapshot) Release() {

}

func main() {
	buf, err := ioutil.ReadFile("/Users/baishui/go/src/mygo/raft/hashicorp/config.json")
	if err != nil {
		log.Fatal(err)
	}

	var v Config
	err = json.Unmarshal(buf, &v)

	dataDir := v.DataDir
	os.MkdirAll(dataDir, 0755)

	if err != nil {
		log.Fatal(err)
	}

	cfg := raft.DefaultConfig()
	// cfg.EnableSingleNode = true
	fsm := new(Word)
	fsm.words = "hahaha"

	dbStore, err := raftboltdb.NewBoltStore(path.Join(dataDir, "raft_db"))
	if err != nil {
		log.Fatal(err)
	}
	fileStore, err := raft.NewFileSnapshotStore(dataDir, 1, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
	trans, err := raft.NewTCPTransport(v.Bind, nil, 3, 5*time.Second, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

	//peers := make([]string, 0, 10)
	//peers = raft.AddUniquePeer(peers, "192.168.78.151:12345")
	//peers = raft.AddUniquePeer(peers, "192.168.78.151:12346")
	//peers = raft.AddUniquePeer(peers, "192.168.78.151:12347")
	//peerStore := raft.NewJSONPeers(dataDir, trans)
	//peerStore.SetPeers(peers)

	r, err := raft.NewRaft(cfg, fsm, dbStore, dbStore, fileStore, trans)

	// t := time.NewTicker(time.Duration(5) * time.Second)

	// for {
	// 	select {
	// 	case <-t.C:
	// 		fmt.Println(r.Leader())
	// 	}
	// }

	//defer r.Shutdown()

	if err != nil {
		fmt.Errorf("create raft error: %v", err)
	}
	fmt.Println("raft ", r.Leader)
}
