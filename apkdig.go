package main

/*
 * Copyright (c) 2014 Floor Terra <floort@gmail.com>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

import (
	"errors"
	"fmt"
	"os"

	"github.com/floort/apkdig/apk"
	"github.com/floort/apkdig/axml"
	"github.com/floort/apkdig/dex"
)

var COMMANDS = map[string](func([]string) error){}

// Errors
var ERR_NOARGS = errors.New("No arguments specified.")
var ERR_TOOMANYARGS = errors.New("Too many arguments specified.")

func Help(args []string) (err error) {
	if len(args) == 1 {
		fmt.Println("Usage:", os.Args[0], "[command] [arguments..]")
		fmt.Println("")
		fmt.Println("Commands:")
		for cmd := range COMMANDS {
			fmt.Println("\t", cmd)
		}
		fmt.Println("")
		fmt.Println("For detailed help use:", os.Args[0], " help [command]")
	}
	return nil
}

func Manifest(args []string) (err error) {
	if len(args) == 0 {
		return ERR_NOARGS
	}
	if len(args) == 1 {
		return errors.New("No file specified.")
	}
	if len(args) > 2 {
		return ERR_TOOMANYARGS
	}
	if args[1] == "-h" {
		// Print help information
		fmt.Println(args[0], "-h\t\tPrint usage.")
		fmt.Println(args[0], "[filename]\tPrint manifest xml.")
		return nil
	} else {
		a, err := apk.OpenAPK(args[1])
		if err != nil {
			return err
		}
		defer a.Close()
		manifestfile, err := a.OpenFile("AndroidManifest.xml")
		if err != nil {
			return err
		}
		manifest, err := axml.ReadAXML(manifestfile)
		if err != nil {
			return err
		}
		fmt.Println(manifest.XML)
		return nil
	}
}

func DexStrings(args []string) (err error) {
	if len(args) == 0 {
		return ERR_NOARGS
	}
	if len(args) == 1 {
		return errors.New("No file specified.")
	}
	if len(args) > 2 {
		return ERR_TOOMANYARGS
	}
	if args[1] == "-h" {
		// Print help information
		fmt.Println(args[0], "-h\t\tPrint usage.")
		fmt.Println(args[0], "[filename]\tPrint strings in dex file.")
		return nil
	} else {
		a, err := apk.OpenAPK(args[1])
		if err != nil {
			return err
		}
		defer a.Close()
		dexfile, err := a.OpenFile("classes.dex")
		if err != nil {
			return err
		}
		dx, err := dex.ReadDex(dexfile)
		if err != nil {
			return err
		}
		for i := range dx.Strings {
			fmt.Printf("%v\n", dx.Strings[i])
		}
		return nil
	}
}

func Code(args []string) (err error) {
	if len(args) == 0 {
		return ERR_NOARGS
	}
	if len(args) == 1 {
		return errors.New("No file specified.")
	}
	if len(args) > 2 {
		return ERR_TOOMANYARGS
	}
	if args[1] == "-h" {
		// Print help information
		fmt.Println(args[0], "-h\t\tPrint usage.")
		fmt.Println(args[0], "[filename]\tPrint classes and methods in dex file.")
		return nil
	} else {
		a, err := apk.OpenAPK(args[1])
		if err != nil {
			return err
		}
		defer a.Close()
		dexfile, err := a.OpenFile("classes.dex")
		if err != nil {
			return err
		}
		dx, err := dex.ReadDex(dexfile)
		if err != nil {
			return err
		}
		for classname, class := range dx.Classes {
			fmt.Println(classname)
			for _, method := range class.Methods {
				fmt.Printf("\t%s\n", method.Name)
			}
		}
		return nil
	}
}

func main() {
	// Fill Commands table
	COMMANDS["?"] = Help
	COMMANDS["help"] = Help
	COMMANDS["manifest"] = Manifest
	COMMANDS["dexstrings"] = DexStrings
	COMMANDS["code"] = Code
	if len(os.Args) == 1 {
		// No arguments are given; print help
		_ = Help([]string{"help"})
		return
	}
	f, ok := COMMANDS[os.Args[1]]
	if ok {
		err := f(os.Args[1:])
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Unknown command:", os.Args[1])
	}

}
