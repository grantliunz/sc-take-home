package folders_test

import (
	"testing"

	"github.com/georgechieng-sc/interns-2022/folders"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAllFolders(t *testing.T) {
	t.Run("ValidOrgID_ReturnsMatchingFolders", func(t *testing.T) {
		req := &folders.FetchFolderRequest{
			OrgID: uuid.FromStringOrNil(folders.DefaultOrgID),
		}
		
		response, err := folders.GetAllFolders(req)
		
		assert.Nil(t, err)
		assert.NotNil(t, response)
		
		if len(response.Folders) > 0 {
			for _, folder := range response.Folders {
				assert.Equal(t, req.OrgID, folder.OrgId)
			}
		}
	})

	t.Run("NilOrgID_ReturnsEmptyList", func(t *testing.T) {
		req := &folders.FetchFolderRequest{
			OrgID: uuid.Nil,
		}
		
		response, err := folders.GetAllFolders(req)
		
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Empty(t, response.Folders)
	})
}

func TestGetPaginatedFolders(t *testing.T) {
	// Create a test org ID
	orgID := uuid.FromStringOrNil(folders.DefaultOrgID)
	

	t.Run("FirstPage_DefaultPageSize", func(t *testing.T) {
		req := &folders.PaginatedFetchFolderRequest{
			OrgID: orgID,
		}

		response, err := folders.GetPaginatedFolders(req)

		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Folders, folders.DefaultPageSize)
		assert.NotEmpty(t, response.NextToken)
	})

	t.Run("CustomPageSize", func(t *testing.T) {
		customPageSize := 5
		req := &folders.PaginatedFetchFolderRequest{
			OrgID:    orgID,
			PageSize: customPageSize,
		}

		response, err := folders.GetPaginatedFolders(req)

		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Folders, customPageSize)
		assert.NotEmpty(t, response.NextToken)
	})

	t.Run("FetchAllPages", func(t *testing.T) {
		req := &folders.PaginatedFetchFolderRequest{
			OrgID:    orgID,
			PageSize: 2, 
		}

		var allFolders []*folders.Folder
		var nextToken string

		for {
			req.PageToken = nextToken
			response, err := folders.GetPaginatedFolders(req)

			assert.Nil(t, err)
			assert.NotNil(t, response)

			allFolders = append(allFolders, response.Folders...)
			nextToken = response.NextToken

			if nextToken == "" {
				break
			}
		}

		// Fetch all folders without pagination to compare
		allFoldersResponse, err := folders.GetAllFolders(&folders.FetchFolderRequest{OrgID: orgID})
		assert.Nil(t, err)

		assert.Equal(t, len(allFoldersResponse.Folders), len(allFolders))
	})

	t.Run("InvalidPageToken", func(t *testing.T) {
		req := &folders.PaginatedFetchFolderRequest{
			OrgID:     orgID,
			PageToken: "invalid_token",
		}

		response, err := folders.GetPaginatedFolders(req)

		assert.NotNil(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "invalid page token")
	})

	t.Run("EmptyResult", func(t *testing.T) {
		emptyOrgID := uuid.Must(uuid.NewV4()) // Assuming this org has no folders
		req := &folders.PaginatedFetchFolderRequest{
			OrgID: emptyOrgID,
		}

		response, err := folders.GetPaginatedFolders(req)

		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Empty(t, response.Folders)
		assert.Empty(t, response.NextToken)
	})
}
