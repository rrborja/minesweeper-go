/*
 * Minesweeper API
 * Copyright (C) 2017  Ritchie Borja
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License along
 * with this program; if not, write to the Free Software Foundation, Inc.,
 * 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
 */

package minesweeper

import (
	"fmt"
	"testing"

	"github.com/rrborja/minesweeper/rendering"
	"github.com/rrborja/minesweeper/visited"
	"github.com/stretchr/testify/assert"
)

func SampleRenderedGame() Minesweeper {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(Easy)
	minesweeper.Play()
	return minesweeper
}

func TestGameActualBombLocations(t *testing.T) {
	minesweeper := SampleRenderedGame()
	game := minesweeper.(*game)
	properties := minesweeper.(rendering.Tracker)

	bombPlacements := make([]rendering.Position, int(float32(game.Height*game.Width)*game.difficultyMultiplier))

	var counter int
	for _, row := range game.blocks {
		for _, block := range row {
			if block.Node == Bomb {
				bombPlacements[counter] = block
				counter++
			}
		}
	}

	assert.NotEmpty(t, bombPlacements)

	actualBombLocations := properties.BombLocations()
	for i, bomb := range bombPlacements {
		assert.Equal(t, bomb, actualBombLocations[i])
	}
}

func TestGameActualHintLocations(t *testing.T) {
	minesweeper := SampleRenderedGame()
	game := minesweeper.(*game)
	properties := minesweeper.(rendering.Tracker)

	hintPlacements := make([]rendering.Position, 0)

	for _, row := range game.blocks {
		for _, block := range row {
			if block.Node == Number {
				hintPlacements = append(hintPlacements, block)
			}
		}
	}

	assert.NotEmpty(t, hintPlacements)

	actualHintPlacements := properties.HintLocations()
	for i, bomb := range hintPlacements {
		assert.Equal(t, bomb, actualHintPlacements[i])
	}
}

func TestBothBombsAndHintsDoNotShareSameLocations(t *testing.T) {
	minesweeper := SampleRenderedGame()
	properties := minesweeper.(rendering.Tracker)

	hintPlacements := properties.HintLocations()
	bombPlacements := properties.BombLocations()

	assert.NotEmpty(t, hintPlacements)
	assert.NotEmpty(t, bombPlacements)
	for _, hint := range hintPlacements {
		for _, bomb := range bombPlacements {
			if hint.X() == bomb.X() && hint.Y() == bomb.Y() {
				assert.Fail(t, fmt.Sprintf("A hint at %v:%v shares the same location with a bomb at %v:%v", hint.X(), hint.Y(), bomb.X(), bomb.Y()))
			}
		}
	}
}

func TestRecentPlayersMove(t *testing.T) {
	minesweeper, _ := NewGame(Grid{sampleGridWidth, sampleGridHeight})
	minesweeper.SetDifficulty(Medium)
	minesweeper.Play()

	var story visited.StoryTeller = minesweeper.(*game)

	var recentMove visited.Record

	maxMoves := 10
	for i := 0; i < maxMoves; i++ {
		randomX := randomNumber(sampleGridWidth)
		randomY := randomNumber(sampleGridHeight)
		blocks, err := minesweeper.Visit(randomX, randomY)

		if len(blocks) == 0 { // Either already visited block or flagged block
			continue
		}

		var expectedAction visited.Action
		switch blocks[0].Node {
		case Unknown:
			expectedAction = visited.Unknown
		case Number:
			expectedAction = visited.Number
		case Bomb:
			expectedAction = visited.Bomb
		}

		recentMove = visited.Record{Position: blocks[0], Action: expectedAction}

		if err != nil {
			break
		}
	}

	assert.Equal(t, recentMove, story.LastAction())
}
