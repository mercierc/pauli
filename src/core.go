package src


// Execute bash functions presents in .pauli/pauli.sh
func RunPauliShell(command string, args []string)
// Check the pauli folder is present with its files.
	path, err := exec.LookPath("./.pauli/pauli.sh")
	if err != nil {
		logger.Fatal().Msg("installing pauli.sh is in your future " + path)
	}
	logger.Info().Msgf("Command %s exists", path)

	
	// Construction de la commande Docker
	shell := exec.Command("./.pauli/pauli.sh", "build")

	// Configuration de la sortie standard pour afficher les résultats de la commande
	shell.Stdout = os.Stdout
	shell.Stderr = os.Stderr
	shell.Stdin = os.Stdin

	// Exécution de la commande Docker
	shell.Run()
