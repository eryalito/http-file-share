package services

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

type DialogResponse struct {
	Path  string
	Error error
}

type DialogsService struct{}

func (ds *DialogsService) SelectFile(displayName string) DialogResponse {
	dialog := application.OpenFileDialog()
	dialog.SetTitle("Select a file")
	dialog.SetOptions(&application.OpenFileDialogOptions{
		CanChooseFiles:          true,
		AllowsMultipleSelection: false,
		CanChooseDirectories:    false,
		Directory:               "",
		Filters: []application.FileFilter{
			{
				DisplayName: displayName,
				Pattern:     "*",
			}},
	})

	path, err := dialog.PromptForSingleSelection()
	if err != nil {
		// capture the error if the user cancels the dialog
		if err.Error() == "shellItem is nil" {
			return DialogResponse{}
		}
		return DialogResponse{Error: err}
	}
	return DialogResponse{Path: path}
}
