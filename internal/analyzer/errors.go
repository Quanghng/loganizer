package analyzer

import "fmt"

type InaccessibleFileError struct {
	filePath string
	Err 		 error
}

type ParsingError struct {
	filePath string
	Err 		 error
}

func (e *InaccessibleFileError) Error() string {
	return fmt.Sprintf("Le fichier est inaccessible : %s (%v)", e.filePath, e.Err)
}

func (e *ParsingError) Error() string {
	return fmt.Sprintf("Erreur de parsing : %s (%v)", e.filePath, e.Err)
}

func (e *InaccessibleFileError) Unwrap() error {
	return e.Err
}

func (e *ParsingError) Unwrap() error {
	return e.Err
}