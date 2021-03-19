package main

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

type field rune

const (
	empty   field = ' '
	player1 field = 'o'
	player2 field = 'x'

	width    int = 4
	height   int = 4
	winCount int = 4
)

type won field

type coord struct{ x, y int }

type direc struct{ x, y int }

type gameBoard [width][height]field

type model struct {
	board         gameBoard
	currentPlayer field
	emptyFields   int
	info          string
}

func main() {
	game := NewModel()
	p := tea.NewProgram(game)
	p.Start()

}

func NewModel() tea.Model {
	m := model{}
	m.emptyFields = height * width
	m.currentPlayer = player1
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			m.board[x][y] = empty
		}
	}
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// is a key press
	case tea.KeyMsg:
		k := msg.String()
		switch k {
		// should we stop the programm?
		case "ctrl+c", "q":
			return m, tea.Quit
		case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0":
			k, err := strconv.Atoi(k)
			if err != nil {
				return m, func() tea.Msg { return err }
			}
			cmd := m.handleTurn(k)
			return m, cmd
		default:
			return m, func() tea.Msg {
				return fmt.Errorf("wrong input: %q (input is not a valid key, use 'q' to quit or a column number to set a stone)", k)
			}
		}
	case error:
		m.info = msg.Error()
		return m, nil
	case won:
		m.info = fmt.Sprintf("player with stone: %q", msg)
		return m, tea.Quit
	}
	return m, nil
}

func (m model) View() string {
	var boardString string
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			boardString += string(m.board[x][y])
		}
		boardString += "\n"
	}

	return boardString + "\n" + m.info + "\n"
}

func (m *model) handleTurn(userInput int) tea.Cmd {
	// handle full board
	if m.emptyFields <= 0 {
		return func() tea.Msg { return fmt.Errorf("no empty fields left - game ends with a draw") }
	}

	// validate input and return appropriate error message if necessary
	column := userInput - 1
	if column < 0 || column > width-1 {
		return func() tea.Msg { return fmt.Errorf("wrong input: %d (column out of board)", userInput) }
	}

	// switch currentPlayer
	tempPlayer := player1
	if m.currentPlayer == player1 {
		tempPlayer = player2
	}

	// set stone
	var lastPosition coord
	for y := height - 1; y >= 0; y-- {
		if m.board[column][y] == empty {
			m.board[column][y] = tempPlayer
			lastPosition = coord{column, y}
			break
		}
		if y == 0 {
			return func() tea.Msg { return fmt.Errorf("column %d already full", userInput) }
		}
	}
	m.currentPlayer = tempPlayer
	m.emptyFields -= 1

	// check victory condition
	directions := [4]direc{{0, 1}, {1, 0}, {1, -1}, {-1, -1}}
	for i := 0; i < len(directions); i++ {
		if m.board.checkVictory(lastPosition, m.currentPlayer, directions[i]) {
			return func() tea.Msg { return won(m.currentPlayer) }
		}
	}
	return nil
}

func (b *gameBoard) checkVictory(lastPosition coord, currentPlayer field, direction direc) bool {
	// count the start stone
	stoneCount := 1
	// count in one direction
	for i := 1; i <= winCount; i++ {
		var position coord
		position.x = lastPosition.x + direction.x*i
		position.y = lastPosition.y + direction.y*i
		if !b.validPosition(position) {
			break
		}
		if !b.sameStone(position, currentPlayer) {
			break
		}
		stoneCount++
	}
	// and the other direction:
	direction.x *= -1
	direction.y *= -1
	for i := 1; i <= winCount; i++ {
		var position coord
		position.x = lastPosition.x + direction.x*i
		position.y = lastPosition.y + direction.y*i
		if !b.validPosition(position) {
			break
		}
		if !b.sameStone(position, currentPlayer) {
			break
		}
		stoneCount++
	}

	return stoneCount >= winCount
}

func (b *gameBoard) validPosition(position coord) bool {
	if position.x > width-1 || position.x < 0 || position.y > height-1 || position.y < 0 {
		return false
	}
	return true
}

func (b *gameBoard) sameStone(position coord, lastStone field) bool {
	return b[position.x][position.y] == lastStone
}
