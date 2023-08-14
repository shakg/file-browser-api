# Go File Browser API

This Go File Browser API is a simple HTTP server that serves two endpoints: `/folder-metadata` and `/file-content`. The API allows you to retrieve metadata for a given folder and download files.

## How to Use

To use this Go File Browser API, you need to provide two command-line arguments when starting the server:

```bash
./filebrowser <folderPath> <PORT>
```

- `<folderPath>`: The absolute or relative path to the folder you want to browse. This folder's metadata will be served through the `/folder-metadata` endpoint.

- `<PORT>`: The port number on which the server will run.

## Downloadables

> You can download binaries from releases page. Or with curl but you need to take a look at the latest release name and replace 'v1.1.0' with release you want to download.

```bash
curl -LJO https://github.com/shakg/file-browser-api/releases/download/v1.1.0/file-browser-api
```

> You can download service file with curl or in Github UI.

```bash
curl -LJO https://raw.githubusercontent.com/shakg/file-browser-api/main/file-browser-api.service
```

## Endpoints

### 1. `/folder-metadata`

**Method:** GET

**Description:** This endpoint retrieves metadata for the specified folder.

**Query Parameters:**
- `path` (string): The path to the folder for which you want to retrieve metadata.

**Response:**
The response will be a JSON array containing metadata for the files and subfolders in the specified folder. Each item in the array will have the following properties:
- `name` (string): The name of the file or folder.
- `size` (int64): The size of the file in bytes. For directories, this will be 0.
- `isDirectory` (bool): A boolean value indicating whether the item is a directory (true) or a file (false).
- `lastModified` (string): The last modification time of the file or folder in ISO 8601 format.

### Example Response
```json
{
  "name": "root",
  "isDirectory": true,
  "children": [
    {
      "name": "folder1",
      "isDirectory": true,
      "children": [
        {
          "name": "file1.txt",
          "size": 1024,
          "isDirectory": false,
          "lastModified": "2023-08-14T12:00:00Z"
        },
        {
          "name": "subfolder1",
          "isDirectory": true,
          "children": [
            {
              "name": "file2.txt",
              "size": 2048,
              "isDirectory": false,
              "lastModified": "2023-08-14T14:30:00Z"
            }
          ]
        }
      ]
    },
    {
      "name": "folder2",
      "isDirectory": true,
      "children": []
    }
  ]
}

```
### 2. `/file-content`

**Method:** GET

**Description:** This endpoint allows you to download the content of a specific file.

**Query Parameters:**
- `fileName` (string): The name of the file you want to download.

**Response:**
The response will be the content of the file with the appropriate `Content-Disposition` header set, allowing the browser to prompt for download with the file's original name.

## File Metadata Struct

The `FileMetadata` struct is used to represent the metadata for a file or folder. It has the following properties:

```go

type FileMetadata struct {
	Name         string        `json:"name"`
	Size         int64         `json:"size"`
	IsDirectory  bool          `json:"isDirectory"`
	LastModified time.Time     `json:"lastModified"`
	Children     []FileMetadata `json:"children,omitempty"` 
}

```

- `Name` (string): The name of the file or folder.
- `Size` (int64): The size of the file in bytes. For directories, this will be 0.
- `IsDirectory` (bool): A boolean value indicating whether the item is a directory (true) or a file (false).
- `LastModified` (time.Time): The last modification time of the file or folder.

## Example Usage

To start the server and browse a folder located at `/path/to/your/folder` on port 8080, run the following command:

```bash
./filebrowser /path/to/your/folder 8080
```

You can then access the endpoints:

- `http://localhost:8080/folder-metadata?path=/path/to/your/folder`
- `http://localhost:8080/file-content?fileName=yourfile.txt`

Please replace `/path/to/your/folder` with the actual path to the folder you want to browse and `yourfile.txt` with the name of the file you want to download.

Remember to test the API with your desired folder and file paths to ensure it meets your requirements.

## Known Issues
> TESTS ARE FAILING, NEED TO REFACTOR TESTS
