package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	UploadDir = "/data/uploads"
	OutputDir = "/data/outputs"
	ScriptDir = "/data/scripts"
)

type UploadResponse struct {
	File string `json:"file"`
}

type OutputFile struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type Script struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
}

type RunRequest struct {
	File   string `json:"file"`
	Script string `json:"script"`
}

type RunResponse struct {
	Output string `json:"output"`
}

func deleteOutputHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	name := r.URL.Query().Get("name")

	if name == "" {
		http.Error(
			w,
			"missing name",
			http.StatusBadRequest,
		)
		return
	}
	name = filepath.Base(name)

	err := os.Remove(
		filepath.Join(
			OutputDir,
			name,
		),
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func outputsHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(OutputDir)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var outputs []OutputFile

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()

		if err != nil {
			continue
		}

		outputs = append(outputs, OutputFile{
			Name: file.Name(),
			Size: info.Size(),
		})
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(outputs)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func scriptsHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(ScriptDir)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	var scripts []Script

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		path := filepath.Join(
			ScriptDir,
			file.Name(),
		)

		content, err := os.ReadFile(path)

		if err != nil {
			continue
		}

		lines := strings.Split(
			string(content),
			"\n",
		)

		displayName := file.Name()
		description := "No description"

		for _, line := range lines {
			if strings.HasPrefix(
				line,
				"# NAME:",
			) {
				displayName =
					strings.TrimSpace(
						strings.TrimPrefix(
							line,
							"# NAME:",
						),
					)
			}

			if strings.HasPrefix(
				line,
				"# DESCRIPTION:",
			) {
				description =
					strings.TrimSpace(
						strings.TrimPrefix(
							line,
							"# DESCRIPTION:",
						),
					)
			}
		}

		scripts = append(
			scripts,
			Script{
				Name:        file.Name(),
				DisplayName: displayName,
				Description: description,
			},
		)
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(scripts)
}

func runHandler(w http.ResponseWriter, r *http.Request) {
	var request RunRequest

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := exec.Command(
		ScriptDir+request.Script,
		UploadDir+request.File,
	)

	output, err := cmd.CombinedOutput()

	response := RunResponse{
		Output: string(output),
	}

	if err != nil {
		response.Output += "\nERROR: " + err.Error()
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(100 << 20)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	dst, err := os.Create(UploadDir + header.Filename)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := UploadResponse{
		File: header.Filename,
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	fs := http.FileServer(http.Dir("./web"))

	http.Handle("/", fs)
	http.HandleFunc("/api/upload", uploadHandler)
	http.HandleFunc("/api/scripts", scriptsHandler)
	http.HandleFunc("/api/run", runHandler)
	http.HandleFunc("/api/outputs", outputsHandler)
	http.Handle(
		"/downloads/",
		http.StripPrefix(
			"/downloads/",
			http.FileServer(
				http.Dir(OutputDir),
			),
		),
	)
	http.HandleFunc(
		"/api/output",
		deleteOutputHandler,
	)

	log.Println("Listening on :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
