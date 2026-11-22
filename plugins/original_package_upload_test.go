package plugins

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
	"go.uber.org/zap"
)

func TestSkipVariableReplace(t *testing.T) {
	if !skipVariableReplace("", "k=v") {
		t.Fatal("empty conf file should skip unpack and replace")
	}
	if !skipVariableReplace("app.conf", "") {
		t.Fatal("empty variable list should skip unpack and replace")
	}
	if skipVariableReplace("app.conf", "k=v") {
		t.Fatal("conf file and variable list should run replace")
	}
}

func TestStagePackageForUploadKeepsSourceWhenCopying(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.zip")
	dst := filepath.Join(dir, "dst.zip")
	if err := os.WriteFile(src, []byte("original-bytes"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := stagePackageForUpload(src, dst, false); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(src); err != nil {
		t.Fatal("cache source was removed")
	}
	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "original-bytes" {
		t.Fatalf("staged content = %q", string(got))
	}
}

func TestStagePackageForUploadMovesSource(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "rand_app.zip")
	dst := filepath.Join(dir, "app.zip")
	if err := os.WriteFile(src, []byte("move-me"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := stagePackageForUpload(src, dst, true); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Fatal("downloaded source should be moved")
	}
	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "move-me" {
		t.Fatalf("staged content = %q", string(got))
	}
}

func TestUploadOriginalPackageAlreadyAtS3Target(t *testing.T) {
	log.Logger = zap.NewNop()
	old := DefaultS3TmpAddress
	DefaultS3TmpAddress = ""
	t.Cleanup(func() {
		DefaultS3TmpAddress = old
	})

	endpoint := "http://127.0.0.1:9000/bucket/app.zip"
	got, err := uploadOriginalPackageToS3(endpoint, "", "en")
	if err != nil {
		t.Fatal(err)
	}
	if got != endpoint {
		t.Fatalf("s3PkgPath = %s, want original endpoint when it is already the s3 target", got)
	}
}

func TestNewS3EndpointUsesTmpAddress(t *testing.T) {
	log.Logger = zap.NewNop()
	old := DefaultS3TmpAddress
	DefaultS3TmpAddress = "http://127.0.0.1:9000/salt-tmp"
	t.Cleanup(func() {
		DefaultS3TmpAddress = old
	})

	got := getNewS3EndpointName("http://core.example/artifacts/repository/app.zip", "app.zip")
	want := "http://127.0.0.1:9000/salt-tmp/app.zip"
	if got != want {
		t.Fatalf("s3 endpoint = %s, want %s", got, want)
	}
}
