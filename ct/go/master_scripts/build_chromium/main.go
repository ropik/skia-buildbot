// Application that builds chromium with or without patches and uploads the build
// to Google Storage.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/skia-dev/glog"

	"skia.googlesource.com/buildbot.git/ct/go/util"
	"skia.googlesource.com/buildbot.git/go/common"
)

var (
	emails         = flag.String("emails", "", "The comma separated email addresses to notify when the task is picked up and completes.")
	gaeTaskID      = flag.Int("gae_task_id", -1, "The key of the App Engine task. This task will be updated when the task is completed.")
	runID          = flag.String("run_id", "", "The unique run id (typically requester + timestamp).")
	targetPlatform = flag.String("target_platform", util.PLATFORM_ANDROID, "The platform the benchmark will run on (Android / Linux).")
	applyPatches   = flag.Bool("apply_patches", false, "If true looks for Chromium/Blink/Skia patches in temp dir and runs once with the patches and once without.")
	chromiumHash   = flag.String("chromium_hash", "", "The Chromium commit hash the checkout should be synced to. If not specified then Chromium's ToT hash is used.")
	skiaHash       = flag.String("skia_hash", "", "The Skia commit hash the checkout should be synced to. If not specified then Skia's LKGR hash is used (the hash in Chromium's DEPS file).")

	taskCompletedSuccessfully = false
	chromiumBuildTimestamp    = ""
)

func sendEmail(recipients []string) {
	emailSubject := "Chromium build task has completed"
	failureHtml := ""
	if !taskCompletedSuccessfully {
		emailSubject += " with failures"
		failureHtml = util.FailureEmailHtml
	}
	bodyTemplate := `
	The Cluster telemetry queued task to create a new chromium build has completed.<br/>
	%s
	You can schedule more runs <a href="%s">here</a>.<br/><br/>
	Thanks!
	`
	emailBody := fmt.Sprintf(bodyTemplate, failureHtml, util.ChromiumBuildTasksWebapp)
	if err := util.SendEmail(recipients, emailSubject, emailBody); err != nil {
		glog.Errorf("Error while sending email: %s", err)
		return
	}
}

func updateWebappTask() {
	extraData := map[string]string{
		"chromium_rev_date": chromiumBuildTimestamp,
		"build_log_link":    util.MASTER_LOGSERVER_LINK,
	}
	if err := util.UpdateWebappTask(*gaeTaskID, util.UpdateChromiumBuildTasksWebapp, extraData); err != nil {
		glog.Errorf("Error while updating webapp task: %s", err)
		return
	}
}

func main() {
	common.Init()

	// Send start email.
	emailsArr := util.ParseEmails(*emails)
	emailsArr = append(emailsArr, util.CtAdmins...)
	if len(emailsArr) == 0 {
		glog.Error("At least one email address must be specified")
		return
	}
	util.SendTaskStartEmail(emailsArr, "Build chromium")
	// Ensure webapp is updated and completion email is sent even if task fails.
	defer updateWebappTask()
	defer sendEmail(emailsArr)
	// Cleanup tmp files after the run.
	defer util.CleanTmpDir()
	// Finish with glog flush and how long the task took.
	defer util.TimeTrack(time.Now(), "Running build chromium")
	defer glog.Flush()

	if *chromiumHash == "" {
		glog.Error("Must specify --chromium_hash")
		return
	}
	if *skiaHash == "" {
		glog.Error("Must specify --skia_hash")
		return
	}

	if _, _, err := util.CreateChromiumBuild("", *targetPlatform, *chromiumHash, *skiaHash, *applyPatches); err != nil {
		glog.Errorf("Error while creating the Chromium build: %s", err)
		return
	}

	// Find when the requested Chromium revision was submitted.
	stdoutFileName := *runID + ".out"
	stdoutFilePath := filepath.Join(os.TempDir(), stdoutFileName)
	stdoutFile, err := os.Create(stdoutFilePath)
	defer stdoutFile.Close()
	defer os.Remove(stdoutFilePath)
	if err != nil {
		glog.Errorf("Could not create %s: %s", stdoutFilePath, err)
		return
	}
	var chromiumBuildDir string
	if *targetPlatform == "Android" {
		chromiumBuildDir = filepath.Join(util.ChromiumBuildsDir, "android_base")
	} else if *targetPlatform == "Linux" {
		chromiumBuildDir = filepath.Join(util.ChromiumBuildsDir, "linux_base")
	}
	if err := os.Chdir(filepath.Join(chromiumBuildDir, "src")); err != nil {
		glog.Errorf("Could not chdir to %s: %s", chromiumBuildDir, err)
		return
	}
	// Run git log --pretty=format="%at" -1
	if err := util.ExecuteCmd(util.BINARY_GIT, []string{"log", "--pretty=format:%at", "-1"}, []string{}, 5*time.Minute, stdoutFile, nil); err != nil {
		glog.Errorf("Could not run git log cmd: %s", err)
		return
	}
	content, err := ioutil.ReadFile(stdoutFilePath)
	if err != nil {
		glog.Errorf("Could not read %s: %s", stdoutFilePath, err)
		return
	}
	chromiumBuildTimestamp = string(content)

	taskCompletedSuccessfully = true
}
