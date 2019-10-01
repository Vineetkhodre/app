package e2e

import (
	"encoding/json"
	"github.com/deislabs/cnab-go/bundle"
	"gotest.tools/assert"
	"gotest.tools/fs"
	"io/ioutil"
	"path"
	"testing"

	"gotest.tools/icmd"
)

func TestBuild(t *testing.T) {
	cmd, cleanup := dockerCli.createTestCmd()
	defer cleanup()

	testDir := path.Join("testdata", "build")
	dir := fs.NewDir(t, "test-name")
	defer dir.Remove()
	f := dir.Join("bundle.json")
	cmd.Command = dockerCli.Command("app", "build", path.Join(testDir, "single"), "--output", f)
	icmd.RunCmd(cmd).Assert(t, icmd.Success)

	data, err := ioutil.ReadFile(f)
	assert.NilError(t, err)
	var bndl bundle.Bundle
	err = json.Unmarshal(data, &bndl)
	assert.NilError(t, err)

	built := []string { bndl.InvocationImages[0].Digest, bndl.Images["web"].Digest, bndl.Images["worker"].Digest }
	for _, ref := range built {
		cmd.Command = dockerCli.Command("inspect", ref)
		icmd.RunCmd(cmd).Assert(t, icmd.Success)
	}
}
