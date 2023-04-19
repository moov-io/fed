# \FEDApi

All URIs are relative to *http://localhost:8086*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Ping**](FEDApi.md#Ping) | **Get** /ping | Ping the FED service to check if running
[**SearchFEDACH**](FEDApi.md#SearchFEDACH) | **Get** /fed/ach/search | Search FEDACH names and metadata
[**SearchFEDWIRE**](FEDApi.md#SearchFEDWIRE) | **Get** /fed/wire/search | Search FEDWIRE names and metadata



## Ping

> Ping(ctx, )

Ping the FED service to check if running

### Required Parameters

This endpoint does not need any parameter.

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SearchFEDACH

> AchDictionary SearchFEDACH(ctx, optional)

Search FEDACH names and metadata

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***SearchFEDACHOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a SearchFEDACHOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xRequestID** | **optional.String**| Optional Request ID allows application developer to trace requests through the systems logs | 
 **xUserID** | **optional.String**| Optional User ID used to perform this search | 
 **name** | **optional.String**| FEDACH Financial Institution Name | 
 **routingNumber** | **optional.String**| FEDACH Routing Number for a Financial Institution | 
 **state** | **optional.String**| FEDACH Financial Institution State | 
 **city** | **optional.String**| FEDACH Financial Institution City | 
 **postalCode** | **optional.String**| FEDACH Financial Institution Postal Code | 
 **limit** | **optional.Int32**| Maximum results returned by a search | 

### Return type

[**AchDictionary**](ACHDictionary.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SearchFEDWIRE

> WireDictionary SearchFEDWIRE(ctx, optional)

Search FEDWIRE names and metadata

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***SearchFEDWIREOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a SearchFEDWIREOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xRequestID** | **optional.String**| Optional Request ID allows application developer to trace requests through the systems logs | 
 **xUserID** | **optional.String**| Optional User ID used to perform this search | 
 **name** | **optional.String**| FEDWIRE Financial Institution Name | 
 **routingNumber** | **optional.String**| FEDWIRE Routing Number for a Financial Institution | 
 **state** | **optional.String**| FEDWIRE Financial Institution State | 
 **city** | **optional.String**| FEDWIRE Financial Institution City | 
 **limit** | **optional.Int32**| Maximum results returned by a search | 

### Return type

[**WireDictionary**](WIREDictionary.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

