# onedriver-api
A Golang client library for Microsoft Graph REST API, specifically designed for working with OneDrive.

The library provides a simple and easy-to-use API for interacting with OneDrive drives, folders, files, and more.

The library is currently a work in progress and is not yet ready for production use. If you're interested in helping out, please take a look at the [issues](https://github.com/bearcatat/onedriver-api/issues) page and let me know if you have any ideas or suggestions!

## Supported APIs

The library currently supports the following Microsoft Graph REST API endpoints:

### Drives

* [GET /me/drive](https://docs.microsoft.com/en-us/graph/api/drive-get?view=graph-rest-1.0): Get the signed-in user's default drive.

### Drive Items
* [GET /drives/{drive-id}/items/{item-id}](https://docs.microsoft.com/en-us/graph/api/driveitem-get?view=graph-rest-1.0): Retrieve the metadata of a DriveItem by its ID.
* [GET /drives/{drive-id}/root:/{item-path}](https://docs.microsoft.com/en-us/graph/api/driveitem-get?view=graph-rest-1.0): Retrieve the metadata of a DriveItem by its path.
* [POST /drives/{drive-id}/items/{parent-id}/children](https://docs.microsoft.com/en-us/graph/api/driveitem-post-children?view=graph-rest-1.0): Create a new folder under a specified parent DriveItem.
* [POST /drives/{drive-id}/items/{item-id}:/createUploadSession](https://docs.microsoft.com/en-us/graph/api/driveitem-createuploadsession?view=graph-rest-1.0): Create an upload session to upload a large file.
* [PATCH /drives/{drive-id}/items/{item-id}](https://docs.microsoft.com/en-us/graph/api/driveitem-update?view=graph-rest-1.0): Update the properties of a DriveItem.
* [POST /drives/{drive-id}/items/{item-id}/copy](https://docs.microsoft.com/en-us/graph/api/driveitem-copy?view=graph-rest-1.0): Copy a DriveItem to a specified location.
* [DELETE /drives/{drive-id}/items/{item-id}](https://docs.microsoft.com/en-us/graph/api/driveitem-delete?view=graph-rest-1.0): Delete a DriveItem by its ID.
* [PATCH /drives/{drive-id}/items/{item-id}](https://docs.microsoft.com/en-us/graph/api/driveitem-update?view=graph-rest-1.0): Move a DriveItem to a specified location.
* [GET /drives/{drive-id}/items/{item-id}/children](https://docs.microsoft.com/en-us/graph/api/driveitem-list-children?view=graph-rest-1.0): List the children of a DriveItem.
* [GET /drives/{drive-id}/items/{item-id}/content](https://docs.microsoft.com/en-us/graph/api/driveitem-get-content?view=graph-rest-1.0): Download the contents of a DriveItem.

