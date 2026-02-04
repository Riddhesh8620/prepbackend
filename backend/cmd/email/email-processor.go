package email

import (
	"fmt"
	"strings"
)

type TopicItem struct {
	ParentCourse string   // e.g., "Docker Mastery"
	Topics       []string // The specific topics unlocked
}

type CourseItem struct {
	CourseTitle string
}

type ActivationPayload struct {
	Topics       []TopicItem // Supports 1 or many courses
	Courses      []CourseItem
	DashboardURL string
}

func BuildActivationEmail(data ActivationPayload) string {
	var contentRows strings.Builder

	// 1. Handle Full Courses First
	for _, course := range data.Courses {
		contentRows.WriteString(fmt.Sprintf(`
			<tr>
				<td style="padding: 25px 0 10px 0;">
					<span style="background-color: #6366f1; color: #ffffff; font-size: 10px; padding: 3px 8px; border-radius: 6px; font-weight: 900; text-transform: uppercase; letter-spacing: 0.5px;">Full Course</span>
					<h2 style="color: #1e293b; font-size: 19px; font-weight: 800; margin: 10px 0 5px 0;">%s</h2>
					<p style="color: #64748b; font-size: 13px; margin: 0;">You now have lifetime access to all topics in this course.</p>
				</td>
			</tr>`, course.CourseTitle))
	}

	// 2. Handle Individual Topics (Partial Access)
	for _, item := range data.Topics {
		// Course Header for the topics
		contentRows.WriteString(fmt.Sprintf(`
			<tr>
				<td style="padding: 25px 0 10px 0;">
					<span style="background-color: #a855f7; color: #ffffff; font-size: 10px; padding: 3px 8px; border-radius: 6px; font-weight: 900; text-transform: uppercase; letter-spacing: 0.5px;">Partial Access</span>
					<h2 style="color: #1e293b; font-size: 19px; font-weight: 800; margin: 10px 0 5px 0;">%s</h2>
				</td>
			</tr>`, item.ParentCourse))

		// List the specific topics
		for _, topic := range item.Topics {
			contentRows.WriteString(fmt.Sprintf(`
				<tr>
					<td style="padding: 6px 0; color: #475569; font-size: 14px;">
						<table border="0" cellpadding="0" cellspacing="0">
							<tr>
								<td width="22" valign="top" style="color: #6366f1; font-weight: bold;">✓</td>
								<td style="color: #475569;">%s</td>
							</tr>
						</table>
					</td>
				</tr>`, topic))
		}
	}

	// Final Replacement
	body := strings.Replace(Templates.PaymentSuccess, "{{content}}", contentRows.String(), 1)
	body = strings.Replace(body, "{{dashboardUrl}}", data.DashboardURL, 1)

	return body
}
