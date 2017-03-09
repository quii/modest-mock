package testdata

import (
	"fmt"
	"time"
)

type ClockMock struct {
	Calls struct {
		Now []ClockMock_NowArgs
	}

	Returns struct {
		Now ClockMock_NowReturnsMap
	}
}

func NewClockMock() *ClockMock {
	newMock := new(ClockMock)

	newMock.Returns.Now = make(ClockMock_NowReturnsMap)

	return newMock
}

func (c *ClockMock) Now() (now time.Time) {
	call := ClockMock_NowArgs{}
	c.Calls.Now = append(c.Calls.Now, call)

	if vals, exists := c.Returns.Now[call]; exists {
		return vals.now
	}

	panic(fmt.Sprintf("no return values found for args %+v, ive got %+v", call, c.Returns.Now))
}

type ClockMock_NowReturnsMap map[ClockMock_NowArgs]ClockMock_NowReturns

type ClockMock_NowArgs struct {
}

type ClockMock_NowReturns struct {
	now time.Time
}
