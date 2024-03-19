/*
   DoDefuse - Removes the time bomb from Cave Story playtest builds
   Copyright (C) 2024  patapancakes

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"os"
	"path/filepath"
)

var pattern = [...]byte{0x74, 0x0B, 0x83, 0xBD, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x74, 0x2B}

func main() {
	if len(os.Args) == 0 {
		//fmt.Print("no file specified\n")
		os.Exit(1)
	}

	file, err := os.ReadFile(os.Args[1])
	if err != nil {
		//fmt.Printf("couldn't read file %s: %s\n", os.Args[1], err)
		os.Exit(1)
	}

	var fileOffset int
	var patternFound bool

	for ; fileOffset+len(pattern) < len(file); fileOffset++ {
		for patternOffset, b := range pattern {
			// the CMP arguments can change between versions
			if b == 0xFF {
				continue
			}

			if file[fileOffset+patternOffset] != b {
				break
			}

			if patternOffset+1 == len(pattern) {
				patternFound = true
				break
			}
		}

		if patternFound {
			break
		}
	}

	if !patternFound {
		//fmt.Print("couldn't find bytes to patch\n")
		os.Exit(1)
	}

	file[fileOffset] = 0x90
	file[fileOffset+1] = 0x90

	file[fileOffset+9] = 0x90
	file[fileOffset+10] = 0x90

	err = os.Rename(os.Args[1], os.Args[1][:len(os.Args[1])-len(filepath.Ext(os.Args[1]))] + ".bak")
	if err != nil {
		//fmt.Printf("couldn't rename file %s: %s", os.Args[1], err)
		os.Exit(1)
	}

	err = os.WriteFile(os.Args[1], file, 0777)
	if err != nil {
		//fmt.Printf("couldn't write file %s: %s\n", os.Args[1], err)
		os.Exit(1)
	}

	os.Exit(0)
}
