package cursor

import "fmt"

/*
  - Position the Cursor:
    \033[<L>;<C>H
    Or
    \033[<L>;<C>f
    puts the cursor at line L and column C.
  - Move the cursor up N lines:
    \033[<N>A
  - Move the cursor down N lines:
    \033[<N>B
  - Move the cursor forward N columns:
    \033[<N>C
  - Move the cursor backward N columns:
    \033[<N>D

  - Clear the screen, move to (0,0):
    \033[2J
  - Erase to end of line:
    \033[K
*/
const (
	pos        = "\033[%d;%dH"
	up         = "\033[%dA"
	down       = "\033[%dB"
	right      = "\033[%dC"
	left       = "\033[%dD"
	clear      = "\033[2J"
	eraseToEOL = "\033[K"
)

func Position(x int, y int) {
	fmt.Printf(pos, y+1, x+1)
}

func Up(n int) {
	fmt.Printf(up, n)
}

func Down(n int) {
	fmt.Printf(down, n)
}

func Right(n int) {
	fmt.Printf(right, n)
}

func Left(n int) {
	fmt.Printf(left, n)
}

func Clear() {
	fmt.Printf(clear)
}

func EraseToEOL() {
	fmt.Printf(eraseToEOL)
}
