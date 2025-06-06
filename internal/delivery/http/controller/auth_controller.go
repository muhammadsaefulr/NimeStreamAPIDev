package controller

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/muhammadsaefulr/NimeStreamAPI/config"
	auth_request_dto "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/auth/request"
	auth_response_dto "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/auth/response"
	user_dto_request "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/util/response"

	user_model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model/user"
	auth_service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/auth_service"
	system_service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/system_service"
	user_service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/user_service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthController struct {
	AuthService  auth_service.AuthService
	UserService  user_service.UserService
	TokenService system_service.TokenService
	EmailService system_service.EmailService
}

func NewAuthController(
	authService auth_service.AuthService, userService user_service.UserService,
	tokenService system_service.TokenService, emailService system_service.EmailService,
) *AuthController {
	return &AuthController{
		AuthService:  authService,
		UserService:  userService,
		TokenService: tokenService,
		EmailService: emailService,
	}
}

// @Tags         Auth
// @Summary      Register as user
// @Accept       json
// @Produce      json
// @Param        request  body  auth_request_dto.Register  true  "Request body"
// @Router       /auth/register [post]
// @Success      201  {object}  example.RegisterResponse
// @Failure      409  {object}  example.DuplicateEmail  "Email already taken"
func (a *AuthController) Register(c *fiber.Ctx) error {
	req := new(auth_request_dto.Register)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user, err := a.AuthService.Register(c, req)
	if err != nil {
		return err
	}

	tokens, err := a.TokenService.GenerateAuthTokens(c, user)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).
		JSON(response.SuccessWithTokens{
			Code:    fiber.StatusCreated,
			Status:  "success",
			Message: "Register successfully",
			User_id: user.ID.String(),
			Tokens:  *tokens,
		})
}

// @Tags         Auth
// @Summary      Login
// @Accept       json
// @Produce      json
// @Param        request  body  auth_request_dto.Login  true  "Request body"
// @Router       /auth/login [post]
// @Success      200  {object}  example.LoginResponse
// @Failure      401  {object}  example.FailedLogin  "Invalid email or password"
func (a *AuthController) Login(c *fiber.Ctx) error {
	req := new(auth_request_dto.Login)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user, err := a.AuthService.Login(c, req)
	if err != nil {
		return err
	}

	tokens, err := a.TokenService.GenerateAuthTokens(c, user)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.SuccessWithTokens{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Login successfully",
			User_id: user.ID.String(),
			Tokens:  *tokens,
		})
}

// @Tags         Auth
// @Summary      Logout
// @Accept       json
// @Produce      json
// @Param        request  body  example.RefreshToken  true  "Request body"
// @Router       /auth/logout [post]
// @Success      200  {object}  example.LogoutResponse
// @Failure      404  {object}  example.NotFound  "Not found"
func (a *AuthController) Logout(c *fiber.Ctx) error {
	req := new(auth_request_dto.Logout)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := a.AuthService.Logout(c, req); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Common{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Logout successfully",
		})
}

// @Tags         Auth
// @Summary      Refresh auth tokens
// @Accept       json
// @Produce      json
// @Param        request  body  example.RefreshToken  true  "Request body"
// @Router       /auth/refresh-tokens [post]
// @Success      200  {object}  example.RefreshTokenResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
func (a *AuthController) RefreshTokens(c *fiber.Ctx) error {
	req := new(auth_request_dto.RefreshToken)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	tokens, err := a.AuthService.RefreshAuth(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(auth_response_dto.RefreshToken{
			Code:   fiber.StatusOK,
			Status: "success",
			Tokens: *tokens,
		})
}

// @Tags         Auth
// @Summary      Forgot password
// @Description  An email will be sent to reset password.
// @Accept       json
// @Produce      json
// @Param        request  body  auth_request_dto.ForgotPassword  true  "Request body"
// @Router       /auth/forgot-password [post]
// @Success      200  {object}  example.ForgotPasswordResponse
// @Failure      404  {object}  example.NotFound  "Not found"
func (a *AuthController) ForgotPassword(c *fiber.Ctx) error {
	req := new(auth_request_dto.ForgotPassword)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	resetPasswordToken, err := a.TokenService.GenerateResetPasswordToken(c, req)
	if err != nil {
		return err
	}

	if errEmail := a.EmailService.SendResetPasswordEmail(req.Email, resetPasswordToken); errEmail != nil {
		return errEmail
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Common{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "A password reset link has been sent to your email address.",
		})
}

// @Tags         Auth
// @Summary      Reset password
// @Accept       json
// @Produce      json
// @Param        token   query  string  true  "The reset password token"
// @Param        request  body  user_dto_request.UpdatePassOrVerify  true  "Request body"
// @Router       /auth/reset-password [post]
// @Success      200  {object}  example.ResetPasswordResponse
// @Failure      401  {object}  example.FailedResetPassword  "Password reset failed"
func (a *AuthController) ResetPassword(c *fiber.Ctx) error {
	req := new(user_dto_request.UpdatePassOrVerify)
	query := &auth_request_dto.Token{
		Token: c.Query("token"),
	}

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := a.AuthService.ResetPassword(c, query, req); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Common{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Update password successfully",
		})
}

// @Tags         Auth
// @Summary      Send verification email
// @Description  An email will be sent to verify email.
// @Security BearerAuth
// @Produce      json
// @Router       /auth/send-verification-email [post]
// @Success      200  {object}  example.SendVerificationEmailResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
func (a *AuthController) SendVerificationEmail(c *fiber.Ctx) error {
	user, _ := c.Locals("user").(*user_model.User)

	verifyEmailToken, err := a.TokenService.GenerateVerifyEmailToken(c, user)
	if err != nil {
		return err
	}

	if errEmail := a.EmailService.SendVerificationEmail(user.Email, *verifyEmailToken); errEmail != nil {
		return errEmail
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Common{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Please check your email for a link to verify your account",
		})
}

// @Tags         Auth
// @Summary      Verify email
// @Produce      json
// @Param        token   query  string  true  "The verify email token"
// @Router       /auth/verify-email [post]
// @Success      200  {object}  example.VerifyEmailResponse
// @Failure      401  {object}  example.FailedVerifyEmail  "Verify email failed"
func (a *AuthController) VerifyEmail(c *fiber.Ctx) error {
	query := &auth_request_dto.Token{
		Token: c.Query("token"),
	}

	if err := a.AuthService.VerifyEmail(c, query); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Common{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Verify email successfully",
		})
}

// @Tags         Auth
// @Summary      Login with google
// @Description  This route initiates the Google OAuth2 login flow. Please try this in your browser.
// @Router       /auth/google [get]
// @Success      200  {object}  example.GoogleLoginResponse
func (a *AuthController) GoogleLogin(c *fiber.Ctx) error {
	// Generate a random state
	state := uuid.New().String()

	c.Cookie(&fiber.Cookie{
		Name:   "oauth_state",
		Value:  state,
		MaxAge: 30,
	})

	url := config.AppConfig.GoogleLoginConfig.AuthCodeURL(state)

	return c.Status(fiber.StatusSeeOther).Redirect(url)
}

func (a *AuthController) GoogleCallback(c *fiber.Ctx) error {
	state := c.Query("state")
	storedState := c.Cookies("oauth_state")

	if state != storedState {
		return fiber.NewError(fiber.StatusUnauthorized, "States don't Match!")
	}

	code := c.Query("code")
	googlecon := config.GoogleConfig()

	token, err := googlecon.Exchange(context.Background(), code)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		c.Context(), http.MethodGet,
		"https://www.googleapis.com/oauth2/v2/userinfo?access_token="+token.AccessToken,
		nil,
	)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	googleUser := new(user_dto_request.GoogleLogin)
	if errJSON := json.Unmarshal(userData, googleUser); errJSON != nil {
		return errJSON
	}

	user, err := a.UserService.CreateGoogleUser(c, googleUser)
	if err != nil {
		return err
	}

	tokens, err := a.TokenService.GenerateAuthTokens(c, user)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.SuccessWithTokens{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Login successfully",
			User_id: user.ID.String(),
			Tokens:  *tokens,
		})

	// TODO: replace this url with the link to the oauth google success page of your front-end app
	// googleLoginURL := fmt.Sprintf("http://link-to-github.com/muhammadsaefulr/NimeStreamAPI/google/success?access_token=%s&refresh_token=%s",
	// 	tokens.Access.Token, tokens.Refresh.Token)

	// return c.Status(fiber.StatusSeeOther).Redirect(googleLoginURL)
}
