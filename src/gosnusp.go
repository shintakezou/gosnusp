/*
   snusp - a SNUSP esolang interpreter written in Go
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

package main

import (
	"./lang"
	"bufio"
	"container/list"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
)

type Pos struct {
	x int
	y int
}

type Dir struct {
	dx int
	dy int
}

type State struct {
	pos Pos
	dir Dir
}

type Size struct {
	w int
	h int
}

type Snusp struct {
	debug    bool
	loaded   bool
	modular  bool
	bloated  bool
	twist    bool
	pos      Pos
	size     Size
	code     map[int]string
	codeLock sync.RWMutex
	mem      map[Pos]byte
	lock     sync.RWMutex
	sg       sync.WaitGroup
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func (e *Snusp) Get(p Pos) string {
	var res string
	if p.x < 0 || p.y < 0 {
		return lang.NoOp
	}
	e.codeLock.RLock()
	if s, ok := e.code[p.y]; ok {
		if p.x >= len(s) {
			res = lang.NoOp
		} else {
			res = s[p.x : p.x+1]
		}
	}
	e.codeLock.RUnlock()
	return res
}

func (e *Snusp) SetMem(p Pos, val int, overwrite bool) {
	e.lock.Lock()
	if v, ok := e.mem[p]; ok {
		if overwrite {
			e.mem[p] = byte(val % 256)
		} else {
			e.mem[p] = byte((int(v) + val) % 256)
		}
	} else {
		e.mem[p] = byte(val % 256)
	}
	e.lock.Unlock()
}

func (e *Snusp) GetMem(p Pos) byte {
	var res byte
	e.lock.RLock()
	if v, ok := e.mem[p]; ok {
		res = v
	} else {
		res = 0
	}
	e.lock.RUnlock()
	return res
}

func (e *Snusp) Load(fileName string) {
	e.loaded = false
	e.pos.x = 0
	e.pos.y = 0
	e.size.w = 0
	e.size.h = 0
	e.code = make(map[int]string)
	e.mem = make(map[Pos]byte)
	fh, errOpen := os.Open(fileName)
	if errOpen != nil {
		log.Fatal(errOpen)
	}
	defer fh.Close()
	r := bufio.NewReader(fh)
	var (
		s   string
		err error
		ly  int
	)
	for {
		if s, err = r.ReadString('\n'); err == nil {
			e.code[ly] = s
			e.size.w = Max(e.size.w, len(s))
			if px := strings.Index(s, lang.Start); px != -1 {
				e.pos.y = ly
				e.pos.x = px
				if e.debug {
					log.Printf("Start at (%d, %d)", e.pos.x, e.pos.y)
				}
			}
			ly++
		} else if err == io.EOF {
			e.code[ly] = s
			ly++
			e.loaded = true
			e.size.h = ly
			if e.debug {
				log.Printf("w=%d, h=%d", e.size.w, e.size.h)
			}
			break
		} else {
			log.Fatal(err)
			break //?
		}
	}
}

func (e *Snusp) Interpret(p Pos, d Dir, m Pos) {
	if e.debug {
		log.Printf("interpret at (%d,%d), dir (%d,%d)", p.x, p.y, d.dx, d.dy)
	}
	defer e.sg.Done()
	stdin := bufio.NewReader(os.Stdin)
	cstack := list.New()
	for p.x >= 0 && p.x < e.size.w && p.y >= 0 && p.y < e.size.h {
		c := e.Get(p)
		if e.debug {
			log.Print(p, c, d)
		}
		switch c {
		case lang.Left:
			m.x--
		case lang.Right:
			m.x++
		case lang.Up:
			if e.bloated {
				m.y--
			}
		case lang.Down:
			if e.bloated {
				m.y++
			}
		case lang.Incr:
			e.SetMem(m, 1, false)
		case lang.Decr:
			e.SetMem(m, -1, false)
		case lang.Lurd:
			if d.dx != 0 {
				d = Dir{0, d.dx}
			} else {
				d = Dir{d.dy, 0}
			}
		case lang.Ruld:
			if d.dx != 0 {
				d = Dir{0, -d.dx}
			} else {
				d = Dir{-d.dy, 0}
			}
		case lang.Leave:
			if e.modular {
				if cstack.Len() == 0 {
					if e.debug {
						log.Printf("leave at (%d,%d)", p.x, p.y)
					}
					return
				}
				el := cstack.Back()
				st := el.Value.(State)
				p = st.pos
				d = st.dir
				if e.twist {
					p.x += d.dx
					p.y += d.dy
				}
				cstack.Remove(el)
			}
		case lang.Enter:
			if e.modular {
				cstack.PushBack(State{Pos{p.x, p.y}, d})
				if !e.twist {
					p.x += d.dx
					p.y += d.dy
				}
			}
		case lang.Skip:
			p.x += d.dx
			p.y += d.dy
		case lang.SkipZ:
			if e.GetMem(m) == 0 {
				p.x += d.dx
				p.y += d.dy
			}
		case lang.Split:
			if e.bloated {
				e.sg.Add(1)
				go e.Interpret(p, d, m)
				p.x += d.dx
				p.y += d.dy
			}
		case lang.Rand:
			if e.bloated {
				e.SetMem(m, rand.Intn(256), true)
			}
		case lang.Write:
			rb := e.GetMem(m)
			fmt.Printf("%c", rb)
		case lang.Read:
			b_in, b_err := stdin.ReadByte()
			if b_err == io.EOF {
				return
			} else if b_err == nil {
				if e.debug {
					log.Printf("read: %d", int(b_in))
				}
				e.SetMem(m, int(b_in&0xFF), true)
			}
		}
		p.x += d.dx
		p.y += d.dy
	}
	if e.debug {
		log.Printf("exit at (%d,%d)", p.x, p.y)
	}
}

func (e *Snusp) Run() {
	if !e.loaded {
		return
	}
	e.sg.Add(1)
	go e.Interpret(e.pos, Dir{1, 0}, Pos{0, 0})
	e.sg.Wait()
}

func main() {
	snusp := new(Snusp)
	flag.BoolVar(&snusp.debug, "debug", false, "debug")
	flag.BoolVar(&snusp.modular, "modular", true, "modular SNUSP")
	flag.BoolVar(&snusp.bloated, "bloated", false, "bloated SNUSP")
	flag.BoolVar(&snusp.twist, "twist", false, "modular SNUSP flavour: twist")
	flag.Parse()
	if flag.NArg() > 0 {
		snusp.Load(flag.Arg(0))
		snusp.Run()
		if snusp.debug {
			for k, v := range snusp.mem {
				log.Print(k, v)
			}
		}
	} else {
		fmt.Printf("%s [flags] FILENAME\nFlags and default values are:\n", os.Args[0])
		flag.PrintDefaults()
	}
}
