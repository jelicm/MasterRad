# MasterRad

## API Endpoints

### Run Application
- **POST** `/runApp`  
  Triggers execution of an application. 

#### Request Body
The request body must be a JSON object matching the following structure:

```go
type AppDTO struct {
    ApplicationId     string `json:"applicationID"`
    ParentNamespaceId string `json:"parentNamespaceId"`
    SizeKB            int    `json:"sizeKB"`
}

``json
{
  "applicationID": "app-123",
  "parentNamespaceId": "namespace-456",
  "sizeKB": 1024
}



### Data Discovery
- **GET** `/dataDiscovery/{nsId}`  
  Runs data discovery and retrieves information (schemas) about all data space items in open state for given namespace ID.

### Add Data Item
- **POST** `/addDataSpaceItem`  
  Adds a new dataspace item to the dataspace.

### Delete Application
- **DELETE** `/deleteApp`  
  Deletes an existing application.

### Create Softlink
- **POST** `/softlink`  
  Creates a soft link between data items.

### Change DataSpaceItem State
- **PUT** `/changeState`  
  Updates the state of a data space item.

### Change Permissions
- **PUT** `/changePermissions`  
  Modifies permissions for a data item or application.

### Put Schema
- **PUT** `/putSchema`  
  Uploads or updates a data schema.

### Delete Application with Merge
- **DELETE** `/deleteAppMerge`  
  Deletes an application and merges related resources.

