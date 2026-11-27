package handlers

import "net/http"

func (h *AppHandler) UploadDocument(w http.ResponseWriter, r *http.Request) {
	const maxUploadSize = 20<<20 //20MB
	const maxMemory = 20<<20 //10MB
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)//governs upload
	if err:= r.ParseMultipartForm(maxUploadSize); err!= nil{//governs loading what protion to memory
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	file, header, err:= r.FormFile("file")
	if err!= nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	
	
}