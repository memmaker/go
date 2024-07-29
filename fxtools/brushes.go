package fxtools

import "github.com/memmaker/go/geometry"

type Brush interface {
    StartDrawing(pos geometry.Point)
    DraggedOver(pos geometry.Point) []geometry.Point
    StopDrawing(pos geometry.Point) []geometry.Point
    Icon() rune
    Name() string
}
type PointSet map[geometry.Point]bool

func NewPointSet() *PointSet {
    set := PointSet{}
    return &set
}

func (s *PointSet) Add(pos geometry.Point) {
    (*s)[pos] = true
}

func (s *PointSet) ToSlice() []geometry.Point {
    slice := make([]geometry.Point, 0, len(*s))
    for pos := range *s {
        slice = append(slice, pos)
    }
    return slice
}

func (s *PointSet) Clear() {
    clear(*s)
}

func (s *PointSet) Cardinality() int {
    return len(*s)
}

func (s *PointSet) Pop() geometry.Point {
    for pos := range *s {
        delete(*s, pos)
        return pos
    }
    return geometry.Point{}
}

func (s *PointSet) Contains(neighbor geometry.Point) bool {
    _, ok := (*s)[neighbor]
    return ok
}
func NewPencil() *Pencil {
    return &Pencil{
        drawPositions: NewPointSet(),
    }
}

type Pencil struct {
    drawPositions *PointSet
}

func (p *Pencil) Name() string {
    return "pencil"
}

func (p *Pencil) StartDrawing(pos geometry.Point) {
    p.drawPositions.Add(pos)
}

func (p *Pencil) Icon() rune {
    return 'p'
}

func (p *Pencil) DraggedOver(pos geometry.Point) []geometry.Point {
    p.drawPositions.Add(pos)
    return p.drawPositions.ToSlice()
}

func (p *Pencil) StopDrawing(pos geometry.Point) []geometry.Point {
    p.drawPositions.Add(pos)
    slice := p.drawPositions.ToSlice()
    p.drawPositions.Clear()
    return slice
}

func NewLineBrush() *LineBrush {
    return &LineBrush{}
}

type LineBrush struct {
    startPos geometry.Point
}

func (p *LineBrush) Name() string {
    return "line"
}

func (p *LineBrush) StartDrawing(pos geometry.Point) {
    p.startPos = pos
}

func (p *LineBrush) Icon() rune {
    return 'L'
}

func (p *LineBrush) DraggedOver(pos geometry.Point) []geometry.Point {
    drawPositions := make([]geometry.Point, 0)
    for _, linePos := range geometry.LineOfSight(p.startPos, pos, func(p geometry.Point) bool {
        return true
    }) {
        drawPositions = append(drawPositions, linePos)
    }
    return drawPositions
}

func (p *LineBrush) StopDrawing(pos geometry.Point) []geometry.Point {
    return p.DraggedOver(pos)
}

func NewOutlinedRectangleBrush() *RectangleBrush {
    return &RectangleBrush{
        drawRect: geometry.Rect{},
        fill:     false,
    }
}
func NewFilledRectangleBrush() *RectangleBrush {
    return &RectangleBrush{
        drawRect: geometry.Rect{},
        fill:     true,
    }
}

type RectangleBrush struct {
    drawRect geometry.Rect
    startPos geometry.Point
    fill     bool
}

func (r *RectangleBrush) Name() string {
    return "rect"
}

func (r *RectangleBrush) Icon() rune {
    if r.fill {
        return 'F'
    }
    return 'f'
}

func (r *RectangleBrush) StartDrawing(pos geometry.Point) {
    r.startPos = pos
    r.drawRect = geometry.NewRect(pos.X, pos.Y, pos.X, pos.Y)
}

func (r *RectangleBrush) DraggedOver(pos geometry.Point) []geometry.Point {
    r.drawRect = geometry.NewRect(r.startPos.X, r.startPos.Y, pos.X, pos.Y)
    if r.fill {
        return r.filledRect()
    }
    return r.outlinedRect()
}

func (r *RectangleBrush) StopDrawing(pos geometry.Point) []geometry.Point {
    r.drawRect = geometry.NewRect(r.startPos.X, r.startPos.Y, pos.X, pos.Y)
    if r.fill {
        return r.filledRect()
    }
    return r.outlinedRect()
}
func (r *RectangleBrush) filledRect() []geometry.Point {
    var filled []geometry.Point
    for x := r.drawRect.Min.X; x <= r.drawRect.Max.X; x++ {
        for y := r.drawRect.Min.Y; y <= r.drawRect.Max.Y; y++ {
            filled = append(filled, geometry.Point{X: x, Y: y})
        }
    }
    return filled
}

func (r *RectangleBrush) outlinedRect() []geometry.Point {
    outine := NewPointSet()
    for x := r.drawRect.Min.X; x <= r.drawRect.Max.X; x++ {
        outine.Add(geometry.Point{X: x, Y: r.drawRect.Min.Y})
        outine.Add(geometry.Point{X: x, Y: r.drawRect.Max.Y})
    }
    for y := r.drawRect.Min.Y; y <= r.drawRect.Max.Y; y++ {
        outine.Add(geometry.Point{X: r.drawRect.Min.X, Y: y})
        outine.Add(geometry.Point{X: r.drawRect.Max.X, Y: y})
    }
    return outine.ToSlice()
}

func NewOutlinedCircleBrush() *CircleBrush {
    return &CircleBrush{
        drawRect: geometry.Rect{},
        fill:     false,
        square:   true,
    }
}
func NewFilledCircleBrush() *CircleBrush {
    return &CircleBrush{
        drawRect: geometry.Rect{},
        fill:     true,
        square:   true,
    }
}

type CircleBrush struct {
    drawRect geometry.Rect
    startPos geometry.Point
    fill     bool
    square   bool
}

func (r *CircleBrush) Name() string {
    return "circle"
}

func (r *CircleBrush) Icon() rune {
    if r.fill {
        return 'C'
    }
    return 'c'
}

func (r *CircleBrush) StartDrawing(pos geometry.Point) {
    r.startPos = pos
    r.drawRect = geometry.NewRect(pos.X, pos.Y, pos.X, pos.Y)
}

func (r *CircleBrush) DraggedOver(pos geometry.Point) []geometry.Point {
    r.drawRect = geometry.NewRect(r.startPos.X, r.startPos.Y, pos.X, pos.Y)
    if r.fill {
        ellipse := r.outlinedEllipse()
        if !ellipse.Contains(r.drawRect.Mid()) {
            for _, fillPos := range floodFillFrom(r.drawRect.Mid(), func(src, from, to geometry.Point) bool {
                return !ellipse.Contains(to)
            }) {
                ellipse.Add(fillPos)
            }
        }
        return ellipse.ToSlice()
    }
    return r.outlinedEllipse().ToSlice()
}

func (r *CircleBrush) StopDrawing(pos geometry.Point) []geometry.Point {
    return r.DraggedOver(pos)
}
func (r *CircleBrush) filledCircle() []geometry.Point {
    var filled []geometry.Point
    // we want to draw a circle inside the bounds of the rectangle
    shortestSide := IntMin(r.drawRect.Size().X, r.drawRect.Size().Y) - 1
    radius := shortestSide / 2
    for x := r.drawRect.Min.X + 1; x < r.drawRect.Max.X; x++ {
        for y := r.drawRect.Min.Y + 1; y < r.drawRect.Max.Y; y++ {
            if geometry.Distance(geometry.Point{X: x, Y: y}, r.drawRect.Mid())+0.5 < float64(radius) {
                filled = append(filled, geometry.Point{X: x, Y: y})
            }
        }
    }
    return filled
}
func (r *CircleBrush) outlinedEllipse() *PointSet {
    result := NewPointSet()
    x0 := r.drawRect.Min.X
    y0 := r.drawRect.Min.Y
    x1 := r.drawRect.Max.X
    y1 := r.drawRect.Max.Y
    a := geometry.Abs(x1 - x0)
    b := geometry.Abs(y1 - y0)
    b1 := b & 1 /* values of diameter */
    dx := 4 * (1 - a) * b * b
    dy := 4 * (b1 + 1) * a * a /* error increment */

    err := dx + dy + b1*a*a /* error of 1.step */
    e2 := 0                 /* error of 2.step */

    if x0 > x1 {
        x0 = x1
        x1 += a
    } /* if called with swapped points */
    if y0 > y1 {
        y0 = y1
    } /* .. exchange them */
    y0 += (b + 1) / 2
    y1 = y0 - b1 /* starting pixel */
    a *= 8 * a
    b1 = 8 * b * b

    for ok := true; ok; ok = x0 <= x1 {
        result.Add(geometry.Point{X: x1, Y: y0}) /*   I. Quadrant */
        result.Add(geometry.Point{X: x0, Y: y0}) /*  II. Quadrant */
        result.Add(geometry.Point{X: x0, Y: y1}) /* III. Quadrant */
        result.Add(geometry.Point{X: x1, Y: y1}) /*  IV. Quadrant */
        e2 = 2 * err
        if e2 <= dy {
            y0++
            y1--
            dy += a
            err += dy
        } /* y step */
        if e2 >= dx || 2*err > dy {
            x0++
            x1--
            dx += b1
            err += dx
        } /* x step */
    }

    for y0-y1 < b { /* too early stop of flat ellipses a=1 */
        result.Add(geometry.Point{X: x0 - 1, Y: y0}) /* -> finish tip of ellipse */
        result.Add(geometry.Point{X: x1 + 1, Y: y0})
        y0++
        result.Add(geometry.Point{X: x0 - 1, Y: y1})
        result.Add(geometry.Point{X: x1 + 1, Y: y1})
        y1--
    }
    return result
}
func (r *CircleBrush) outlinedCircle() []geometry.Point {
    var outline []geometry.Point
    shortestSide := IntMin(r.drawRect.Size().X, r.drawRect.Size().Y) - 1
    radius := shortestSide / 2
    for x := r.drawRect.Min.X + 1; x < r.drawRect.Max.X; x++ {
        for y := r.drawRect.Min.Y + 1; y < r.drawRect.Max.Y; y++ {
            if int(geometry.Distance(geometry.Point{X: x, Y: y}, r.drawRect.Mid())+0.5) == radius {
                outline = append(outline, geometry.Point{X: x, Y: y})
            }
        }
    }
    return outline
}

func NewPaintBucket(traversable func(geometry.Point, geometry.Point, geometry.Point) bool) *PaintBucket {
    return &PaintBucket{
        traversable: traversable,
    }
}

type PaintBucket struct {
    traversable func(geometry.Point, geometry.Point, geometry.Point) bool
}

func (p *PaintBucket) Name() string {
    return "fill"
}

func (p *PaintBucket) StartDrawing(pos geometry.Point) {}

func (p *PaintBucket) Icon() rune {
    return 'B'
}

func (p *PaintBucket) DraggedOver(pos geometry.Point) []geometry.Point {
    return floodFillFrom(pos, p.traversable)
}

func (p *PaintBucket) StopDrawing(pos geometry.Point) []geometry.Point {
    return p.DraggedOver(pos)
}

func IntMin(one int, two int) int {
    if one < two {
        return one
    }
    return two
}

func floodFillFrom(sourceOfFill geometry.Point, canTraverse func(geometry.Point, geometry.Point, geometry.Point) bool) []geometry.Point {
    closedSet := NewPointSet()
    openSet := NewPointSet()
    openSet.Add(sourceOfFill)
    neighbors := geometry.Neighbors{}
    for openSet.Cardinality() > 0 {
        current := openSet.Pop()
        closedSet.Add(current)
        for _, neighbor := range neighbors.Cardinal(current, func(potentialNeighbor geometry.Point) bool {
            return canTraverse(sourceOfFill, current, potentialNeighbor)
        }) {
            if !closedSet.Contains(neighbor) {
                openSet.Add(neighbor)
            }
        }
    }

    return closedSet.ToSlice()
}
