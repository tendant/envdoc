package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/tendant/envdoc"
)

var version = "dev"

func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	rulesPath := flag.String("rules", "", "path to YAML rules file")
	listenAddr := flag.String("listen", "", "HTTP listen address (overrides ENVDOC_LISTEN_ADDR)")
	flag.Parse()

	if *showVersion {
		fmt.Println("envdoc " + version)
		return
	}

	var opts []envdoc.Option

	if *rulesPath != "" {
		rules, err := envdoc.LoadRulesFile(*rulesPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "envdoc: %v\n", err)
			os.Exit(1)
		}
		opts = append(opts, envdoc.WithRules(rules))
	}

	inspector := envdoc.New(opts...)
	_, err := inspector.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "envdoc: %v\n", err)
		os.Exit(1)
	}

	cfg := inspector.Config()
	if cfg.EnableHTTP {
		addr := *listenAddr
		if addr == "" {
			addr = cfg.ListenAddr
		}
		if addr == "" {
			addr = "127.0.0.1:9090"
		}

		mux := http.NewServeMux()
		mux.Handle("/debug/env", inspector.Handler())

		srv := &http.Server{Addr: addr, Handler: mux}

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		go func() {
			fmt.Fprintf(os.Stderr, "envdoc: HTTP server listening on %s\n", addr)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				fmt.Fprintf(os.Stderr, "envdoc: HTTP server error: %v\n", err)
				os.Exit(1)
			}
		}()

		<-ctx.Done()
		fmt.Fprintln(os.Stderr, "envdoc: shutting down HTTP server...")
		srv.Shutdown(context.Background())
	}
}
