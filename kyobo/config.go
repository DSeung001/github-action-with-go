package kyobo

import "fmt"

type KyoboConfig struct {
	Domain           string `json:"domain"`
	Path             string `json:"path"`
	Page             int    `json:"page"`
	Per              int    `json:"per"`
	saleCmdtDvsnCode string `json:"saleCmdtDvsnCode"`
	saleCmdtClstCode int    `json:"saleCmdtClstCode"`
	sort             string `json:"sort"`
}

var Config = KyoboConfig{}

func (KyoboConfig) GetFullURL() string {
	return fmt.Sprintf("%s%s?page=%d&per=%d&saleCmdtDvsnCode=%s&saleCmdtClstCode=%d&sort=%s", Config.Domain, Config.Path, Config.Page, Config.Per, Config.saleCmdtDvsnCode, Config.saleCmdtClstCode, Config.sort)
}
