// Copyright 2015 by David A. Golden. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: grave <dirname>")
	}

	home := os.Getenv("HOME")
	configdir := filepath.Join(home, ".grave")
	s, err := os.Stat(configdir)
	if err != nil || !s.Mode().IsDir() {
		log.Fatal(configdir, " is not a valid config directory")
	}

	profile := "default" // XXX eventually set with flags

	generate(configdir, profile, os.Args[1])

}

func generate(configdir, profile, target string) error {

	profiledir := filepath.Join(configdir, profile)
	s, err := os.Stat(profiledir)
	if err != nil || !s.Mode().IsDir() {
		log.Fatal(profiledir, " is not a valid profile directory")
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	targetdir := filepath.Join(cwd, target)

	err = os.Mkdir(targetdir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	err = filepath.Walk(profiledir, genWalker(profiledir, cwd, target))
	if err != nil {
		log.Fatalf("Error creating %s: %s", targetdir, err)
	}

	return nil
}

// genWalker generates a file path walker that contains a closure to
// the original path so that a path relative to it can be constructed
func genWalker(pd, cwd, td string) filepath.WalkFunc {

	absTd := filepath.Join(cwd, td)

	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(pd, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(absTd, relPath)

		fmt.Printf("%s/%s\n", td, relPath) // XXX eventually add verbose flag for this

		err = os.MkdirAll(filepath.Dir(relPath), 755)
		if err != nil {
			return err
		}

		src, err := os.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()

		dst, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			return err
		}

		return nil
	}
}
