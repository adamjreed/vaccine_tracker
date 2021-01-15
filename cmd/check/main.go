package main

import (
	"covid/internal/app"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Handler() error {
	container, err := app.CreateContainer()
	if err != nil {
		return err
	}

	resp, err := http.Get(container.Config.Check.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if !strings.Contains(string(content), container.Config.Check.Phrase) {
		param, err := container.ParameterStore.GetParam(
			fmt.Sprintf("/%s/vaccine_tracker/triggered", container.Config.Stage),
			false,
		)
		if err != nil {
			return err
		}

		triggered, err := strconv.ParseBool(param)
		if err != nil {
			return err
		}

		if !triggered {
			log.Println("Page change detected. Appointments may be available.")
			log.Println(strings.Replace(string(content), "\n", "", -1))

			for _, to := range container.Config.SMS.To {
				resp, _, _ := container.Twilio.SendSMS(
					container.Config.SMS.From,
					to,
					fmt.Sprintf("COVID vaccine page change detected. Appointments may be available. Go check now: %s. \n\n Reply STOP to stop.", container.Config.Check.URL),
					"",
					"",
				)
				log.Println(fmt.Sprintf("Sent text to %s: %s", to, resp.Sid))
			}

			err = container.ParameterStore.PutParam(
				fmt.Sprintf("/%s/vaccine_tracker/triggered", container.Config.Stage),
				"true",
			)
			if err != nil {
				return err
			}
		} else {
			log.Println("Alarm already triggered.")
		}
	} else {
		err = container.ParameterStore.PutParam(
			fmt.Sprintf("/%s/vaccine_tracker/triggered", container.Config.Stage),
			"false",
		)
		if err != nil {
			return err
		}
		log.Println("No appointments were available.")
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
