package acceptance

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf-experimental/disaster-recovery-acceptance-tests/common"
)

type CfAppTestCase struct {
	uniqueTestID string
}

func NewCfAppTestCase(uniqueTestID string) CfAppTestCase {
	return CfAppTestCase{uniqueTestID: uniqueTestID}
}

func (tc CfAppTestCase) PopulateState() {
	By("creating new orgs and spaces")
	urlForDeploymentToBackup, usernameForDeploymentToBackup, passwordForDeploymentToBackup := FindCredentialsFor(DeploymentToBackup())

	RunCommandSuccessfully("cf login --skip-ssl-validation -a", urlForDeploymentToBackup, "-u", usernameForDeploymentToBackup, "-p", passwordForDeploymentToBackup)
	RunCommandSuccessfully("cf create-org acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf create-space acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf target -s acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)
	var testAppFixturePath = "../fixtures/test_app/"
	RunCommandSuccessfully("cf push test_app_" + tc.uniqueTestID + " -p " + testAppFixturePath)
}

func (tc CfAppTestCase) CheckState() {
	By("finding credentials for the deployment to restore")
	urlForDeploymentToRestore, usernameForDeploymentToRestore, passwordForDeploymentToRestore := FindCredentialsFor(DeploymentToRestore())
	RunCommandSuccessfully("cf login --skip-ssl-validation -a", urlForDeploymentToRestore, "-u", usernameForDeploymentToRestore, "-p", passwordForDeploymentToRestore)

	By("verifying apps are back")
	RunCommandSuccessfully("cf target -s acceptance-test-space-" + tc.uniqueTestID + " -o acceptance-test-org-" + tc.uniqueTestID)
	url := GetAppUrl("test_app_" + tc.uniqueTestID)

	Eventually(StatusCode("https://"+url), 5*time.Minute, 5*time.Second).Should(Equal(200))

	By("verify orgs and spaces have been re-created")
	RunCommandSuccessfully("cf org acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf target -o acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf space acceptance-test-space-" + tc.uniqueTestID)
}

func (tc CfAppTestCase) Cleanup() {
	By("cleaning up orgs and spaces")
	urlForDeploymentToBackup, usernameForDeploymentToBackup, passwordForDeploymentToBackup := FindCredentialsFor(DeploymentToBackup())

	RunCommandSuccessfully("cf login --skip-ssl-validation -a", urlForDeploymentToBackup, "-u", usernameForDeploymentToBackup, "-p", passwordForDeploymentToBackup)
	RunCommandSuccessfully("cf target -o acceptance-test-org-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf delete-space -f acceptance-test-space-" + tc.uniqueTestID)
	RunCommandSuccessfully("cf delete-org -f acceptance-test-org-" + tc.uniqueTestID)
}