package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/samoslab/galtcoin/src/testutil"
	"github.com/samoslab/galtcoin/src/util/file"
)

func Example() {
	// In cmd/cli/cli.go:
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app, err := NewApp(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func TestLoadConfig(t *testing.T) {
	t.Run("set COIN", func(t *testing.T) {
		val := "foocoin"
		os.Setenv("COIN", val)
		defer os.Unsetenv("COIN")

		cfg, err := LoadConfig()
		require.NoError(t, err)
		require.Equal(t, cfg.Coin, val)
	})

	t.Run("set RPC_ADDR", func(t *testing.T) {
		val := "http://111.22.33.44:5555"
		os.Setenv("RPC_ADDR", val)
		defer os.Unsetenv("RPC_ADDR")

		cfg, err := LoadConfig()
		require.NoError(t, err)
		require.Equal(t, cfg.RPCAddress, val)
	})

	t.Run("set RPC_ADDR invalid", func(t *testing.T) {
		val := "111.22.33.44:5555"
		os.Setenv("RPC_ADDR", val)
		defer os.Unsetenv("RPC_ADDR")

		_, err := LoadConfig()
		testutil.RequireError(t, err, "RPC_ADDR must be in scheme://host format")
	})

	t.Run("set WALLET_DIR", func(t *testing.T) {
		val := "/home/foo/bar"
		os.Setenv("WALLET_DIR", val)
		defer os.Unsetenv("WALLET_DIR")

		cfg, err := LoadConfig()
		require.NoError(t, err)
		require.Equal(t, cfg.WalletDir, val)
	})

	t.Run("set WALLET_NAME", func(t *testing.T) {
		val := "bar.wlt"
		os.Setenv("WALLET_NAME", val)
		defer os.Unsetenv("WALLET_NAME")

		cfg, err := LoadConfig()
		require.NoError(t, err)
		require.Equal(t, cfg.WalletName, val)
	})

	t.Run("set WALLET_NAME invalid", func(t *testing.T) {
		val := "badwltext.foo"
		os.Setenv("WALLET_NAME", val)
		defer os.Unsetenv("WALLET_NAME")

		_, err := LoadConfig()
		require.Error(t, err)
		require.Equal(t, ErrWalletName, err)
	})

	t.Run("set DATA_DIR", func(t *testing.T) {
		val := "/home/foo/"
		os.Setenv("DATA_DIR", val)
		defer os.Unsetenv("DATA_DIR")

		cfg, err := LoadConfig()
		require.NoError(t, err)
		require.Equal(t, cfg.DataDir, val)
	})

	t.Run("set WALLET_DIR by DATA_DIR", func(t *testing.T) {
		val := "/home/foo/"
		os.Setenv("DATA_DIR", val)
		defer os.Unsetenv("DATA_DIR")
		valWallet := val + "wallets"
		cfg, err := LoadConfig()
		require.NoError(t, err)
		require.Equal(t, cfg.DataDir, val)
		require.Equal(t, cfg.WalletDir, valWallet)
	})

	t.Run("don't set USE_CSRF", func(t *testing.T) {
		os.Unsetenv("USE_CSRF")

		c, err := LoadConfig()
		require.NoError(t, err)
		require.False(t, c.UseCSRF)
	})

	t.Run("set USE_CSRF false", func(t *testing.T) {
		val := "0"
		os.Setenv("USE_CSRF", val)
		defer os.Unsetenv("USE_CSRF")

		c, err := LoadConfig()
		require.NoError(t, err)
		require.False(t, c.UseCSRF)
	})

	t.Run("set USE_CSRF true", func(t *testing.T) {
		val := "1"
		os.Setenv("USE_CSRF", val)
		defer os.Unsetenv("USE_CSRF")

		c, err := LoadConfig()
		require.NoError(t, err)
		require.True(t, c.UseCSRF)
	})

	t.Run("set USE_CSRF invalid", func(t *testing.T) {
		val := "not_boolean"
		os.Setenv("USE_CSRF", val)
		defer os.Unsetenv("USE_CSRF")

		_, err := LoadConfig()
		testutil.RequireError(t, err, "Invalid USE_CSRF value, must be interpretable as a boolean e.g. 0, 1, true, false")
	})
}

func TestResolveWalletPath(t *testing.T) {
	mustLoadConfig := func() Config {
		cfg, err := LoadConfig()
		require.NoError(t, err)
		return cfg
	}

	defaultCfg := mustLoadConfig()
	expectedFullPath := fmt.Sprintf("%[1]s/.%[2]s/wallets/%[2]s_cli%[3]s", file.UserHome(), defaultCoin, walletExt)
	require.Equal(t, expectedFullPath, defaultCfg.FullWalletPath())

	absPathInput := "./foo/bar.wlt"
	absPathOutput, err := filepath.Abs(absPathInput)
	require.NoError(t, err)

	tests := []struct {
		Name       string
		WalletDir  string
		WalletName string
		WalletArg  string
		Error      error
		Expected   string
	}{
		{
			Name:     "default config, empty cli arg",
			Expected: defaultCfg.WalletDir + "/" + defaultCfg.WalletName,
		},
		{
			Name:      "default config, nonempty cli arg, invalid extension",
			WalletArg: "badwltext.foo",
			Error:     ErrWalletName,
		},
		{
			Name:      "default config, nonempty cli arg, valid",
			WalletArg: "foo.wlt",
			Expected:  defaultCfg.WalletDir + "/foo.wlt",
		},
		{
			Name:       "nondefault config wallet name, empty cli arg",
			WalletName: "bar.wlt",
			Expected:   defaultCfg.WalletDir + "/bar.wlt",
		},
		{
			Name:       "nondefault config wallet name invalid, empty cli arg",
			WalletName: "badwltext.foo",
			Error:      ErrWalletName,
		},
		{
			Name:      "nondefault config wallet path, empty cli arg",
			WalletDir: "/home/foo/somewhere",
			Expected:  "/home/foo/somewhere/" + defaultCfg.WalletName,
		},
		{
			Name:       "nondefault config wallet path and name, empty cli arg",
			WalletDir:  "/home/foo/somewhere",
			WalletName: "bar.wlt",
			Expected:   "/home/foo/somewhere/bar.wlt",
		},
		{
			Name:       "nondefault config wallet path and name, cli arg, valid, name only",
			WalletDir:  "/home/foo/somewhere",
			WalletName: "bar.wlt",
			WalletArg:  "foo.wlt",
			Expected:   "/home/foo/somewhere/foo.wlt",
		},
		{
			Name:       "nondefault config wallet path and name, cli arg, valid, with path",
			WalletDir:  "/home/foo/somewhere",
			WalletName: "bar.wlt",
			WalletArg:  "/usr/lib/foo.wlt",
			Expected:   "/usr/lib/foo.wlt",
		},
		{
			Name:      "wallet arg relative path",
			WalletArg: absPathInput,
			Expected:  absPathOutput,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cfg := mustLoadConfig()
			if tc.WalletDir != "" {
				cfg.WalletDir = tc.WalletDir
			}
			if tc.WalletName != "" {
				cfg.WalletName = tc.WalletName
			}

			w, err := resolveWalletPath(cfg, tc.WalletArg)

			if tc.Error == nil {
				require.NoError(t, err)
				require.Equal(t, tc.Expected, w)
				require.True(t, filepath.IsAbs(w))
			} else {
				require.Error(t, err)
				require.Equal(t, tc.Error, err)
			}
		})
	}
}
