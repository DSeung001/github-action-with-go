package mail

import (
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
	"log"
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

		email, err := importEmailToName(watcher)
		if err != nil {
			fmt.Println("Error getting watchers:", err)
			return
		}

		fmt.Println("email")
		fmt.Println(email)

		err = sendEmail(watcher)

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
	bodyMap, err := githubGetResBodyMap(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var userNames []string
	for _, watcher := range bodyMap {
		userNames = append(userNames, watcher["login"].(string))
	}

	return userNames, nil
}

func importEmailToName(userName string) ([]string, error) {

	/*
		1. https://api.github.com/users/DSeung001/events/public 에서 payload/commit으로 이메일 찾기
		2. 두번째 방법
	*/

	url := fmt.Sprintf("%s/users/%s/emails", githubAPIURL, userName)
	bodyMap, err := githubGetResBodyMap(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var emails []string
	for _, watcher := range bodyMap {
		fmt.Println(watcher)
		//fmt.Println(watcher["login"])
		emails = append(emails, watcher["email"].(string))
	}
	fmt.Println(emails)

	return emails, nil

}

// DSeung과 같은 형태로 옴
func sendEmail(username string) error {

	// 사용자 명으로 이메일 정보 가져오기 => 이메일 정보를 허용하지 않으면 이메일 보내지 않음
	// https://api.github.com/users/DSeung001/emails

	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)                   // 송신자 이메일 주소
	m.SetHeader("To", username)                        // 수신자 이메일 주소
	m.SetHeader("Subject", "GitHub Repository Update") // 이메일 제목
	m.SetBody("text/plain", "Hello, you are receiving this email because you are watching the GitHub repository.")

	d := gomail.NewDialer("smtp.gmail.com", 587, senderEmail, senderEmailPwd)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
