// Copyright 2017 Ritchie Borja
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package minesweeper

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	SAMPLE_GRID_WIDTH  = 10
	SAMPLE_GRID_HEIGHT = 40
)

func newBlankGame() Minesweeper {
	game, _ := NewGame()
	return game
}

func newSampleGame() Minesweeper {
	game, _ := NewGame(Grid{SAMPLE_GRID_WIDTH, SAMPLE_GRID_HEIGHT})
	return game
}

func TestGridMustNotBeSquaredForTheSakeOfTesting(t *testing.T) {
	assert.True(t, SAMPLE_GRID_WIDTH != SAMPLE_GRID_HEIGHT)
}

func TestGame_SetGrid(t *testing.T) {
	minesweeper := newBlankGame()
	minesweeper.SetGrid(SAMPLE_GRID_WIDTH, SAMPLE_GRID_HEIGHT)
	assert.Equal(t, minesweeper.(*game).Board.Grid, &Grid{SAMPLE_GRID_WIDTH, SAMPLE_GRID_HEIGHT})
}

func TestGameWithGridArgument(t *testing.T) {
	minesweeper := newSampleGame()
	assert.Equal(t, minesweeper.(*game).Board.Grid, &Grid{SAMPLE_GRID_WIDTH, SAMPLE_GRID_HEIGHT})
}

func TestNewGridWhenStartedGame(t *testing.T) {
	minesweeper := newSampleGame()
	err := minesweeper.SetGrid(10, 20)
	assert.Error(t, err)
	assert.NotNil(t, err, "Must report an error upon setting a new grid from an already started game")
	assert.IsType(t, new(GameAlreadyStarted), err, "The error must be GameAlreadyStarted error type")
}

func TestFlaggedBlock(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.Flag(3, 6)
	assert.Equal(t, minesweeper.(*game).Blocks[3][6].flagged, true)
}

func TestGame_SetDifficulty(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(EASY)
	assert.Equal(t, minesweeper.(*game).Difficulty, EASY)
}

func TestShiftFromMaxPosition(t *testing.T) {
	grid := Grid{5, 5}
	x, y := shiftPosition(&grid, 4, 4)
	assert.Equal(t, struct {
		x int
		y int
	}{0, 0}, struct {
		x int
		y int
	}{x, y})
}

func TestBombsInPlace(t *testing.T) {

	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(EASY)
	minesweeper.Play()

	game := minesweeper.(*game)

	numOfBombs := int(float32(game.Width*game.Height) * EASY_MULTIPLIER)
	countedBombs := 0
	for _, row := range game.Blocks {
		for _, block := range row {
			if block.Node == BOMB {
				countedBombs++
			}
		}
	}
	assert.Equal(t, numOfBombs, countedBombs)
}

func TestTalliedBomb(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(EASY)
	minesweeper.Play()

	game := minesweeper.(*game)
	width := game.Width
	height := game.Height

	count := func(blocks Blocks, x, y int) (has int) {
		if x >= 0 && y >= 0 &&
			x < width && y < height &&
			blocks[x][y].Node&BOMB == 1 {
			return 1
		}
		return
	}

	hasSurroundingTally := func(blocks Blocks, x, y int) int {
		if x >= 0 && y >= 0 &&
			x < width && y < height {
			switch blocks[x][y].Node {
			case NUMBER:
				return 1
			case BOMB:
				return -1
			default:
				return 0
			}
		}
		return -1
	}
	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node == BOMB {
				assert.NotEqual(t, 0, hasSurroundingTally(game.Blocks, x-1, y-1))
				assert.NotEqual(t, 0, hasSurroundingTally(game.Blocks, x-1, y))
				assert.NotEqual(t, 0, hasSurroundingTally(game.Blocks, x-1, y+1))
				assert.NotEqual(t, 0, hasSurroundingTally(game.Blocks, x, y-1))
				assert.NotEqual(t, 0, hasSurroundingTally(game.Blocks, x, y+1))
				assert.NotEqual(t, 0, hasSurroundingTally(game.Blocks, x+1, y-1))
				assert.NotEqual(t, 0, hasSurroundingTally(game.Blocks, x+1, y))
				assert.NotEqual(t, 0, hasSurroundingTally(game.Blocks, x+1, y+1))
			}
		}
	}

	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node == NUMBER {
				var counted int
				counted = count(game.Blocks, x-1, y-1) +
					count(game.Blocks, x-1, y) +
					count(game.Blocks, x-1, y+1) +
					count(game.Blocks, x, y-1) +
					count(game.Blocks, x, y+1) +
					count(game.Blocks, x+1, y-1) +
					count(game.Blocks, x+1, y) +
					count(game.Blocks, x+1, y+1)
				assert.Equal(t, counted, block.Value)
			}
		}
	}
}

func TestVisitedBlock(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(EASY)
	minesweeper.Play()

	game := minesweeper.(*game)

	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node == NUMBER {
				game.Visit(x, y)
				assert.True(t, game.Blocks[x][y].visited)
			}
		}
	}

}

func TestVisitedBombToGameOver(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(EASY)
	minesweeper.Play()

	game := minesweeper.(*game)
	var x, y int
	var err error

	for i, row := range game.Blocks {
		for j, block := range row {
			if block.Node == BOMB {
				x, y = i, j
				_, err = game.Visit(x, y)
				assert.Error(t, err)
				assert.NotNil(t, err)
				assert.IsType(t, new(Exploded), err)
			}
		}
	}

}

func TestVisitedBombToGameOverWithCorrectLocationReason(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(EASY)
	minesweeper.Play()

	game := minesweeper.(*game)
	var x, y int
	var err error

	for i, row := range game.Blocks {
		for j, block := range row {
			if block.Node == BOMB {
				x, y = i, j
				_, err = game.Visit(x, y)
				assert.Error(t, err)
				assert.EqualError(t, err,
					fmt.Sprintf("Game over at X=%v Y=%v",
						x, y))
				assert.IsType(t, new(Exploded), err)
			}
		}
	}

}

func TestVisitedUnmarkedBlockDistributeVisit(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(EASY)
	minesweeper.Play()

	game := minesweeper.(*game)

	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node == UNKNOWN && !block.visited {
				minesweeper.Visit(x, y)
			}
		}
	}

	for _, row := range game.Blocks {
		for _, block := range row {
			if block.Node == UNKNOWN {
				assert.True(t, block.visited)
			}
		}
	}
}

func TestVisitAFlaggedBlock(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(EASY)
	minesweeper.Play()

	game := minesweeper.(*game)

	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node == BOMB {
				minesweeper.Flag(x, y)
				_, err := minesweeper.Visit(x, y)
				assert.NoError(t, err)
				if err != nil {
					assert.IsType(t, new(Exploded), err)
				}
			}
		}
	}
}

func TestVisitedBlocksReturnOneBlockWhenAHintBlockIsVisited(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(EASY)
	minesweeper.Play()

	game := minesweeper.(*game)

	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node == NUMBER {
				visitedBlocks, err := minesweeper.Visit(x, y)
				assert.NoError(t, err)
				assert.Equal(t, 1, len(visitedBlocks))
				assert.Equal(t, block.Value, visitedBlocks[0].Value)
				assert.Equal(t, visitedBlocks[0], game.Blocks[x][y])
			}
		}
	}
}

func TestVisitedBlocksWhenBlockIsABomb(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(EASY)
	minesweeper.Play()

	game := minesweeper.(*game)

	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node == BOMB {
				visitedBlocks, err := minesweeper.Visit(x, y)
				assert.Error(t, err)
				assert.EqualError(t, err, (&Exploded{struct{ x, y int }{x: x, y: y}}).Error())
				assert.Equal(t, visitedBlocks[0], game.Blocks[x][y])
			}
		}
	}
}

func TestVisitedBlockWhenBlockIsUnknownAndSpreadVisits(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(EASY)
	minesweeper.Play()

	game := minesweeper.(*game)

	var x, y int
	var actualVisitedBlocks []Block
first:
	for i, row := range game.Blocks {
		for j, block := range row {
			if block.Node == UNKNOWN && !block.visited {
				x, y = i, j
				var err error
				actualVisitedBlocks, err = minesweeper.Visit(x, y)
				assert.NoError(t, err)
				break first
			}
		}
	}

	var visitedBlocks []Block
	for _, row := range game.Blocks {
		for _, block := range row {
			if block.visited {
				visitedBlocks = append(visitedBlocks, block)
			}
		}
	}

	assert.NotEmpty(t, actualVisitedBlocks)

	for _, block1 := range visitedBlocks {
		found := false
		for _, block2 := range actualVisitedBlocks {
			if block1 == block2 {
				found = true
				break
			}
		}
		assert.Truef(t, found, "%v not found in list %v", block1, actualVisitedBlocks)
	}
}

func TestBlockLocationAfterNewGame(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(EASY)
	minesweeper.Play()

	game := minesweeper.(*game)

	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node == BOMB {
				assert.Equal(t, struct{ X, Y int }{X: x, Y: y}, block.Location)
			}
		}
	}

	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node == NUMBER {
				assert.Equal(t, struct{ X, Y int }{X: x, Y: y}, block.Location)
			}
		}
	}

	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node == UNKNOWN {
				assert.Equal(t, struct{ X, Y int }{X: x, Y: y}, block.Location)
			}
		}
	}
}

func TestCheckEventOfGameWhenWinning(t *testing.T) {
	minesweeper, event := NewGame(Grid{SAMPLE_GRID_WIDTH, SAMPLE_GRID_HEIGHT})
	minesweeper.SetDifficulty(EASY)
	minesweeper.Play()

	game := minesweeper.(*game)

	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node != BOMB && !block.visited {
				minesweeper.Visit(x, y)
			}
		}
	}

	go func() {
		time.Sleep(5 * time.Second)
		assert.Fail(t, "Was expecting any event in less than 5 seconds of runtime")
		close(event)
	}()

	if won, ok := <-event; ok {
		assert.Equal(t, WIN, won, "Expecting a winning event")
	} else {
		assert.Fail(t, "Channel event closed. Broken code.")
	}
}

func TestCheckEventOfGameWhenLosing(t *testing.T) {
	minesweeper, event := NewGame(Grid{SAMPLE_GRID_WIDTH, SAMPLE_GRID_HEIGHT})
	minesweeper.SetDifficulty(EASY)
	minesweeper.Play()

	game := minesweeper.(*game)

mainLoop:
	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node == BOMB {
				minesweeper.Visit(x, y)
				break mainLoop
			}
		}
	}

	go func() {
		time.Sleep(5 * time.Second)
		assert.Fail(t, "Was expecting any event in less than 5 seconds of runtime")
		close(event)
	}()

	if won, ok := <-event; ok {
		assert.Equal(t, LOSE, won, "Expecting a losing event")
	} else {
		assert.Fail(t, "Channel event closed. Broken code.")
	}
}

func TestGameEasyDifficultyIsSet(t *testing.T) {
	minesweeper, _ := NewGame(Grid{SAMPLE_GRID_WIDTH, SAMPLE_GRID_HEIGHT})
	minesweeper.SetDifficulty(EASY)
	minesweeper.Play()

	game := minesweeper.(*game)

	assert.Equal(t, EASY, game.Difficulty)
	assert.Equal(t, int(SAMPLE_GRID_WIDTH*SAMPLE_GRID_HEIGHT*EASY_MULTIPLIER), game.totalBombs())
}

func TestGameMediumDifficultyIsSet(t *testing.T) {
	minesweeper, _ := NewGame(Grid{SAMPLE_GRID_WIDTH, SAMPLE_GRID_HEIGHT})
	minesweeper.SetDifficulty(MEDIUM)
	minesweeper.Play()

	game := minesweeper.(*game)

	assert.Equal(t, MEDIUM, game.Difficulty)
	assert.Equal(t, int(SAMPLE_GRID_WIDTH*SAMPLE_GRID_HEIGHT*MEDIUM_MULTIPLIER), game.totalBombs())
}

func TestGameHardDifficultyIsSet(t *testing.T) {
	minesweeper, _ := NewGame(Grid{SAMPLE_GRID_WIDTH, SAMPLE_GRID_HEIGHT})
	minesweeper.SetDifficulty(HARD)
	minesweeper.Play()

	game := minesweeper.(*game)

	assert.Equal(t, HARD, game.Difficulty)
	assert.Equal(t, int(SAMPLE_GRID_WIDTH*SAMPLE_GRID_HEIGHT*HARD_MULTIPLIER), game.totalBombs())
}

func print(game *game) {
	for _, row := range game.Blocks {
		fmt.Println()
		for _, block := range row {
			if block.Node == BOMB {
				fmt.Print("* ")
			} else if block.Node == UNKNOWN {
				fmt.Print("  ")
			} else {
				fmt.Printf("%v ", block.Value)
			}
		}
	}
}
