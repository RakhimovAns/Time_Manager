package parser

import (
	"errors"
	"github.com/RakhimovAns/Time_Manager/model"
	"strconv"
	"strings"
	"time"
)

func Pars(data []string, botMessage *model.BotMessage) ([]model.Doing, error) {

	var Doings []model.Doing

	for _, doing := range data {

		SplitedData := strings.Fields(doing)
		if len(SplitedData) != 4 {
			botMessage.Text = "invalid type of doing"
			return nil, errors.New("invalid type of doing")
		}

		var Do model.Doing

		Do.Name = SplitedData[0]
		DateTimeStr := SplitedData[1] + " " + SplitedData[2]
		layout := "2.01.2006 15:04"

		dateTime, err := time.Parse(layout, DateTimeStr)
		if err != nil {
			botMessage.Text = "invalid type of doing"
			return nil, errors.New("invalid type of doing")
		}

		Do.Data = dateTime
		Do.Importance, _ = strconv.Atoi(SplitedData[3])
		Doings = append(Doings, Do)
	}
	return Doings, nil
}
