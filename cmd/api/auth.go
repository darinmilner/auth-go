package main

import (
	"auth/api/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/pascaldekloe/jwt"
	"github.com/teris-io/shortid"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(pw string) string {
	hashedPW, _ := bcrypt.GenerateFromPassword([]byte(pw), 12)

	fmt.Println(string(hashedPW))
	return string(hashedPW)
}

type ForgotPassword struct {
	Email       string `json:"email"`
	NewPassword string `json:"newPassword"`
	ResetCode   string `json:"resetCode"`
}
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Signup struct {
	Name            string `json:"name" `
	Email           string `json:"email"`
	UserName        string `json:"username"`
	Password        string `json:"password"`
	PassWordConfirm string `json:"passwordConfirm"`
	AccessLevel     string `json:"accessLevel"`
}

//Register function registers Admins
func (app *application) Register(w http.ResponseWriter, r *http.Request) {
	var data Signup

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		app.errorJSON(w, errors.New("error registering"))
		app.logger.Print("decode error", err)
		return
	}

	if data.Name == "" {
		app.errorJSON(w, errors.New("name field is empty"))
		app.logger.Println("name fields are required")
		return
	}

	if data.Email == "" {
		app.errorJSON(w, errors.New("email is required"))
		app.logger.Println("email is required")
		return
	}

	if data.UserName == "" {
		app.errorJSON(w, errors.New("username is required"))
		app.logger.Println("username is required")
		return
	}
	if data.Password == "" || len(data.Password) < 6 {
		app.errorJSON(w, errors.New("password is required and should be at least six characters"))
		app.logger.Println("password is required and should be 6 characters")
		return
	}
	if data.Password != data.PassWordConfirm {
		app.errorJSON(w, errors.New("passwords do not match"))
		app.logger.Println("passwords do not match")
		return
	}

	user := models.User{
		Name:     data.Name,
		Email:    data.Email,
		Password: data.Password,
		UserName: data.UserName,
	}

	i, err := app.db.DB.InsertUser(user)
	user.ID = i
	log.Println(i)
	if err != nil {
		app.logger.Print(err)

		app.errorJSON(w, err)
		return
	}

	userId := int(user.ID)
	userJson := jsonResp{
		OK:      true,
		Message: "New user signed up",
		UserID:  userId,
	}
	err = app.writeJSON(w, http.StatusCreated, userJson, "newUser")
	if err != nil {
		app.logger.Print(err)
		app.errorJSON(w, errors.New("error writing json"))
		return
	}

}

func (app *application) Signin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		app.errorJSON(w, errors.New("Unauthorized"))
		return
	}

	//query DB
	// hashedPw := hashPassword(validUser.Password)
	// log.Print(hashedPw)
	pw := creds.Password
	userName := creds.Username

	id, userPW, err := app.db.DB.Authenticate(userName, pw)

	if err != nil {
		app.errorJSON(w, errors.New("unauthorized, check your login details"), http.StatusForbidden)
		return
	}
	var user models.User
	user.ID = id
	user.Password = userPW
	err = bcrypt.CompareHashAndPassword([]byte(userPW), ([]byte(creds.Password)))
	if err != nil {
		app.errorJSON(w, errors.New("Unauthorized"))
		return
	}

	//create JWT
	var claims jwt.Claims
	claims.Subject = fmt.Sprint(id)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(time.Hour * 24))
	claims.Issuer = "mydomain.com"
	claims.Audiences = []string{"mydomain.com"}

	//create token
	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secret))
	if err != nil {
		app.errorJSON(w, errors.New("error signing in"))
		return
	}

	app.writeJSON(w, http.StatusOK, string(jwtBytes), "token")
}

func (app *application) StatusHandler(w http.ResponseWriter, r *http.Request) {
	currentStatus := AppStatus{
		Status:      "Healthy",
		Environment: app.config.env,
		Version:     version,
	}

	json, err := json.MarshalIndent(currentStatus, "", "\t")
	if err != nil {
		app.logger.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func (app *application) SecuredRoute(w http.ResponseWriter, r *http.Request) {
	// params := httprouter.ParamsFromContext(r.Context())

	// id, err := strconv.Atoi(params.ByName("id"))
	// if err != nil {
	// 	app.logger.Print(errors.New("error getting id from params"))
	// 	app.errorJSON(w, err)
	// 	return
	// }

	// app.logger.Println(id)
	ok := jsonResp{
		OK:      true,
		Message: "You are logged in and can access the secure route",
	}

	err := app.writeJSON(w, http.StatusOK, ok, "products")
	if err != nil {
		app.logger.Print(errors.New("error creating products json"))
		app.errorJSON(w, err)
		return
	}

}

func (app *application) OpenRoute(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Print(errors.New("error getting id from params"))
		app.errorJSON(w, err)
		return
	}

	app.logger.Println(id)
	ok := jsonResp{
		OK:      true,
		Message: "Open route, Anyone can see!!",
		//User's ID can be added here
		UserID: id,
	}

	err = app.writeJSON(w, http.StatusOK, ok, "response")
	if err != nil {
		app.logger.Print(errors.New("error creating response json"))
		app.errorJSON(w, err)
		return
	}

}

func (app *application) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data ForgotPassword
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		app.errorJSON(w, errors.New("error posting forgot password credentials"))
		app.logger.Print("decode error", err)
		return
	}
	log.Println(data.Email)
	var u models.User

	u.Email = data.Email
	user, err := app.db.DB.GetUserByEmail(u.Email)
	if err != nil {
		app.logger.Print(err)
		app.errorJSON(w, errors.New("no user with this email"), http.StatusForbidden)
		return

	}
	resetCode, err := shortid.Generate()
	if err != nil {
		app.logger.Print(err)
		app.errorJSON(w, errors.New("error generating password reset code"))
		return

	}
	app.logger.Println("Reset code: ", resetCode)
	user.Email = data.Email
	user.PasswordResetCode = resetCode
	app.logger.Println("DB Reset code: ", user.PasswordResetCode)

	err = app.db.DB.AddResetPasswordCodeToUser(user)
	if err != nil {
		app.logger.Print(err)
		app.errorJSON(w, errors.New("error adding password reset code to the DB"))
		return

	}

	//prepare to send email

	params := models.ForgotPasswordEmailPayload{
		Source:            "darinmcodingprojects@gmail.com",
		Destination:       user.Email,
		PasswordResetCode: user.PasswordResetCode,
	}

	app.SendEmail(w, params)

}

func (app *application) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var data ForgotPassword
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		app.errorJSON(w, errors.New("error posting forgot password credentials"))
		app.logger.Print("decode error", err)
		return
	}

	resetCode := data.ResetCode
	email := data.Email
	newPassword := data.NewPassword

	password := []byte(newPassword)

	u, err := app.db.DB.GetUserByEmail(data.Email)
	if err != nil {
		app.logger.Print(err)
		app.errorJSON(w, errors.New("no user with this email"), http.StatusForbidden)
		return

	}

	if u.PasswordResetCode != resetCode {
		app.logger.Println("Invalid reset password code ", resetCode)
		app.logger.Println("DB reset code", u.PasswordResetCode)
		app.errorJSON(w, errors.New("invalid reset password code"), http.StatusForbidden)
		return
	}
	// Create a bcrypt hash of the plain-text password.
	hashPassword, err := bcrypt.GenerateFromPassword(password, 12)

	if err != nil {
		app.errorJSON(w, errors.New("error hashing password"))
		app.logger.Print("Hash PW error", err)
		return
	}
	var user models.User
	user.Email = email
	user.Password = string(hashPassword)
	user.PasswordResetCode = resetCode

	err = app.db.DB.UpdateUserPassword(user)
	if err != nil {
		app.errorJSON(w, errors.New("error updating password in DB"))
		app.logger.Print("Save PW in DB error", err)
		return
	}

	app.logger.Println("reset code: ", resetCode, "\npassword: ", newPassword, "\nEmail: ", email)

	respJson := jsonResp{
		OK:      true,
		Message: "Password successfully updated",
		UserID:  u.ID,
	}
	app.writeJSON(w, http.StatusOK, respJson, "pwUpdated")
}
