/*
   this is part of gosnusp - a SNUSP esolang interpreter written in Go
   Copyright (C) 2014  Mauro Panigada

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.

====

   Main email address: shintakezou AT@AT gmail DOT.DOT com

*/

package snusp

import "./dir"

const (
	Start string = "$"
	Up           = ":"
	Down         = ";"
	Left         = "<"
	Right        = ">"
	Lurd         = "\\"
	Ruld         = "/"
	Incr         = "+"
	Decr         = "-"
	Read         = ","
	Write        = "."
	Skip         = "!"
	SkipZ        = "?"
	Enter        = "@"
	Leave        = "#"
	Split        = "&"
	Rand         = "%"
	NoOp         = " "
)

type Pos struct {
	X int
	Y int
}

type Size struct {
	W int
	H int
}

var LurdMap map[dir.Dir]dir.Dir
var RuldMap map[dir.Dir]dir.Dir

func init() {
	LurdMap = make(map[dir.Dir]dir.Dir)
	RuldMap = make(map[dir.Dir]dir.Dir)

	LurdMap[dir.Left] = dir.Up
	LurdMap[dir.Up] = dir.Left
	LurdMap[dir.Right] = dir.Down
	LurdMap[dir.Down] = dir.Right

	RuldMap[dir.Left] = dir.Down
	RuldMap[dir.Up] = dir.Right
	RuldMap[dir.Right] = dir.Up
	RuldMap[dir.Down] = dir.Left
}
