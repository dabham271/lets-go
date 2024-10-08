package main

import (
  "flag"
  "log/slog"
  "net/http"
  "os"
)

type application struct {
  logger *slog.Logger
}

func main() {

  type config struct {
    addr      string
    staticDir string
  }

  var cfg config

  // Define a new command-line flag with the name 'addr', a default value of ":4000"
  // and some short help text explaining what the flag controls. THe value of the
  // flag will be stored in the addr variable at runtime.
  flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")

  // Importantly, we use the flag.Parse() function to parse the command-line flag.
  // This reads in the command-line flag value and assigns it to the addr
  // variable. You need to call this *before* you use the addr variable
  // otherwise it will always contain the default value of ":4000". If any errors are
  // encountered during parsing, the application will be terminated.
  flag.Parse()

  // Use the slog.New() function to initialize a new structured logger, which
  // writes to the standard out stream and uses the default settings.
  logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

  // Initialize a new instance of our application struct, containing the
  // dependencies (for now, just the structured logger).
  app := &application{
    logger: logger,
  }

  mux := http.NewServeMux()

  // Create a file server which serves files out of the "./ui/static" directory.
  // Note that the path given to the http.Dir function is relative to the project
  // directory root.
  fileServer := http.FileServer(http.Dir("./ui/static/"))

  // Use the mux.Handle() function to register the file server as the handler for
  // all URL paths that start with "/static/". For matching paths, we strip the
  // "/static" prefix before the request reaches the file server.
  mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

  // Swap the route declarations to use the application struct's methods as the
  // handler functions.
  mux.HandleFunc("GET /{$}", app.home)
  mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
  mux.HandleFunc("GET /snippet/create", app.snippetCreate)
  mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

  // Use the Info() method to log the starting server message at Info severity
  // (along with the listen address as an attribute).
  logger.Info("starting server", slog.String("addr", cfg.addr))

  // And we pass the dereferenced addr pointer to http.ListenAndServe() too.
  err := http.ListenAndServe(cfg.addr, mux)

  // And we also use the Error() method to log any error message returned by
  // http.ListenAndServe() at Error severity (with no additional attributes),
  // and then call os.Exit(1) to terminate the application with exit code 1.
  logger.Error(err.Error())
  os.Exit(1)
}
