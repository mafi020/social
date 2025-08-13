package templates

import (
	"fmt"

	"github.com/mafi020/social/internal/env"
)

func EmailInvitation(token string) (string, string) {
	// Build invite link
	baseURL := env.GetEnvOrPanic("BASE_URL")
	inviteLink := fmt.Sprintf("%s/api/invitations/accept?token=%s", baseURL, token)

	// Prepare email contents
	plainTextContent := fmt.Sprintf(
		"Hello!\n\nYou have been invited to join OurApp.\nPlease click the link below to accept the invitation:\n%s\n\nThis link will expire in 48 hours.\n\nBest regards,\nThe OurApp Team",
		inviteLink,
	)

	htmlContent := fmt.Sprintf(`
		<html>
			<body style="font-family: Arial, sans-serif; line-height: 1.5;">
				<p>Hello,</p>
				<p>You have been invited to join <strong>Social</strong>.</p>
				<p>
					<a href="%s" style="display: inline-block; padding: 10px 20px; color: white; background-color: #4CAF50; text-decoration: none; border-radius: 5px;">
						Accept Invitation
					</a>
				</p>
				<p>This link will expire in 48 hours.</p>
				<p>Best regards,<br>The Social Team</p>
			</body>
		</html>
	`, inviteLink)

	return plainTextContent, htmlContent
}
