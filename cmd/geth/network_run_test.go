// Copyright 2016 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cespare/cp"
)

// Quorum
//
// Start a raft network of 2 nodes and check if the `geth` process shutdowns gracefully on both
func TestGethRaftStopsGracefully(t *testing.T) {
	geth1, dataDir1 := startGethRaft(t, 1)
	defer os.RemoveAll(dataDir1)
	geth2, dataDir2 := startGethRaft(t, 2)
	defer os.RemoveAll(dataDir2)

	time.Sleep(10 * time.Second) // wait until network is running

	geth1.Interrupt()
	geth2.Interrupt()

	time.Sleep(10 * time.Second) // wait until network is running

	geth1.ExpectExit()
	geth2.ExpectExit()
}

func startGethRaft(t *testing.T, id int) (*testgeth, string) {
	nodeDir := fmt.Sprintf("testdata/raft-network/node%d", id)

	dataDir := nodeDir + "/data"

	os.RemoveAll(dataDir)
	if err := cp.CopyAll(dataDir, nodeDir+"/original-data"); err != nil {
		t.Fatal(err)
	}

	runGeth(t, "--nousb", "--datadir", dataDir, "init", "testdata/raft-network/genesis.json").WaitExit()

	geth := runGethWithRaftConsensus(t,
		"--datadir", dataDir, "--nodiscover", "--networkid", "31337", "--raftport", fmt.Sprintf("%d", 50000+id), "--raftblocktime", "1", "--port", fmt.Sprintf("%d", 21000+id), "--nousb", "--rpc", "--rpcaddr", "0.0.0.0", "--rpcport", fmt.Sprintf("%d", 22000+id), "--rpcapi", "admin,db,eth,debug,miner,net,shh,txpool,personal,web3,quorum,raft", "--ipcdisable")

	return geth, dataDir
}
