package taskrunner

import (
	"net"
	"net/http"
	"net/url"
	"time"
)

func (t *TaskRunner) WaitUntilPortIsReady(port string) {
	t.retryOperation(func() (bool, error) {
		conn, err := net.DialTimeout("tcp", "localhost:"+port, 1*time.Second)
		if err == nil {
			conn.Close()
			return true, nil
		}
		return false, err
	}, "port", "localhost:"+port, 30)
}

func (t *TaskRunner) retryOperation(operation func() (bool, error), description, target string, maxAttempts int) {
	attempt := 0
	for attempt < maxAttempts {
		success, err := operation()
		if success && err == nil {
			Log.Info("%s was requested successfully at %s", description, target)
			return
		} else {
			if attempt%5 == 0 {
				Log.Info("attempt %v/%v: %s is not yet reachable at %s. error: %v. Trying again...", attempt, maxAttempts, description, target, err)
			}
			attempt++
			time.Sleep(1 * time.Second)
		}
	}
	Log.Error("error: %s could not be reached in time at %s. Cleanup and exit...", description, target)
	t.ExitWithError()
}

func (t *TaskRunner) WaitForWebPageToBeReady(targetUrl string) {
	t.retryOperation(func() (bool, error) {
		parsedURL, err := url.Parse(targetUrl)
		if err != nil {
			return false, err
		}

		req, err := http.NewRequest("GET", targetUrl, nil)
		if err != nil {
			return false, err
		}
		req.Header.Set("Origin", parsedURL.Scheme+"://"+parsedURL.Host)

		response, err := http.DefaultClient.Do(req)
		if err == nil && response.StatusCode == 200 {
			return true, nil
		}
		return false, err
	}, "Index page", targetUrl, 60)
}
