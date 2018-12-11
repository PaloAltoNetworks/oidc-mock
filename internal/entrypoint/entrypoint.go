package entrypoint

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/common/version"
	flag "github.com/spf13/pflag"
	"go.aporeto.io/oidc-mock/internal/oidcserver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Configuration stuct is used to populate the various fields used by apowine
type Configuration struct {
	ServerIP   string
	ServerPort string

	PrivateKeyPath string
	PublicKeyPath  string

	LogFormat string
	LogLevel  string
}

func usage() {
	flag.PrintDefaults()
	os.Exit(2)
}

func banner(version, revision string) {

	fmt.Printf("\n\x1b[1m\x1b[38;5;6m◼︎ %s\x1b[0m \x1b[38;5;242m%s (%s) %s\n\n\x1b[0m", strings.ToTitle("OIDC-MOCK"), version, revision, "v1.0.0")
}

// StartServer starts the server
func StartServer(cfg *Configuration) {

	banner(version.Version, version.Revision)

	if err := setLogs(cfg.LogFormat, cfg.LogLevel); err != nil {
		log.Fatalf("Error setting up logs: %s", err)
	}

	zap.L().Info("Configuration", zap.Reflect("config", cfg))

	r := mux.NewRouter()
	oidc := oidcserver.NewOIDCServer(cfg.ServerIP, cfg.ServerPort, cfg.PublicKeyPath, cfg.PrivateKeyPath)

	r.HandleFunc("/.well-known/openid-configuration", oidc.ProviderEndpoints).Methods(http.MethodGet)
	r.HandleFunc("/auth", oidc.Authenticate).Methods(http.MethodGet)
	r.HandleFunc("/userInfo", oidc.UserInfo).Methods(http.MethodGet)
	r.HandleFunc("/token", oidc.IssueToken).Methods(http.MethodPost)
	r.HandleFunc("/cert", oidc.IssueCertificate).Methods(http.MethodGet)

	go func() {
		if err := http.ListenAndServeTLS(cfg.ServerPort, ".data/system.crt", ".data/system.key", r); err != nil {
			zap.L().Fatal("Unable to start server",
				zap.String("port", cfg.ServerPort),
				zap.Error(err))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	zap.L().Info("Everything started. Waiting for Stop signal")
	// Waiting for a Sig
	<-c

	zap.L().Info("Server stopped")
}

// setLogs setups Zap to log at the specified log level and format
func setLogs(logFormat, logLevel string) error {
	var zapConfig zap.Config

	switch logFormat {
	case "json":
		zapConfig = zap.NewProductionConfig()
		zapConfig.DisableStacktrace = true
	default:
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.DisableStacktrace = true
		zapConfig.DisableCaller = true
		zapConfig.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {}
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Set the logger
	switch logLevel {
	case "trace":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "debug":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "fatal":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	logger, err := zapConfig.Build()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(logger)
	return nil
}
