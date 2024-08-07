package geometry

// code of this file is a modified version of code from
// https://github.com/anaseto/gruid, which has the following license:
//
// Copyright (c) 2020 Yon <anaseto@bardinflor.perso.aquilenet.fr>
//
// Permission to use, copy, modify, and distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

// CCMapAll computes a map of the connected components. It makes the
// assumption that the paths are bidirectional, allowing for efficient
// computation. This means, in particular, that the pather should return no
// neighbors for obstacles.
func (pr *PathRange) CCMapAll(nb Pather) {
	max := pr.Rg.Size()
	w, h := max.X, max.Y
	if pr.CC == nil {
		pr.CC = make([]int, w*h)
	} else {
		for i := range pr.CC {
			pr.CC[i] = 0
		}
	}
	pr.CCStack = pr.CCStack[:0]
	ccid := 1
	for i := 0; i < len(pr.CC); i++ {
		if pr.CC[i] > 0 {
			continue
		}
		pr.CC[i] = ccid
		pr.CCStack = append(pr.CCStack, i)
		for len(pr.CCStack) > 0 {
			idx := pr.CCStack[len(pr.CCStack)-1]
			pr.CCStack = pr.CCStack[:len(pr.CCStack)-1]
			p := idxToPos(idx, w)
			for _, q := range nb.Neighbors(p) {
				if !q.In(pr.Rg) {
					continue
				}
				nidx := pr.idx(q)
				if pr.CC[nidx] > 0 {
					continue
				}
				pr.CC[nidx] = ccid
				pr.CCStack = append(pr.CCStack, nidx)
			}
		}
		ccid++
	}
}

// CCMap computes the connected component which contains a given position.
// It returns a cached slice with all the positions in the same connected
// component as p, or nil if p is out of range.  It makes the assumption that
// the paths are bidirectional, allowing for efficient computation. This means,
// in particular, that the pather should return no neighbors for obstacles.
//
// It makes uses of the same caching structures as ComputeCCAll, so CCAt will
// return -1 on all unreachable positions from p.
func (pr *PathRange) CCMap(nb Pather, p Point) []Point {
	max := pr.Rg.Size()
	w, h := max.X, max.Y
	if pr.CC == nil {
		pr.CC = make([]int, w*h)
	} else {
		for i := range pr.CC {
			pr.CC[i] = 0
		}
	}
	if pr.CCIterCache == nil {
		pr.CCIterCache = make([]Point, w*h)
	}
	pr.CCIterCache = pr.CCIterCache[:0]
	pr.CCStack = pr.CCStack[:0]
	if !p.In(pr.Rg) {
		return nil
	}
	idx := pr.idx(p)
	ccid := 1
	pr.CC[idx] = ccid
	pr.CCStack = append(pr.CCStack, idx)
	for len(pr.CCStack) > 0 {
		idx = pr.CCStack[len(pr.CCStack)-1]
		pr.CCStack = pr.CCStack[:len(pr.CCStack)-1]
		p := idxToPos(idx, w)
		pr.CCIterCache = append(pr.CCIterCache, p)
		for _, q := range nb.Neighbors(p) {
			if !q.In(pr.Rg) {
				continue
			}
			nidx := pr.idx(q)
			if pr.CC[nidx] > 0 {
				continue
			}
			pr.CC[nidx] = ccid
			pr.CCStack = append(pr.CCStack, nidx)
		}
	}
	return pr.CCIterCache
}

// CCMapAt returns a positive number identifying the position's connected
// component as computed by either the last CCMap or CCMapAll call. It
// returns -1 on out of range positions.
func (pr *PathRange) CCMapAt(p Point) int {
	if !p.In(pr.Rg) || pr.CC == nil {
		return -1
	}
	return pr.CC[pr.idx(p)] - 1
}

// idxToPos returns a grid position given an index and the width of the grid.
func idxToPos(i, w int) Point {
	return Point{X: i % w, Y: i / w}
}
