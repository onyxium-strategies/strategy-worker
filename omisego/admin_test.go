package omisego

import (
	"net/url"
	"os"
	"testing"
)

// Move this to a dotenv file https://github.com/joho/godotenv
var (
	apiKeyId = "api_01cd29gxqbk0a7c859t5v8g4bx"
	apiKey   = "C-b0WYz2L6gvUB-HAwBlcANu0ktoMFTCxJkzKlo3FmU"
	email    = "admin@example.com"
	pwd      = "u22rNF38veC5acIDS1flgA"
	userId   = "usr_01cd29gyb4yrtnf1dmxqm33kbs"

	adminURL = &url.URL{
		Scheme: "http",
		Host:   "localhost:4000",
		Path:   "/admin/api",
	}

	loginBody = LoginParams{
		Email:    email,
		Password: pwd,
	}
)

func TestLogin(t *testing.T) {
	c, _ := NewClient(apiKeyId, apiKey, adminURL)
	adminClient := AdminAPI{
		Client: c,
	}
	body := LoginParams{
		Email:    email,
		Password: pwd,
	}
	_, err := adminClient.Login(body)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLogout(t *testing.T) {
	client, _ := NewClient(apiKeyId, apiKey, adminURL)
	adminClient := AdminAPI{client}
	adminClient.Login(loginBody)
	err := adminClient.Logout()
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccessKeyCreate(t *testing.T) {
	client, _ := NewClient(apiKeyId, apiKey, adminURL)
	adminClient := AdminAPI{client}
	adminClient.Login(loginBody)
	_, err := adminClient.AccessKeyCreate()
	if err != nil {
		t.Fatal(err)
	}
}

func TestAuthTokenSwitchAccount(t *testing.T) {
	client, _ := NewClient(apiKeyId, apiKey, adminURL)
	adminClient := AdminAPI{client}
	adminClient.Login(loginBody)
	body := AuthTokenSwitchAccountParams{
		AccountId: "the_account_id",
	}
	_, err := adminClient.AuthTokenSwitchAccount(body)
	if err.Error() != "{Object:error Code:account:not_found Description:There is no user corresponding to the provided account id Messages:map[]}" {
		t.Fatal(err)
	}
}

func TestPasswordReset(t *testing.T) {
	client, _ := NewClient(apiKeyId, apiKey, adminURL)
	adminClient := AdminAPI{client}
	body := PasswordResetParams{
		Email:       "test@example.com",
		RedirectUrl: "https://example.com/admin/update_password?email={email}&token={token}",
	}
	err := adminClient.PasswordReset(body)
	if err.Error() != "{Object:error Code:user:email_not_found Description:There is no user corresponding to the provided email Messages:map[]}" {
		t.Fatal(err)
	}
}

func TestPasswordUpdate(t *testing.T) {
	client, _ := NewClient(apiKeyId, apiKey, adminURL)
	adminClient := AdminAPI{client}
	body := PasswordUpdateParams{
		Email:                "test@example.com",
		Token:                "26736ca1-43a0-442b-803e-76220cd3cb1d",
		Password:             "nZi9Enc5$l#",
		PasswordConfirmation: "nZi9Enc5$l#",
	}
	err := adminClient.PasswordUpdate(body)
	if err.Error() != "{Object:error Code:user:email_not_found Description:There is no user corresponding to the provided email Messages:map[]}" {
		t.Fatal(err)
	}
}
