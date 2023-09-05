package web

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/alemelomeza/improved-octo-memory.git/internal/usecase"
)

type Handlers struct {
	summaryUseCase    *usecase.SummaryUseCase
	evaluationUseCase *usecase.EvaluationUseCase
}

func NewHandlers(summaryUseCase *usecase.SummaryUseCase, evaluationUseCase *usecase.EvaluationUseCase) *Handlers {
	return &Handlers{
		summaryUseCase:    summaryUseCase,
		evaluationUseCase: evaluationUseCase,
	}
}

func (h *Handlers) GetSummaryHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: embed html template
	tmpl := template.Must(template.New("summary").Parse(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Summary conversations</title>
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-9ndCyUaIbzAi2FUVXJi0CjmCapSmO7SnpJef0486qhLnuZ2cdeRhO02iuK6FUUVM" crossorigin="anonymous">
	</head>
	<body>
		<div class="container">
			<main>
				<div class="py-5 text-center">
					<h2>Summary conversations</h2>
					<p class="lead">To get a summary of the conversation enter de <b>clientID</b> and press the botton <b>Send</b>.</span></code> </p>
				</div>
				<form action="/summary" method="post" class="needs-validation" novalidate>
					<div class="row justify-content-center">
						<div class="col-3 align-self-center">
							<input type="number" name="clientId" id="client-id" class="form-control" placeholder="Client ID">
							<hr class="my-4">
							<button type="submit" class="w-100 btn btn-primary btn-lg">Send</button>
							<hr class="my-4">
						</div>
					</div>
				</form>
			</main>
		</div>
	</body>
	</html>
	`))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handlers) PostSummaryHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	input := usecase.SummaryInputDto{
		ClientID: r.Form.Get("clientId"),
	}

	output, err := h.summaryUseCase.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: embed html template
	tmpl := template.Must(template.New("summary").Parse(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Summary conversations</title>
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-9ndCyUaIbzAi2FUVXJi0CjmCapSmO7SnpJef0486qhLnuZ2cdeRhO02iuK6FUUVM" crossorigin="anonymous">
	</head>
	<body>
		<div class="container">
			<main>
				<div class="py-5 text-center">
					<h2>Summary conversations</h2>
					<p class="lead">To get a summary of the conversation enter de <b>clientID</b> and press the botton <b>Send</b>.</span></code> </p>
				</div>
				<form action="/summary" method="post" class="needs-validation" novalidate>
					<div class="row justify-content-center">
						<div class="col-3 align-self-center">
							<input type="number" name="clientId" id="client-id" class="form-control" placeholder="ClientID">
							<hr class="my-4">
							<button type="submit" class="w-100 btn btn-primary btn-lg">Send</button>
							<hr class="my-4">
						</div>
					</div>
				</form>
				{{if ne .Summary ""}}
                <form action="/evaluation" method="post" class="needs-validation" novalidate>
                    <div class="row py-4 justify-content-center">
                        <div class="col-6 align-self-top">
                            <label for="conversation">Conversation</label>
                            <textarea name="conversation" id="conversation" class="form-control h-100 p-5 bg-body-tertiary border rounded-3">
                                {{.Conversation}}
                            </textarea>
                            
                        </div>
                        <div class="col-6 align-self-center">
                            <label for="summary">Summary</label>
                            <textarea name="summary" id="summary" class="form-control h-100 p-5 bg-body-secondary border rounded-3">
                                {{.Summary}}
                            </textarea>
                            <input type="range" class="form-range my-4" min="0" max="4" name="score" id="score">
							<div class="row align-self-center">
								<div class="col text-start">⭐️(1)</div>
								<div class="col text-start">⭐️(2)</div>
								<div class="col text-center">⭐️(3)</div>
								<div class="col text-end">⭐️(4)</div>
								<div class="col text-end">⭐️(5)</div>
							</div>
							<hr class="my-4">
							<button type="submit" class="w-100 btn btn-primary btn-md">Evaluate summary</button>
							<hr class="my-4">
                        </div>
                    </div>
                    <input type="number" name="clientId" id="client-id-evaluation" value="{{.ClientID}}" class="invisible" hidden>
                    <textarea name="prompt" id="prompt" class="invisible" hidden>
                        {{.Prompt}}
                    </textarea>
				</form>
				{{end}}
			</main>
		</div>
	</body>
	</html>
	`))
	err = tmpl.Execute(w, output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handlers) PostEvaluationHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	input := usecase.EvaluationInputDto{
		ClientID:     r.Form.Get("clientId"),
		Conversation: r.Form.Get("conversation"),
		Summary:      r.Form.Get("summary"),
		Prompt:       r.Form.Get("prompt"),
	}
	input.Score, err = strconv.ParseInt(r.Form.Get("score"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.evaluationUseCase.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/summary", http.StatusSeeOther)
}
