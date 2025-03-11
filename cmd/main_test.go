package main

import (
	"fmt"
	"github.com/hokaccha/go-prettyjson"
	"github.com/stretchr/testify/assert"
	"github.com/xLeSHka/mentorLinkSchool/internal/app"
	"github.com/xLeSHka/mentorLinkSchool/internal/pkg/config"
	"go.uber.org/fx"
	"testing"
)

func TestApp(t *testing.T) {
	cfg, err := config.New()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}
	s, _ := prettyjson.Marshal(cfg)
	fmt.Println(string(s))
	err = fx.ValidateApp(fx.Supply(cfg), app.App)
	assert.Nil(t, err)
	err = fx.ValidateApp(fx.Supply(cfg), app.WSApp)
	assert.Nil(t, err)
}
