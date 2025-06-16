package main

import (
	"context"
	"log"
	"net/http"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	"html/template"
)

type envConfig struct {
	Port  int    `envconfig:"PORT" default:"8080" required:"true"`
	Sink  string `envconfig:"K_SINK"`
	Red   int    `envconfig:"COLOR_RED" default:"255"`
	Green int    `envconfig:"COLOR_GREEN" default:"0"`
	Blue  int    `envconfig:"COLOR_BLUE" default:"0"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("failed to process env var: %s", err)
	}

	handler := new(Handler)

	if env.Sink != "" {
		if client, err := cloudevents.NewClient(cloudevents.WithTarget(env.Sink)); err != nil {
			log.Printf("failed to make cloudevents client: %v\n", err)
		} else {
			handler.client = client
		}
	}

	log.Printf("Server starting on port :%d\n", env.Port)
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("failed to start server, %s", err.Error())
	}
}

/*

http://colorgrid.default.d2k.n3wscott.com/?url=http://colors.default.d2k.n3wscott.com&count=25

*/

var index = `<!DOCTYPE html>
<html>
  <head>
    <style>
      body {
        background-color: rgb({{ .red }}, {{ .green }}, {{ .blue }});
      }
    </style>
  </head>
  <body></body>
</html>`

var indexTemplate *template.Template

func init() {
	indexTemplate = template.Must(template.New("index").Parse(index))
}

type Handler struct {
	red    int
	green  int
	blue   int
	client cloudevents.Client
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet || r.URL.Path != "/" {
		return
	}

	color := map[string]int{
		"red":   h.red,
		"green": h.green,
		"blue":  h.blue,
	}
	_ = indexTemplate.Execute(w, color)

	if h.client != nil {
		if result := h.client.Send(context.Background(), newEvent(color)); cloudevents.IsUndelivered(result) {
			log.Printf("failed to send cloudevent: %v\n", result.Error())
		}
	}
}

func newEvent(data interface{}) cloudevents.Event {
	event := cloudevents.NewEvent() // Sets version
	event.SetType("com.n3wscott.atlanta.colors")
	event.SetSource("github.com/n3wscott/k8s-meetup-atlanta/cmd/colors")
	if err := event.SetData(cloudevents.ApplicationJSON, data); err != nil {
		log.Printf("failed to cloudevents event: %v\n", err)
	}
	return event
}
