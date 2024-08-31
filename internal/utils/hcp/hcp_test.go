package hcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNewHcp(t *testing.T) {
	assert.EqualValues(t, -1.5, GetNewHcp(-1.4, 72, 71))
	assert.EqualValues(t, 4.3, GetNewHcp(4.4, 72, 71))
	assert.EqualValues(t, 4.1, GetNewHcp(4.5, 72, 70))
	assert.EqualValues(t, 11.7, GetNewHcp(12.3, 72, 70))
	assert.EqualValues(t, 18.8, GetNewHcp(20.0, 72, 69))
	assert.EqualValues(t, 27.4, GetNewHcp(28.4, 72, 70))

	assert.EqualValues(t, -1.3, GetNewHcp(-1.4, 72, 74))
	assert.EqualValues(t, 4.5, GetNewHcp(4.4, 72, 74))
	assert.EqualValues(t, 4.6, GetNewHcp(4.5, 72, 75))
	assert.EqualValues(t, 12.4, GetNewHcp(12.3, 72, 76))
	assert.EqualValues(t, 20.1, GetNewHcp(20.0, 72, 77))
	assert.EqualValues(t, 28.5, GetNewHcp(28.4, 72, 78))

	assert.EqualValues(t, 4.4, GetNewHcp(4.4, 72, 72))
	assert.EqualValues(t, 4.5, GetNewHcp(4.5, 72, 73))
	assert.EqualValues(t, 12.3, GetNewHcp(12.3, 72, 74))
	assert.EqualValues(t, 20.0, GetNewHcp(20.0, 72, 72))
	assert.EqualValues(t, 28.4, GetNewHcp(28.4, 72, 76))
}

func TestGetStrokes(t *testing.T) {
	assert.Equal(t, 3, GetStrokes(4.3, 129, 70.5, 72))
	assert.Equal(t, -4, GetStrokes(-1.8, 129, 70.5, 72))
	assert.Equal(t, 10, GetStrokes(9.7, 129, 70.5, 72))
	assert.Equal(t, 18, GetStrokes(19.0, 129, 70.5, 72))

	assert.Equal(t, 9, GetStrokes(5.2, 137, 74.2, 72))
	assert.Equal(t, 15, GetStrokes(10.5, 137, 74.2, 72))
}
