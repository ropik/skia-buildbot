package alerting

import (
	"fmt"
	"strings"

	"github.com/skia-dev/glog"
	"skia.googlesource.com/buildbot.git/go/email"
)

const (
	EMAIL_FOOTER = `<br/><br/>
This email was generated by the Skia alert server.
To snooze or dismiss this alert, visit https://skiaalerts.com
`
	EMAIL_SUBJECT_TMPL = "Skia Alert: %s triggered at %s"
)

type Action interface {
	Fire()
	Followup(string)
}

type EmailAction struct {
	rule      *Rule
	to        []string
	subject   string
	emailAuth *email.GMail
}

func (a *EmailAction) Fire() {
	// Cache the email subject so we can send followup emails on the same thread.
	a.subject = fmt.Sprintf(EMAIL_SUBJECT_TMPL, a.rule.Name, a.rule.activeAlert.Triggered().String())
	body := a.rule.Message + EMAIL_FOOTER
	if err := a.emailAuth.Send(a.to, a.subject, body); err != nil {
		glog.Errorf("Failed to send email: %s", err)
	}
}

func (a *EmailAction) Followup(msg string) {
	if err := a.emailAuth.Send(a.to, a.subject, msg); err != nil {
		glog.Errorf("Failed to send email: %s", err)
	}
}

func NewEmailAction(r *Rule, to []string, emailAuth *email.GMail) Action {
	return &EmailAction{
		rule:      r,
		to:        to,
		subject:   "",
		emailAuth: emailAuth,
	}
}

type PrintAction struct {
	rule *Rule
}

func (a *PrintAction) Fire() {
	glog.Infof("ALERT FIRED (%s): %s", a.rule.Name, a.rule.Message)
}

func (a *PrintAction) Followup(msg string) {
	glog.Infof("ALERT FOLLOWUP (%s): %s", a.rule.Name, msg)
}

func NewPrintAction(r *Rule) Action {
	return &PrintAction{r}
}

func parseEmailList(str string) []string {
	split := strings.Split(str, ",")
	emails := []string{}
	for _, email := range split {
		emails = append(emails, strings.Trim(email, " "))
	}
	return emails
}

func (r *Rule) parseActions(actionsInterface interface{}, emailAuth *email.GMail, testing bool) error {
	actionsList := []Action{
		NewPrintAction(r),
	}
	actionStrings := actionsInterface.([]interface{})
	for _, actionString := range actionStrings {
		str := actionString.(string)
		if testing {
			// Do nothing; print is added by default.
		} else if strings.HasPrefix(str, "Email(") && strings.HasSuffix(str, ")") {
			to := parseEmailList(str[6 : len(str)-1])
			actionsList = append(actionsList, NewEmailAction(r, to, emailAuth))
		} else if str == "Print" {
			// Do nothing; print is added by default.
		} else {
			return fmt.Errorf("Unknown action: %q", str)
		}
	}
	r.actions = actionsList
	return nil
}
