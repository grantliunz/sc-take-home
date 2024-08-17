package folders

import (
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/gofrs/uuid"
)

/*
This is an implementation of GetAllFolders with a token-based pagination system. The client receives a NextToken with each response,
which they can use to request the next page of results.

The page token is a base64-encoded string representing the index of the first item in the next page.

All folders for an organization are fetched in memory and then sliced based on the requested page.
This approach works well for moderate datasets but may need to be optimized for very large datasets by implementing pagination at the database level.

The solution supports customizable page sizes, with a default size if none is specified.

Error handling is implemented for invalid page tokens and other potential issues.

This pagination implementation provides a balance between simplicity and functionality, allowing users to efficiently retrieve large sets of folder data in manageable chunks.

*/

const DefaultPageSize = 10

type PaginatedFetchFolderRequest struct {
	OrgID     uuid.UUID 
	PageToken string    
	PageSize  int       
}

type PaginatedFetchFolderResponse struct {
	Folders   []*Folder 
	NextToken string   
}

func GetPaginatedFolders(req *PaginatedFetchFolderRequest) (*PaginatedFetchFolderResponse, error) {
	// Fetch all folders for the organization
	allFolders, err := FetchAllFoldersByOrgID(req.OrgID)
	if err != nil {
		return nil, err
	}

	// Use default page size if not specified or invalid
	if req.PageSize <= 0 {
		req.PageSize = DefaultPageSize
	}

	// Determine the starting index for pagination
	startIndex := 0
	if req.PageToken != "" {
		decodedToken, err := decodePageToken(req.PageToken)
		if err != nil {
			return nil, fmt.Errorf("invalid page token: %v", err)
		}
		startIndex = decodedToken
	}

	// Calculate the end index for the current page
	endIndex := startIndex + req.PageSize
	if endIndex > len(allFolders) {
		endIndex = len(allFolders)
	}

	// Slice the folders for the current page
	paginatedFolders := allFolders[startIndex:endIndex]

	// Generate the next page token if there are more results
	var nextToken string
	if endIndex < len(allFolders) {
		nextToken = encodePageToken(endIndex)
	}

	// Return the paginated response
	return &PaginatedFetchFolderResponse{
		Folders:   paginatedFolders,
		NextToken: nextToken,
	}, nil
}

// Converts an index to a base64-encoded string
func encodePageToken(index int) string {
	return base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(index)))
}

// Converts a base64-encoded string back to an index
func decodePageToken(token string) (int, error) {
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(decoded))
}