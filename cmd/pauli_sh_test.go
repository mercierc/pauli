package cmd 

import(
	"testing"
	"fmt"
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"

	"github.com/mercierc/pauli/logs"
	"github.com/mercierc/pauli/src"
)


func TestCLIinit(t *testing.T) {

	// Read inputs from temporary file
	file, err := ioutil.TempFile("/tmp", "test")
	defer file.Close()
	if err != nil {
		logs.Logger.Error().Err(err).Msg("error")	
	}

	// Write and replace the cursor at the begining of the file
	_, err = file.WriteString("War\nAnd\nPeace")
	file.Seek(0, 0)

	rootCmd.AddCommand(initCmd)
	rootCmd.SetArgs([]string{"init"})


	_, err = os.Stat(".pauli")
	
	if os.IsNotExist(err) {
	        t.Fatalf(".pauli folder does not exist %v", err)
	}

	// Load template from config.yaml
	content, err := os.ReadFile(".pauli/config.yaml")
	// Parse yaml config  file.
	var confYaml src.Configuration
	err = yaml.Unmarshal(content, &confYaml)

	if (
		confYaml.Name != "War" ||
		confYaml.Builder.Image != "And" ||
		confYaml.Builder.Tag != "Peace") {
		
		t.Fatal("Input for config.yaml are wrong.")
	}
	fmt.Printf("%+v", confYaml)
}

// func TestCLIbuild(t *testing.T) {
// 	rootCmd.SetArgs([]string{"build"})
// 	rootCmd.Execute()
// }

// func TestCLIclean(t *testing.T) {
// 	t.Parallel()
// 	rootCmd.SetArgs([]string{"clean"})
// 	rootCmd.Execute()
// }

// func TestCLIrun(t *testing.T) {
// 	t.Parallel()
// 	rootCmd.SetArgs([]string{"run"})
// 	rootCmd.Execute()
// }

