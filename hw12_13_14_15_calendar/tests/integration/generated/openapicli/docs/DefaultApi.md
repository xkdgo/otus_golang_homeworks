# \DefaultApi

All URIs are relative to *http://127.0.0.1:8080/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CalendarGet**](DefaultApi.md#CalendarGet) | **Get** /calendar | 
[**CreateEvent**](DefaultApi.md#CreateEvent) | **Post** /calendar/event/create | 



## CalendarGet

> string CalendarGet(ctx, )



Hello World

### Required Parameters

This endpoint does not need any parameter.

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

> string CreateEvent(ctx, eventTemplate)



### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**eventTemplate** | [**EventTemplate**](EventTemplate.md)|  | 

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

