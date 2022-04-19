package main

import "github.com/gdamore/tcell"

type Border struct {
	t  int // Border thickness in characters
	x1 int // Left border boundary x-coordinate
	x2 int // Right border boundary x-coordinate
	y1 int // Top border boundary y-coordinate
	y2 int // Bottom border boundary y-coordinate
}

func (b *Border) Draw(s tcell.Screen) {
	for c := b.x1; c <= b.x2; c++ { // Add top and bottom borders
		s.SetContent(c, b.y1, '#', nil, tcell.StyleDefault)
		s.SetContent(c, b.y2, '#', nil, tcell.StyleDefault)
	}

	for r := b.y1 + 1; r <= b.y2-1; r++ { // Add left and right borders
		s.SetContent(b.x1, r, '#', nil, tcell.StyleDefault)
		s.SetContent(b.x2, r, '#', nil, tcell.StyleDefault)
	}
}

type Tree struct {
	x int // Trunk left corner x-coordinate
	y int // Trunk left corner y-coordinate
}

func (t *Tree) Draw(s tcell.Screen) {
	s.SetContent(t.x, t.y, '|', nil, tcell.StyleDefault)
	s.SetContent(t.x+1, t.y, '|', nil, tcell.StyleDefault)
	s.SetContent(t.x-1, t.y-1, '/', nil, tcell.StyleDefault)
	s.SetContent(t.x, t.y-1, '_', nil, tcell.StyleDefault)
	s.SetContent(t.x+1, t.y-1, '_', nil, tcell.StyleDefault)
	s.SetContent(t.x+2, t.y-1, '\\', nil, tcell.StyleDefault)
	s.SetContent(t.x, t.y-2, '/', nil, tcell.StyleDefault)
	s.SetContent(t.x+1, t.y-2, '\\', nil, tcell.StyleDefault)
}

type Player struct {
	x int // Player x-coordinate
	y int // Player y-coordinate
}

func (p *Player) Clear(s tcell.Screen) {
	s.SetContent(p.x, p.y, ' ', nil, tcell.StyleDefault)
	s.Show()
}

func (p *Player) Draw(s tcell.Screen) {
	s.SetContent(p.x, p.y, '@', nil, tcell.StyleDefault)
	s.Show()
}

func (p *Player) Move(s tcell.Screen, len int, dir int) {
	p.Clear(s)

	// Determine new location.
	if len != 0 {
		switch dir {
		case 0: // up
			p.y = p.y - 1
		case 1: // right
			p.x = p.x + 1
		case 2: // down
			p.y = p.y + 1
		case 3: // left
			p.x = p.x - 1
		}
	}

	// Ensure player doesn't move past boundaries.
	w, h := s.Size()
	if p.x >= w-1 {
		p.x = w - 2
	} else if p.x <= 0 {
		p.x = 1
	}

	if p.y >= h-1 {
		p.y = h - 2
	} else if p.y <= 1 {
		p.y = 1
	}

	p.Clear(s)
	p.Draw(s)

}
