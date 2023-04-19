# AchParticipant

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**RoutingNumber** | **string** | The institution&#39;s routing number | [optional] 
**OfficeCode** | **string** | Main/Head Office or Branch  * &#x60;O&#x60; - Main * &#x60;B&#x60; - Branch  | [optional] 
**ServicingFRBNumber** | **string** | Servicing Fed&#39;s main office routing number | [optional] 
**RecordTypeCode** | **string** | The code indicating the ABA number to be used to route or send ACH items to the RDFI  * &#x60;0&#x60; - Institution is a Federal Reserve Bank * &#x60;1&#x60; - Send items to customer routing number * &#x60;2&#x60; - Send items to customer using new routing number field  | [optional] 
**Revised** | **string** | Date of last revision  * YYYYMMDD * Blank  | [optional] 
**NewRoutingNumber** | **string** | Financial Institution&#39;s new routing number resulting from a merger or renumber | [optional] 
**CustomerName** | **string** | Financial Institution Name | [optional] 
**AchLocation** | [**AchLocation**](ACHLocation.md) |  | [optional] 
**PhoneNumber** | **string** | The Financial Institution&#39;s phone number | [optional] 
**StatusCode** | **string** | Code is based on the customers receiver code  * &#x60;1&#x60; - Receives Gov/Comm  | [optional] 
**ViewCode** | **string** | Code is current view  * &#x60;1&#x60; - Current view | [optional] 
**CleanName** | **string** | Normalized name of ACH participant | [optional] 
**Logo** | [**Logo**](Logo.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


