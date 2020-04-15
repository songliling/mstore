package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/tendermint/iavl"
	tmdb "github.com/tendermint/tendermint/libs/db"
)

// main runs the main program
func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: %v <dbpath>\n", os.Args[0])
		os.Exit(1)
	}
	versionStr := os.Args[3]
	version, err := strconv.ParseInt(versionStr, 10, 64)
	if err != nil {
		fmt.Printf("parse int error: %s", err)
		os.Exit(1)
	}
	err = run(os.Args[1], os.Args[2], version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err.Error())
		os.Exit(1)
	}
}

// run runs the command with normal error handling
func run(dbPath, moduleName string, version int64) error {
	version, _, err := runExportAndImport(dbPath, moduleName, version)
	if err != nil {
		return err
	}
	return nil
}

// runExport runs an export benchmark and returns a map of store names/export nodes
func runExportAndImport(dbPath, moduleName string, targetVersion int64) (int64, map[string][]*iavl.ExportNode, error) {
	start := time.Now()
	ldb, err := tmdb.NewGoLevelDB(moduleName, dbPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db := tmdb.NewPrefixDB(ldb, []byte("s/_/"))
	tree, err := iavl.NewMutableTree(db, 0)
	if err != nil {
		return 0, nil, err
	}

	storeVersion, err := tree.LoadVersion(targetVersion)
	if err != nil {
		return 0, nil, err
	}
	if storeVersion == 0 {
		fmt.Printf("%-13v\n", moduleName)
		os.Exit(1)
	}

	fmt.Printf("tree.WorkingHash(): %x", tree.WorkingHash())

	itree, err := tree.GetImmutable(storeVersion)
	if err != nil {
		return 0, nil, err
	}
	exporter := itree.Export()
	defer exporter.Close()

	// TODO init import new db
	newDB := tmdb.NewDB(moduleName+"_new", tmdb.GoLevelDBBackend, dbPath)
	newTree, err := iavl.NewMutableTree(newDB, 0)
	if err != nil {
		return 0, nil, err
	}
	importer, err := newTree.Import(targetVersion)
	if err != nil {
		return 0, nil, err
	}
	defer importer.Close()

	for {
		node, err := exporter.Next()
		if err == iavl.ExportDone {
			break
		} else if err != nil {
			return 0, nil, err
		}
		// TODO import
		err = importer.Add(node)
		if err != nil {
			fmt.Printf("error importer add node: err: %s", err.Error())
			break
		}
	}
	err = importer.Commit()
	if err != nil {
		return 0, nil, err
	}

	fmt.Printf("newTree.WorkingHash(): %x", newTree.WorkingHash())
	fmt.Printf("program cost time: %s", time.Since(start).String())

	return storeVersion, nil, nil
}
