/*
   gosnusp - a SNUSP esolang interpreter written in Go
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
	. "./snusp"
	"./snusp/dir"
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

type State struct {
	pos Pos
	dir dir.Dir
}

type Snusp struct {
	debug   bool
	loaded  bool
	modular bool
	bloated bool
	twist   bool
	eof0    bool

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
	if p.X < 0 || p.Y < 0 {
		return NoOp
	}
	e.codeLock.RLock()
	if s, ok := e.code[p.Y]; ok {
		if p.X >= len(s) {
			res = NoOp
		} else {
			res = s[p.X : p.X+1]
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
	e.pos.X = 0
	e.pos.Y = 0
	e.size.W = 0
	e.size.H = 0
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
			e.size.W = Max(e.size.W, len(s))
			if px := strings.Index(s, Start); px != -1 {
				e.pos.Y = ly
				e.pos.X = px
				if e.debug {
					log.Printf("Start at (%d, %d)", e.pos.X, e.pos.Y)
				}
			}
			ly++
		} else if err == io.EOF {
			e.code[ly] = s
			ly++
			e.loaded = true
			e.size.H = ly
			e.size.W = Max(e.size.W, len(s))
			if e.debug {
				log.Printf("w=%d, h=%d", e.size.W, e.size.H)
			}
			break
		} else {
			log.Fatal(err)
			break //?
		}
	}
}

func (e *Snusp) Interpret(p Pos, d dir.Dir, m Pos) {
	if e.debug {
		log.Printf("interpret at (%d,%d), dir (%d,%d)", p.X, p.Y, d.Dx, d.Dy)
	}
	defer e.sg.Done()
	stdin := bufio.NewReader(os.Stdin)
	stdout := bufio.NewWriter(os.Stdout)
	defer stdout.Flush()
	cstack := list.New()
	for p.X >= 0 && p.X < e.size.W && p.Y >= 0 && p.Y < e.size.H {
		c := e.Get(p)
		if e.debug {
			log.Print(p, c, d)
		}
		switch c {
		case Left:
			m.X--
		case Right:
			m.X++
		case Up:
			if e.bloated {
				m.Y--
			}
		case Down:
			if e.bloated {
				m.Y++
			}
		case Incr:
			e.SetMem(m, 1, false)
		case Decr:
			e.SetMem(m, -1, false)
		case Lurd:
			// it was a lot better when it was done with the if...
			// but I wanted to try init() in a package, and so...
			d = LurdMap[d]
		case Ruld:
			d = RuldMap[d]
		case Leave:
			if e.modular {
				if cstack.Len() == 0 {
					if e.debug {
						log.Printf("leave at (%d,%d)", p.X, p.Y)
					}
					return
				}
				el := cstack.Back()
				st := el.Value.(State)
				p = st.pos
				d = st.dir
				if e.twist {
					p.X += d.Dx
					p.Y += d.Dy
				}
				cstack.Remove(el)
			}
		case Enter:
			if e.modular {
				cstack.PushBack(State{Pos{p.X, p.Y}, d})
				if !e.twist {
					p.X += d.Dx
					p.Y += d.Dy
				}
			}
		case Skip:
			p.X += d.Dx
			p.Y += d.Dy
		case SkipZ:
			if e.GetMem(m) == 0 {
				p.X += d.Dx
				p.Y += d.Dy
			}
		case Split:
			if e.bloated {
				e.sg.Add(1)
				p.X += d.Dx
				p.Y += d.Dy
				go e.Interpret(p, d, m)
			}
		case Rand:
			if e.bloated {
				e.SetMem(m, rand.Intn(256), true)
			}
		case Write:
			rb := e.GetMem(m)
			stdout.WriteByte(rb)
		case Read:
			stdout.Flush() // flush stdout before reading...
			b_in, b_err := stdin.ReadByte()
			if b_err == io.EOF {
				if e.eof0 {
					if e.debug {
						log.Print("EOF write 0 in cell")
					}
					e.SetMem(m, 0, true)
				} else {
					return
				}
			} else if b_err == nil {
				if e.debug {
					log.Printf("read: %d", int(b_in))
				}
				e.SetMem(m, int(b_in&0xFF), true)
			}
		}
		p.X += d.Dx
		p.Y += d.Dy
	}
	if e.debug {
		log.Printf("exit at (%d,%d)", p.X, p.Y)
	}
}

func (e *Snusp) Run() {
	if !e.loaded {
		return
	}
	e.sg.Add(1)
	go e.Interpret(e.pos, dir.Dir{1, 0}, Pos{0, 0})
	e.sg.Wait()
	if e.debug {
		log.Print("memory dump")
		for k, v := range e.mem {
			log.Print(k, v)
		}
	}
}

func main() {
	snusp := new(Snusp)
	flag.BoolVar(&snusp.debug, "debug", false, "debug")
	flag.BoolVar(&snusp.modular, "modular", true, "modular SNUSP")
	flag.BoolVar(&snusp.bloated, "bloated", false, "bloated SNUSP")
	flag.BoolVar(&snusp.twist, "twist", true, "modular SNUSP flavour: twist")
	flag.BoolVar(&snusp.eof0, "eof0", false, "EOF place 0 in the cell")
	flag.Parse()
	if flag.NArg() > 0 {
		snusp.Load(flag.Arg(0))
		snusp.Run()
	} else {
		fmt.Printf("%s [flags] FILENAME\nFlags and default values are:\n", os.Args[0])
		flag.PrintDefaults()
	}
}
