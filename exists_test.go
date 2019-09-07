package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

type Sail struct {
	Kind string
}

type Engine struct {
	Power float32
}

type Boat struct {
	Name string
	Sails []Sail
	Engine Engine
}

func TestFirstLevelStructField(t *testing.T) {
	boat := Boat{}
	assert.Equal(t, true, Exists(&boat, []string{"Name"}))
	assert.Equal(t, false, Exists(&boat, []string{"Blarg"}))
	assert.Equal(t, false, Exists(&boat, []string{"0"}))	
}

func TestNestedStructField(t *testing.T) {
	boat := Boat{}
	assert.Equal(t, true, Exists(&boat, []string{"Engine", "Power"}))
	assert.Equal(t, false, Exists(&boat, []string{"Engine", "Sugar"}))
	assert.Equal(t, false, Exists(&boat, []string{"Engine", "0"}))
}

func TestNestedArrayField(t *testing.T) {
	boat := Boat{}
	boat.Sails = []Sail{ Sail{} }
	assert.Equal(t, true, Exists(&boat, []string{"Sails", "0"}))
	assert.Equal(t, false, Exists(&boat, []string{"Sails", "1"}))
	assert.Equal(t, false, Exists(&boat, []string{"Sails", "Jib"}))	
	assert.Equal(t, false, Exists(&boat, []string{"Sails", "0", "Luffing"}))
}
