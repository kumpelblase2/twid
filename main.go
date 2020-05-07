package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var (
	client = &http.Client{}
	port   = "2020"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Not enough arguments specified. Requires at least Client-ID")
		os.Exit(1)
	}

	id := os.Args[1]
	secret := os.Args[2]
	var scopes []string
	if len(os.Args) == 4 {
		scopes = strings.Split(os.Args[3], ",")
	} else {
		scopes = []string{}
	}

	existingToken := GetExistingToken(id)
	if existingToken != "" {
		if IsTokenValid(existingToken) {
			fmt.Println(existingToken)
			os.Exit(0)
		}
	}

	result := make(chan string, 1)
	go func() {
		GetNewToken(id, scopes, result)
	}()

	code := <-result
	token := GetTokenFromCode(id, secret, code)
	SaveToken(id, token)
	fmt.Println(token)
}

func GetExistingToken(id string) string {
	dir := GetCacheDir()
	content, err := ioutil.ReadFile(dir + id)
	if err != nil {
		return ""
	} else {
		return string(content)
	}
}

func GetTokenFromCode(id string, secret string, code string) string {
	resp, _ := http.Post("https://id.twitch.tv/oauth2/token?client_id="+id+"&client_secret="+secret+"&code="+code+"&grant_type=authorization_code&redirect_uri=http://localhost:"+port, "application/json", nil)
	var content map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&content)
	accessToken, _ := content["access_token"]
	return fmt.Sprintf("%v", accessToken)
}

func GetCacheDir() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		os.Stderr.WriteString("Could not get user config directory, using local directory.")
		return "./twid/"
	} else {
		return dir + "/twid/"
	}
}

func SaveToken(id string, token string) {
	dir := GetCacheDir()
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0700)
	}

	ioutil.WriteFile(dir+id, []byte(token), 0600)
}

func GetNewToken(clientId string, scopes []string, result chan string) {
	openBrowser("https://id.twitch.tv/oauth2/authorize?client_id=" + clientId + "&redirect_uri=http://localhost:" + port + "&response_type=code&scope=" + strings.Join(scopes, ","))
	server := &http.Server{
		Addr: "localhost:" + port,
	}
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		token := query.Get("code")
		_, _ = writer.Write([]byte("Authentication success. You can close this tab now."))
		result <- token
		go func() {
			server.Close()
		}()
	})
	_ = server.ListenAndServe()
}

func IsTokenValid(token string) bool {
	req, _ := http.NewRequest("GET", "https://id.twitch.tv/oauth2/validate", nil)
	req.Header.Add("Authorization", "OAuth "+token)
	resp, _ := client.Do(req)
	var content map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&content)
	status, ok := content["status"]
	return !ok || status == 200
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
