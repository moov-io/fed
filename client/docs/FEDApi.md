# \FEDApi

All URIs are relative to *http://localhost:8086*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Ping**](FEDApi.md#Ping) | **Get** /ping | Ping the FED service to check if running
[**SearchFEDACH**](FEDApi.md#SearchFEDACH) | **Get** /fed/ach/search | Search FEDACH names and metadata
[**SearchFEDWIRE**](FEDApi.md#SearchFEDWIRE) | **Get** /fed/wire/search | Search FEDWIRE names and metadata



## Ping

> Ping(ctx).Execute()

Ping the FED service to check if running

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.FEDApi.Ping(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `FEDApi.Ping``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiPingRequest struct via the builder pattern


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

> ACHDictionary SearchFEDACH(ctx).XRequestID(xRequestID).XUserID(xUserID).Name(name).RoutingNumber(routingNumber).State(state).City(city).PostalCode(postalCode).Limit(limit).Execute()

Search FEDACH names and metadata

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
    xRequestID := "rs4f9915" // string | Optional Request ID allows application developer to trace requests through the systems logs (optional)
    xUserID := "xUserID_example" // string | Optional User ID used to perform this search (optional)
    name := "Farmers" // string | FEDACH Financial Institution Name (optional)
    routingNumber := "044112187" // string | FEDACH Routing Number for a Financial Institution (optional)
    state := "OH" // string | FEDACH Financial Institution State (optional)
    city := "CALDWELL" // string | FEDACH Financial Institution City (optional)
    postalCode := "43724" // string | FEDACH Financial Institution Postal Code (optional)
    limit := int32(499) // int32 | Maximum results returned by a search (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.FEDApi.SearchFEDACH(context.Background()).XRequestID(xRequestID).XUserID(xUserID).Name(name).RoutingNumber(routingNumber).State(state).City(city).PostalCode(postalCode).Limit(limit).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `FEDApi.SearchFEDACH``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `SearchFEDACH`: ACHDictionary
    fmt.Fprintf(os.Stdout, "Response from `FEDApi.SearchFEDACH`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSearchFEDACHRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xRequestID** | **string** | Optional Request ID allows application developer to trace requests through the systems logs | 
 **xUserID** | **string** | Optional User ID used to perform this search | 
 **name** | **string** | FEDACH Financial Institution Name | 
 **routingNumber** | **string** | FEDACH Routing Number for a Financial Institution | 
 **state** | **string** | FEDACH Financial Institution State | 
 **city** | **string** | FEDACH Financial Institution City | 
 **postalCode** | **string** | FEDACH Financial Institution Postal Code | 
 **limit** | **int32** | Maximum results returned by a search | 

### Return type

[**ACHDictionary**](ACHDictionary.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SearchFEDWIRE

> WIREDictionary SearchFEDWIRE(ctx).XRequestID(xRequestID).XUserID(xUserID).Name(name).RoutingNumber(routingNumber).State(state).City(city).Limit(limit).Execute()

Search FEDWIRE names and metadata

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
    xRequestID := "rs4f9915" // string | Optional Request ID allows application developer to trace requests through the systems logs (optional)
    xUserID := "xUserID_example" // string | Optional User ID used to perform this search (optional)
    name := "MIDWEST" // string | FEDWIRE Financial Institution Name (optional)
    routingNumber := "091905114" // string | FEDWIRE Routing Number for a Financial Institution (optional)
    state := "IA" // string | FEDWIRE Financial Institution State (optional)
    city := "IOWA CITY" // string | FEDWIRE Financial Institution City (optional)
    limit := int32(499) // int32 | Maximum results returned by a search (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.FEDApi.SearchFEDWIRE(context.Background()).XRequestID(xRequestID).XUserID(xUserID).Name(name).RoutingNumber(routingNumber).State(state).City(city).Limit(limit).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `FEDApi.SearchFEDWIRE``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `SearchFEDWIRE`: WIREDictionary
    fmt.Fprintf(os.Stdout, "Response from `FEDApi.SearchFEDWIRE`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSearchFEDWIRERequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xRequestID** | **string** | Optional Request ID allows application developer to trace requests through the systems logs | 
 **xUserID** | **string** | Optional User ID used to perform this search | 
 **name** | **string** | FEDWIRE Financial Institution Name | 
 **routingNumber** | **string** | FEDWIRE Routing Number for a Financial Institution | 
 **state** | **string** | FEDWIRE Financial Institution State | 
 **city** | **string** | FEDWIRE Financial Institution City | 
 **limit** | **int32** | Maximum results returned by a search | 

### Return type

[**WIREDictionary**](WIREDictionary.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

