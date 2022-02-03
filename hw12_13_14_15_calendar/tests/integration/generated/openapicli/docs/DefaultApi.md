# \DefaultApi

All URIs are relative to *http://127.0.0.1:8080/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CalendarGet**](DefaultApi.md#CalendarGet) | **Get** /calendar | 
[**CreateEvent**](DefaultApi.md#CreateEvent) | **Post** /calendar/event/create | 



## CalendarGet

> string CalendarGet(ctx).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultApi.CalendarGet(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.CalendarGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CalendarGet`: string
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.CalendarGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiCalendarGetRequest struct via the builder pattern


### Return type

**string**

### Authorization

[UserAuth](../README.md#UserAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: plain/text

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateEvent

> string CreateEvent(ctx).EventTemplate(eventTemplate).Execute()



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    eventTemplate := *openapiclient.NewEventTemplate("Id_example", "Title_example", "Datetimestart_example", "Duration_example", "Alarmtime_example") // EventTemplate | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultApi.CreateEvent(context.Background()).EventTemplate(eventTemplate).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.CreateEvent``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateEvent`: string
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.CreateEvent`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateEventRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **eventTemplate** | [**EventTemplate**](EventTemplate.md) |  | 

### Return type

**string**

### Authorization

[UserAuth](../README.md#UserAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: plain/text

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

