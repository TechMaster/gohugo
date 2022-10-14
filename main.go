package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"
)

var success_logger, error_logger *log.Logger

type Hugo_site struct {
	Git_url       string
	Git_folder    string
	Minify_folder string
}

func parse_yaml() {
	config_file, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	data := make(map[string]Hugo_site)
	if err := yaml.Unmarshal(config_file, &data); err != nil {
		log.Fatal(err)
	}
	fmt.Println(data)
}
func handle_githook(w http.ResponseWriter, request *http.Request) {
	var jsonbody map[string]interface{}
	if request.Body == nil {
		http.Error(w, "Please send a request body", 400)
		error_logger.Println("Invalid request body")
		return
	}

	err := json.NewDecoder(request.Body).Decode(&jsonbody)
	if err != nil {
		http.Error(w, err.Error(), 400)
		error_logger.Println(err.Error())
		return
	}
	//Lấy thông tin github URL

	if m1, ok := jsonbody["repository"].(map[string]interface{}); ok {
		var repo_name, repo_url string
		if repo_name, ok = m1["name"].(string); ok {
			if repo_url, ok = m1["url"].(string); ok {
				gitpull_hugo_minify(repo_name, repo_url)
			} else {
				error_logger.Println("cannot find url in payload")
			}
		} else {
			error_logger.Println("cannot find name in payload")
		}
	}
}

/*
Cần log lại lỗi vào một file riêng
*/
func gitpull_hugo_minify(repo_name string, repo_url string) {
	os.Chdir(repo_name)
	currDir, _ := os.Getwd()

	output, err := exec.Command("git", "pull", repo_url).Output()
	if err != nil {
		error_logger.Println(err)
	}
	success_logger.Println("git pull " + repo_url + " to: " + currDir)
	success_logger.Println(string(output))

	output_folder := "/Users/cuong/CODE/gohugo/minify"

	//Render minify to output_folder
	output, err = exec.Command("hugo", "--minify", "-d", output_folder).Output()
	if err != nil {
		fmt.Println(err)
	}
	success_logger.Println("hugo --minify -d " + output_folder)
	success_logger.Println(string(output))

}

func main() {
	success_log_file, err := os.OpenFile("success.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	error_log_file, err := os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer success_log_file.Close()
	defer error_log_file.Close()
	parse_yaml()

	success_logger = log.New(success_log_file, "", log.LstdFlags)
	error_logger = log.New(error_log_file, "", log.LstdFlags)

	http.HandleFunc("/githook", handle_githook)
	http.ListenAndServe(":4567", nil)
}
