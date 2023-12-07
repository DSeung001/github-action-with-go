package mail

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
	"io"
	"log"
	"net/http"
	"os"
)

// github watch : https://docs.github.com/en/rest/activity/watching?apiVersion=2022-11-28

var (
	githubAPIURL   = ""
	repoOwner      = ""
	repoName       = ""
	githubToken    = ""
	senderEmail    = ""
	senderEmailPwd = ""
)

func setting() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(".env 파일을 찾을 수 없습니다.")
	}

	githubAPIURL = os.Getenv("GITHUB_API_URL")
	repoOwner = os.Getenv("REPO_OWNER")
	repoName = os.Getenv("REPO_NAME")
	githubToken = os.Getenv("GITHUB_TOKEN")
	senderEmail = os.Getenv("SENDER_EMAIL")
	senderEmailPwd = os.Getenv("SENDER_EMAIL_PWD")
}

func Send() {
	setting()

	// GitHub API에서 repository를 watch한 사용자 목록을 가져오기
	watchers, err := getWatchersList(repoOwner, repoName)
	if err != nil {
		fmt.Println("Error getting watchers:", err)
		return
	}

	// 각 watcher에게 이메일 보내기
	fmt.Println("Sending emails to watchers")
	fmt.Println(watchers)
	for _, watcher := range watchers {
		err := sendEmail(watcher)
		if err != nil {
			fmt.Printf("Error sending email to %s: %s\n", watcher, err)
		} else {
			fmt.Printf("Email sent successfully to %s\n", watcher)
		}
	}
}

func getWatchersList(owner, repo string) ([]string, error) {

	// https://docs.github.com/en/rest/activity/watching?apiVersion=2022-11-28#list-watchers 이걸로
	url := fmt.Sprintf("%s/repos/%s/%s/subscribers", githubAPIURL, owner, repo)
	//fmt.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+githubToken)

	resp, err := http.DefaultClient.Do(req)
	//fmt.Println("resp")
	//fmt.Println(resp)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	//fmt.Println("body")
	//fmt.Println(body)
	//fmt.Println("body string")
	//fmt.Println(string(body))
	if err != nil {
		return nil, err
	}

	// 여기서 에러 발생
	var watchers []map[string]interface{}
	err = json.Unmarshal(body, &watchers)
	if err != nil {
		return nil, err
	}
	//fmt.Println("watchers")
	//fmt.Println(watchers)

	var watcherList []string
	for _, watcher := range watchers {
		//fmt.Println(watcher)
		//fmt.Println(watcher["login"])
		watcherList = append(watcherList, watcher["login"].(string))
	}

	return watcherList, nil
}

// DSeung과 같은 형태로 옴
func sendEmail(to string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)                   // 송신자 이메일 주소
	m.SetHeader("To", to)                              // 수신자 이메일 주소
	m.SetHeader("Subject", "GitHub Repository Update") // 이메일 제목
	m.SetBody("text/plain", "Hello, you are receiving this email because you are watching the GitHub repository.")

	d := gomail.NewDialer("smtp.gmail.com", 587, senderEmail, senderEmailPwd)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
