# ACHParticipant

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**RoutingNumber** | Pointer to **string** | The institution&#39;s routing number | [optional] 
**OfficeCode** | Pointer to **string** | Main/Head Office or Branch  * &#x60;O&#x60; - Main * &#x60;B&#x60; - Branch  | [optional] 
**ServicingFRBNumber** | Pointer to **string** | Servicing Fed&#39;s main office routing number | [optional] 
**RecordTypeCode** | Pointer to **string** | The code indicating the ABA number to be used to route or send ACH items to the RDFI  * &#x60;0&#x60; - Institution is a Federal Reserve Bank * &#x60;1&#x60; - Send items to customer routing number * &#x60;2&#x60; - Send items to customer using new routing number field  | [optional] 
**Revised** | Pointer to **string** | Date of last revision  * YYYYMMDD * Blank  | [optional] 
**NewRoutingNumber** | Pointer to **string** | Financial Institution&#39;s new routing number resulting from a merger or renumber | [optional] 
**CustomerName** | Pointer to **string** | Financial Institution Name | [optional] 
**AchLocation** | Pointer to [**ACHLocation**](ACHLocation.md) |  | [optional] 
**PhoneNumber** | Pointer to **string** | The Financial Institution&#39;s phone number | [optional] 
**StatusCode** | Pointer to **string** | Code is based on the customers receiver code  * &#x60;1&#x60; - Receives Gov/Comm  | [optional] 
**ViewCode** | Pointer to **string** | Code is current view  * &#x60;1&#x60; - Current view | [optional] 
**CleanName** | Pointer to **string** | Normalized name of ACH participant | [optional] 
**Logo** | Pointer to [**Logo**](Logo.md) |  | [optional] 

## Methods

### NewACHParticipant

`func NewACHParticipant() *ACHParticipant`

NewACHParticipant instantiates a new ACHParticipant object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewACHParticipantWithDefaults

`func NewACHParticipantWithDefaults() *ACHParticipant`

NewACHParticipantWithDefaults instantiates a new ACHParticipant object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRoutingNumber

`func (o *ACHParticipant) GetRoutingNumber() string`

GetRoutingNumber returns the RoutingNumber field if non-nil, zero value otherwise.

### GetRoutingNumberOk

`func (o *ACHParticipant) GetRoutingNumberOk() (*string, bool)`

GetRoutingNumberOk returns a tuple with the RoutingNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRoutingNumber

`func (o *ACHParticipant) SetRoutingNumber(v string)`

SetRoutingNumber sets RoutingNumber field to given value.

### HasRoutingNumber

`func (o *ACHParticipant) HasRoutingNumber() bool`

HasRoutingNumber returns a boolean if a field has been set.

### GetOfficeCode

`func (o *ACHParticipant) GetOfficeCode() string`

GetOfficeCode returns the OfficeCode field if non-nil, zero value otherwise.

### GetOfficeCodeOk

`func (o *ACHParticipant) GetOfficeCodeOk() (*string, bool)`

GetOfficeCodeOk returns a tuple with the OfficeCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOfficeCode

`func (o *ACHParticipant) SetOfficeCode(v string)`

SetOfficeCode sets OfficeCode field to given value.

### HasOfficeCode

`func (o *ACHParticipant) HasOfficeCode() bool`

HasOfficeCode returns a boolean if a field has been set.

### GetServicingFRBNumber

`func (o *ACHParticipant) GetServicingFRBNumber() string`

GetServicingFRBNumber returns the ServicingFRBNumber field if non-nil, zero value otherwise.

### GetServicingFRBNumberOk

`func (o *ACHParticipant) GetServicingFRBNumberOk() (*string, bool)`

GetServicingFRBNumberOk returns a tuple with the ServicingFRBNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServicingFRBNumber

`func (o *ACHParticipant) SetServicingFRBNumber(v string)`

SetServicingFRBNumber sets ServicingFRBNumber field to given value.

### HasServicingFRBNumber

`func (o *ACHParticipant) HasServicingFRBNumber() bool`

HasServicingFRBNumber returns a boolean if a field has been set.

### GetRecordTypeCode

`func (o *ACHParticipant) GetRecordTypeCode() string`

GetRecordTypeCode returns the RecordTypeCode field if non-nil, zero value otherwise.

### GetRecordTypeCodeOk

`func (o *ACHParticipant) GetRecordTypeCodeOk() (*string, bool)`

GetRecordTypeCodeOk returns a tuple with the RecordTypeCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRecordTypeCode

`func (o *ACHParticipant) SetRecordTypeCode(v string)`

SetRecordTypeCode sets RecordTypeCode field to given value.

### HasRecordTypeCode

`func (o *ACHParticipant) HasRecordTypeCode() bool`

HasRecordTypeCode returns a boolean if a field has been set.

### GetRevised

`func (o *ACHParticipant) GetRevised() string`

GetRevised returns the Revised field if non-nil, zero value otherwise.

### GetRevisedOk

`func (o *ACHParticipant) GetRevisedOk() (*string, bool)`

GetRevisedOk returns a tuple with the Revised field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRevised

`func (o *ACHParticipant) SetRevised(v string)`

SetRevised sets Revised field to given value.

### HasRevised

`func (o *ACHParticipant) HasRevised() bool`

HasRevised returns a boolean if a field has been set.

### GetNewRoutingNumber

`func (o *ACHParticipant) GetNewRoutingNumber() string`

GetNewRoutingNumber returns the NewRoutingNumber field if non-nil, zero value otherwise.

### GetNewRoutingNumberOk

`func (o *ACHParticipant) GetNewRoutingNumberOk() (*string, bool)`

GetNewRoutingNumberOk returns a tuple with the NewRoutingNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNewRoutingNumber

`func (o *ACHParticipant) SetNewRoutingNumber(v string)`

SetNewRoutingNumber sets NewRoutingNumber field to given value.

### HasNewRoutingNumber

`func (o *ACHParticipant) HasNewRoutingNumber() bool`

HasNewRoutingNumber returns a boolean if a field has been set.

### GetCustomerName

`func (o *ACHParticipant) GetCustomerName() string`

GetCustomerName returns the CustomerName field if non-nil, zero value otherwise.

### GetCustomerNameOk

`func (o *ACHParticipant) GetCustomerNameOk() (*string, bool)`

GetCustomerNameOk returns a tuple with the CustomerName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCustomerName

`func (o *ACHParticipant) SetCustomerName(v string)`

SetCustomerName sets CustomerName field to given value.

### HasCustomerName

`func (o *ACHParticipant) HasCustomerName() bool`

HasCustomerName returns a boolean if a field has been set.

### GetAchLocation

`func (o *ACHParticipant) GetAchLocation() ACHLocation`

GetAchLocation returns the AchLocation field if non-nil, zero value otherwise.

### GetAchLocationOk

`func (o *ACHParticipant) GetAchLocationOk() (*ACHLocation, bool)`

GetAchLocationOk returns a tuple with the AchLocation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAchLocation

`func (o *ACHParticipant) SetAchLocation(v ACHLocation)`

SetAchLocation sets AchLocation field to given value.

### HasAchLocation

`func (o *ACHParticipant) HasAchLocation() bool`

HasAchLocation returns a boolean if a field has been set.

### GetPhoneNumber

`func (o *ACHParticipant) GetPhoneNumber() string`

GetPhoneNumber returns the PhoneNumber field if non-nil, zero value otherwise.

### GetPhoneNumberOk

`func (o *ACHParticipant) GetPhoneNumberOk() (*string, bool)`

GetPhoneNumberOk returns a tuple with the PhoneNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPhoneNumber

`func (o *ACHParticipant) SetPhoneNumber(v string)`

SetPhoneNumber sets PhoneNumber field to given value.

### HasPhoneNumber

`func (o *ACHParticipant) HasPhoneNumber() bool`

HasPhoneNumber returns a boolean if a field has been set.

### GetStatusCode

`func (o *ACHParticipant) GetStatusCode() string`

GetStatusCode returns the StatusCode field if non-nil, zero value otherwise.

### GetStatusCodeOk

`func (o *ACHParticipant) GetStatusCodeOk() (*string, bool)`

GetStatusCodeOk returns a tuple with the StatusCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusCode

`func (o *ACHParticipant) SetStatusCode(v string)`

SetStatusCode sets StatusCode field to given value.

### HasStatusCode

`func (o *ACHParticipant) HasStatusCode() bool`

HasStatusCode returns a boolean if a field has been set.

### GetViewCode

`func (o *ACHParticipant) GetViewCode() string`

GetViewCode returns the ViewCode field if non-nil, zero value otherwise.

### GetViewCodeOk

`func (o *ACHParticipant) GetViewCodeOk() (*string, bool)`

GetViewCodeOk returns a tuple with the ViewCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetViewCode

`func (o *ACHParticipant) SetViewCode(v string)`

SetViewCode sets ViewCode field to given value.

### HasViewCode

`func (o *ACHParticipant) HasViewCode() bool`

HasViewCode returns a boolean if a field has been set.

### GetCleanName

`func (o *ACHParticipant) GetCleanName() string`

GetCleanName returns the CleanName field if non-nil, zero value otherwise.

### GetCleanNameOk

`func (o *ACHParticipant) GetCleanNameOk() (*string, bool)`

GetCleanNameOk returns a tuple with the CleanName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCleanName

`func (o *ACHParticipant) SetCleanName(v string)`

SetCleanName sets CleanName field to given value.

### HasCleanName

`func (o *ACHParticipant) HasCleanName() bool`

HasCleanName returns a boolean if a field has been set.

### GetLogo

`func (o *ACHParticipant) GetLogo() Logo`

GetLogo returns the Logo field if non-nil, zero value otherwise.

### GetLogoOk

`func (o *ACHParticipant) GetLogoOk() (*Logo, bool)`

GetLogoOk returns a tuple with the Logo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLogo

`func (o *ACHParticipant) SetLogo(v Logo)`

SetLogo sets Logo field to given value.

### HasLogo

`func (o *ACHParticipant) HasLogo() bool`

HasLogo returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


