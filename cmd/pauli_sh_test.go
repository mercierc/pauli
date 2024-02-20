package cmd 

import(
	"testing"
	"os"
	"io/ioutil"
	"fmt"
	"io"
	"strings"
)



func CopyFile(sourceFile string, destinationFile string) {
	input, err := ioutil.ReadFile(sourceFile)
        if err != nil {
                fmt.Println(err)
                return
        }
	
        err = ioutil.WriteFile(destinationFile, input, 0644)
        if err != nil {
                fmt.Println("Error creating", destinationFile)
                fmt.Println(err)
                return
        }
}


func CallPauliDotShFunc(command string) int { 
	// Capture stdout to see if the command in pauli.sh file has been called.
	oldStdout := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w
	
	rootCmd.SetArgs([]string{command})
	rootCmd.Execute()

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = oldStdout

	return strings.Count(string(out), command)
}


func TestMain(m *testing.M) {
	Parse()
	os.Mkdir(".pauli", 0777)
	CopyFile("../data/pauli.sh", ".pauli/pauli.sh")
	CopyFile("../data/config.yaml", ".pauli/config.yaml")

	// Exécuter les tests
	exitCode := m.Run()
	os.RemoveAll(".pauli")
	os.Exit(exitCode)
}


// Ensure the build command in pauli.sh file has been called.
func TestCLIbuild(t *testing.T) {
	command := "build"
	if CallPauliDotShFunc(command)!= 1 {
		t.Fatal(command, "never called in pauli.sh")
	}	
}


// Ensure the run command in pauli.sh file has been called.
func TestCLIrun(t *testing.T) {
	command := "run"
	if CallPauliDotShFunc(command)!= 1 {
		t.Fatal(command, "never called in pauli.sh")
	}	
}


// Ensure the clean command in pauli.sh file has been called.
func TestCLIClean(t *testing.T) {
	command := "clean"
	if CallPauliDotShFunc(command)!= 1 {
		t.Fatal(command, "never called in pauli.sh")
	}	
}


// Ensure the unittests command in pauli.sh file has been called.
func TestCLIunittests(t *testing.T) {
	command := "unittests"
	if CallPauliDotShFunc(command)!= 1 {
		t.Fatal(command, "never called in pauli.sh")
	}	
}


// Ensure the inttests command in pauli.sh file has been called.
func TestCLIinttests(t *testing.T) {
	command := "inttests"
	if CallPauliDotShFunc(command)!= 1 {
		t.Fatal(command, "never called in pauli.sh")
	}	
}


// Ensure the staticanalysis command in pauli.sh file has been called.
func TestCLIstaticanalysis(t *testing.T) {
	command := "staticanalysis"
	if CallPauliDotShFunc(command)!= 1 {
		t.Fatal(command, "never called in pauli.sh")
	}	
}


// Ensure config.yaml and pauli.sh are present.
// Let src/config_test.go TestInitiateProject to verify if the project is
// correctly initialized
func TestCLIinit(t *testing.T) {


	rootCmd.SetArgs([]string{"init"})
	rootCmd.Execute()

	f1, err_sh := os.Open(".pauli/pauli.sh")
	f2, err_yaml := os.Open(".pauli/config.yaml")

	if os.IsNotExist(err_sh) || os.IsNotExist(err_yaml) {
		t.Fatal("\n",
			err_sh,
			"\n",
			err_yaml)
	}

	// Clean
	t.Cleanup(func() {
		f1.Close()
		f2.Close()
	})
}
