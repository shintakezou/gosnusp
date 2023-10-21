/*
   this is part of gosnusp - a SNUSP esolang interpreter written in Go
   Copyright (C) 2023  Mauro Panigada

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

package dir

type Dir struct {
	Dx int
	Dy int
}

var (
	Left  Dir = Dir{-1, 0}
	Right     = Dir{1, 0}
	Up        = Dir{0, -1}
	Down      = Dir{0, 1}
)
