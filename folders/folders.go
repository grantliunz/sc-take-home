package folders

import (
	"github.com/gofrs/uuid"
)

/* GetAllFolders retrieves all folders for a given organization ID

Comments on the original code:

The GetAllFolders function was overly complicated and used unnecessary variables and loops.
Error handling was inconsistent, with some errors being ignored.
The code created multiple unnecessary slices, which could impact performance for large datasets.

Improvements made:

Simplified the GetAllFolders function by directly using the result from FetchAllFoldersByOrgID.
Improved error handling by propagating errors from FetchAllFoldersByOrgID.
Removed unnecessary type conversions and intermediate slices.
*/

func GetAllFolders(req *FetchFolderRequest) (*FetchFolderResponse, error) {
	folders, err := FetchAllFoldersByOrgID(req.OrgID)
	if err != nil {
		return nil, err
	}
	return &FetchFolderResponse{Folders: folders}, nil
}

// This code remains the same as the original code
func FetchAllFoldersByOrgID(orgID uuid.UUID) ([]*Folder, error) {
	folders := GetSampleData()

	resFolder := []*Folder{}
	for _, folder := range folders {
		if folder.OrgId == orgID {
			resFolder = append(resFolder, folder)
		}
	}
	return resFolder, nil
}
