package main

import (
	"auth/api/models"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type jsonResp struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
	UserID  int    `json:"userId"`
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, wrap string) error {
	wrapper := make(map[string]interface{})

	wrapper[wrap] = data
	json, err := json.Marshal(wrapper)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(json)
	return nil
}

func (app *application) errorJSON(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}
	type jsonError struct {
		Message string `json:"message"`
	}

	theError := jsonError{
		Message: err.Error(),
	}

	app.writeJSON(w, statusCode, theError, "error")
}

func (app *application) SendEmail(w http.ResponseWriter, emailReq models.ForgotPasswordEmailPayload) {
	// var msg EmailForm
	// err := json.NewDecoder(r.Body).Decode(&msg)
	// if err != nil {
	// 	app.errorJSON(w, errors.New("error getting form data"))
	// 	return
	// }

	// log.Println(msg)
	// email := EmailForm{
	// 	FirstName: msg.FirstName,
	// 	LastName:  msg.LastName,
	// 	Email:     msg.Email,
	// 	Message:   msg.Message,
	// }

	postBody, err := json.Marshal(emailReq)
	if err != nil {
		app.logger.Println("Error marshalling json")
		return
	}
	responseBody := bytes.NewBuffer(postBody)
	//Leverage Go's HTTP Post function to make request
	//resp, err := http.Post("http://localhost:7000/gateway/email/reset-password", "application/json", responseBody)

	//A possible test email Test route if using an email service
	resp, err := http.Post("http://localhost:7000/email/test", "application/json", responseBody)

	//Handle Error
	if err != nil {
		app.logger.Println("Reset API has encountered an error or does not exist: ", err)
		app.errorJSON(w, err, http.StatusNotFound)
		return
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		app.logger.Println(err)
	}
	sb := string(body)
	app.logger.Print(sb)

	type jsonResponse struct {
		OK      bool   `json:"ok"`
		Message string `json:"message"`
	}

	emailResp := jsonResponse{
		OK:      true,
		Message: "email sent",
	}

	err = app.writeJSON(w, http.StatusOK, emailResp, "response")
	if err != nil {
		app.errorJSON(w, errors.New("could not send forgot email response json"))
	}

}
