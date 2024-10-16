package scaffold

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func GenerateAPI(args *Args, tmpFolder string, dest string) {
	log.Println("Adding clutch dependencies to go.mod...")
	if err := os.Chdir(filepath.Join(tmpFolder, "backend")); err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command("go", "get", fmt.Sprintf("github.com/%s/clutch/backend@%s", args.Org, args.GoPin))
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Println(string(out))
		log.Fatal("`go get` backend in the destination dir returned the above error")
	}

	cmd = exec.Command("go", "get", fmt.Sprintf("github.com/%s/clutch/tools@%s", args.Org, args.GoPin))
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Println(string(out))
		log.Fatal("`go get` tools in the destination dir returned the above error")
	}

	if err := os.Chdir(tmpFolder); err != nil {
		log.Fatal(err)
	}

	log.Println("Generating API code from protos...")
	log.Println("cd", tmpFolder, "&& make api")
	cmd = exec.Command("make", "api")
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Println(string(out))
		log.Fatal("`make api` in the destination dir returned the above error")
	}
	log.Println("API generation complete")

	fmt.Println("*** All done!")
	fmt.Println("\n*** Try the following command to get started developing the custom gateway:")
	fmt.Printf("cd %s && make\n", dest)
}
