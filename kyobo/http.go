package kyobo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetITNewBooks() {
	// => 다음 API로 바꾸자 https://product.kyobobook.co.kr/api/gw/pdt/category/new?page=1&per=20&saleCmdtDvsnCode=KOR&saleCmdtClstCode=33&sort=new
	Config.Domain = "https://product.kyobobook.co.kr"
	Config.Path = "/api/gw/pdt/category/new"
	Config.Page = 1
	Config.Per = 20
	Config.saleCmdtDvsnCode = "KOR"
	Config.saleCmdtClstCode = 33
	Config.sort = "new"

	fmt.Println(Config.GetFullURL())
	resp, err := http.Get(Config.GetFullURL())

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	if err != nil {
		fmt.Println(err)
	} else if resp.StatusCode != 200 {
		fmt.Println("요청 실패 status code:", resp.StatusCode)
	} else {
		fmt.Println("요청 성공")
		fmt.Println(resp)
		fmt.Println(resp.Body)
		body, err := io.ReadAll(resp.Body)
		fmt.Println(body)
		if err != nil {
			fmt.Println("응답 데이터 읽기 오류:", err)
			return
		}

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			fmt.Println("Json에서 구조체로 매핑 중 오류:", err)
			return
		}

		fmt.Println(result)
	}
}
