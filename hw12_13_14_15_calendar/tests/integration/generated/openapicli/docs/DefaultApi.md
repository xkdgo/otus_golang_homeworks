# \DefaultApi

All URIs are relative to *http://127.0.0.1:8080/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CalendarGet**](DefaultApi.md#CalendarGet) | **Get** /calendar | 
[**CreateEvent**](DefaultApi.md#CreateEvent) | **Post** /calendar/event/create | 
[**GetEventsByDay**](DefaultApi.md#GetEventsByDay) | **Get** /calendar/event/day/{date} | 
[**GetEventsByMonth**](DefaultApi.md#GetEventsByMonth) | **Get** /calendar/event/month/{date} | 
[**GetEventsByWeek**](DefaultApi.md#GetEventsByWeek) | **Get** /calendar/event/week/{date} | 



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


## GetEventsByDay

> []EventTemplate GetEventsByDay(ctx, date).Execute()



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
    date := "date_example" // string | date in format \"2006-01-02\"

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultApi.GetEventsByDay(context.Background(), date).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.GetEventsByDay``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetEventsByDay`: []EventTemplate
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.GetEventsByDay`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**date** | **string** | date in format \&quot;2006-01-02\&quot; | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetEventsByDayRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**[]EventTemplate**](EventTemplate.md)

### Authorization

[UserAuth](../README.md#UserAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetEventsByMonth

> []EventTemplate GetEventsByMonth(ctx, date).Execute()



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
    date := "date_example" // string | date in format \"2006-01-02\"

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultApi.GetEventsByMonth(context.Background(), date).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.GetEventsByMonth``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetEventsByMonth`: []EventTemplate
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.GetEventsByMonth`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**date** | **string** | date in format \&quot;2006-01-02\&quot; | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetEventsByMonthRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**[]EventTemplate**](EventTemplate.md)

### Authorization

[UserAuth](../README.md#UserAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetEventsByWeek

> []EventTemplate GetEventsByWeek(ctx, date).Execute()



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
    date := "date_example" // string | date in format \"2006-01-02\"

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultApi.GetEventsByWeek(context.Background(), date).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.GetEventsByWeek``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetEventsByWeek`: []EventTemplate
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.GetEventsByWeek`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**date** | **string** | date in format \&quot;2006-01-02\&quot; | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetEventsByWeekRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**[]EventTemplate**](EventTemplate.md)

### Authorization

[UserAuth](../README.md#UserAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

