package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/sign", signform)
	http.HandleFunc("/development", devform)

	http.ListenAndServe(":8090", nil)
}

func signform(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		return
	}

	r.ParseForm()
	imgname := r.FormValue("image")
	key := r.FormValue("key")

	if len(key) == 0 {
		key = "default-signing.key"
	}

	signvars := sign(imgname, key)

	if strings.Contains(signvars, "Error:") {
		fmt.Fprintf(w, signvars)
	} else {
		fmt.Fprintf(w, "Container Image Signed: %s\n", imgname)
		fmt.Fprintf(w, "Signed with: %s\n", key)
	}

}

func devform(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}

	dir, _ := search(r.FormValue("q"))
	for _, r := range dir {
		fmt.Fprintf(w, r+"\n")
	}

}

func search(root string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, nil
}

func sign(image string, key string) string {
	shim := exec.Command("cosign", "sign", "--key", "./safe-signing-keys/"+key, image)
	output, err := shim.CombinedOutput()
	rtr := string(output)

	if err != nil {
		return rtr
	}

	return rtr
}
