package common

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vstdy0/go-diploma/pkg"
)

const (
	ctxKeyConfig = pkg.ContextKey("config")
)

// SetConfigToCmdCtx adds Config to cobra.Command context.
func SetConfigToCmdCtx(cmd *cobra.Command, config Config) {
	v := cmd.Context().Value(ctxKeyConfig)
	if v == nil {
		panic(fmt.Errorf("%s context: not set", ctxKeyConfig))
	}

	ctxPtr := v.(*Config)
	*ctxPtr = config
}

// GetConfigFromCmdCtx gets stored Config from cobra.Command context.
func GetConfigFromCmdCtx(cmd *cobra.Command) Config {
	v := cmd.Context().Value(ctxKeyConfig)
	if v == nil {
		panic(fmt.Errorf("%s context: not set", ctxKeyConfig))
	}
	config := v.(*Config)

	return *config
}

// NewBaseCmdCtx creates an empty base context used for storing values for cobra.Command.
func NewBaseCmdCtx() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxKeyConfig, &Config{})

	return ctx
}
