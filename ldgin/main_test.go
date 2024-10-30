/*
 * Copyright (C) distroy
 */

package ldgin

import (
	"log"
	"os"
	"testing"

	"github.com/distroy/ldgo/v3/ldio"
	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	def := log.Default().Writer()
	log.SetOutput(ldio.Discard())
	defer log.SetOutput(def)

	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
