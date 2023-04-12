# WIREParticipant

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**RoutingNumber** | Pointer to **string** | The institution&#39;s routing number | [optional] 
**TelegraphicName** | Pointer to **string** | Short name of financial institution | [optional] 
**CustomerName** | Pointer to **string** | Financial Institution Name | [optional] 
**WireLocation** | Pointer to [**WIRELocation**](WIRELocation.md) |  | [optional] 
**FundsTransferStatus** | Pointer to **string** | Designates funds transfer status  * &#x60;Y&#x60; - Eligible * &#x60;N&#x60; - Ineligible  | [optional] 
**FundsSettlementOnlyStatus** | Pointer to **string** | Designates funds settlement only status   * &#x60;S&#x60; - Settlement-Only  | [optional] 
**BookEntrySecuritiesTransferStatus** | Pointer to **string** | Designates book entry securities transfer status  * &#x60;Y&#x60; - Eligible * &#x60;N&#x60; - Ineligible  | [optional] 
**Date** | Pointer to **string** | Date of last revision  * YYYYMMDD * Blank  | [optional] 
**CleanName** | Pointer to **string** | Normalized name of Wire participant | [optional] 
**Logo** | Pointer to [**Logo**](Logo.md) |  | [optional] 

## Methods

### NewWIREParticipant

`func NewWIREParticipant() *WIREParticipant`

NewWIREParticipant instantiates a new WIREParticipant object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWIREParticipantWithDefaults

`func NewWIREParticipantWithDefaults() *WIREParticipant`

NewWIREParticipantWithDefaults instantiates a new WIREParticipant object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRoutingNumber

`func (o *WIREParticipant) GetRoutingNumber() string`

GetRoutingNumber returns the RoutingNumber field if non-nil, zero value otherwise.

### GetRoutingNumberOk

`func (o *WIREParticipant) GetRoutingNumberOk() (*string, bool)`

GetRoutingNumberOk returns a tuple with the RoutingNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRoutingNumber

`func (o *WIREParticipant) SetRoutingNumber(v string)`

SetRoutingNumber sets RoutingNumber field to given value.

### HasRoutingNumber

`func (o *WIREParticipant) HasRoutingNumber() bool`

HasRoutingNumber returns a boolean if a field has been set.

### GetTelegraphicName

`func (o *WIREParticipant) GetTelegraphicName() string`

GetTelegraphicName returns the TelegraphicName field if non-nil, zero value otherwise.

### GetTelegraphicNameOk

`func (o *WIREParticipant) GetTelegraphicNameOk() (*string, bool)`

GetTelegraphicNameOk returns a tuple with the TelegraphicName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTelegraphicName

`func (o *WIREParticipant) SetTelegraphicName(v string)`

SetTelegraphicName sets TelegraphicName field to given value.

### HasTelegraphicName

`func (o *WIREParticipant) HasTelegraphicName() bool`

HasTelegraphicName returns a boolean if a field has been set.

### GetCustomerName

`func (o *WIREParticipant) GetCustomerName() string`

GetCustomerName returns the CustomerName field if non-nil, zero value otherwise.

### GetCustomerNameOk

`func (o *WIREParticipant) GetCustomerNameOk() (*string, bool)`

GetCustomerNameOk returns a tuple with the CustomerName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCustomerName

`func (o *WIREParticipant) SetCustomerName(v string)`

SetCustomerName sets CustomerName field to given value.

### HasCustomerName

`func (o *WIREParticipant) HasCustomerName() bool`

HasCustomerName returns a boolean if a field has been set.

### GetWireLocation

`func (o *WIREParticipant) GetWireLocation() WIRELocation`

GetWireLocation returns the WireLocation field if non-nil, zero value otherwise.

### GetWireLocationOk

`func (o *WIREParticipant) GetWireLocationOk() (*WIRELocation, bool)`

GetWireLocationOk returns a tuple with the WireLocation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWireLocation

`func (o *WIREParticipant) SetWireLocation(v WIRELocation)`

SetWireLocation sets WireLocation field to given value.

### HasWireLocation

`func (o *WIREParticipant) HasWireLocation() bool`

HasWireLocation returns a boolean if a field has been set.

### GetFundsTransferStatus

`func (o *WIREParticipant) GetFundsTransferStatus() string`

GetFundsTransferStatus returns the FundsTransferStatus field if non-nil, zero value otherwise.

### GetFundsTransferStatusOk

`func (o *WIREParticipant) GetFundsTransferStatusOk() (*string, bool)`

GetFundsTransferStatusOk returns a tuple with the FundsTransferStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFundsTransferStatus

`func (o *WIREParticipant) SetFundsTransferStatus(v string)`

SetFundsTransferStatus sets FundsTransferStatus field to given value.

### HasFundsTransferStatus

`func (o *WIREParticipant) HasFundsTransferStatus() bool`

HasFundsTransferStatus returns a boolean if a field has been set.

### GetFundsSettlementOnlyStatus

`func (o *WIREParticipant) GetFundsSettlementOnlyStatus() string`

GetFundsSettlementOnlyStatus returns the FundsSettlementOnlyStatus field if non-nil, zero value otherwise.

### GetFundsSettlementOnlyStatusOk

`func (o *WIREParticipant) GetFundsSettlementOnlyStatusOk() (*string, bool)`

GetFundsSettlementOnlyStatusOk returns a tuple with the FundsSettlementOnlyStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFundsSettlementOnlyStatus

`func (o *WIREParticipant) SetFundsSettlementOnlyStatus(v string)`

SetFundsSettlementOnlyStatus sets FundsSettlementOnlyStatus field to given value.

### HasFundsSettlementOnlyStatus

`func (o *WIREParticipant) HasFundsSettlementOnlyStatus() bool`

HasFundsSettlementOnlyStatus returns a boolean if a field has been set.

### GetBookEntrySecuritiesTransferStatus

`func (o *WIREParticipant) GetBookEntrySecuritiesTransferStatus() string`

GetBookEntrySecuritiesTransferStatus returns the BookEntrySecuritiesTransferStatus field if non-nil, zero value otherwise.

### GetBookEntrySecuritiesTransferStatusOk

`func (o *WIREParticipant) GetBookEntrySecuritiesTransferStatusOk() (*string, bool)`

GetBookEntrySecuritiesTransferStatusOk returns a tuple with the BookEntrySecuritiesTransferStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBookEntrySecuritiesTransferStatus

`func (o *WIREParticipant) SetBookEntrySecuritiesTransferStatus(v string)`

SetBookEntrySecuritiesTransferStatus sets BookEntrySecuritiesTransferStatus field to given value.

### HasBookEntrySecuritiesTransferStatus

`func (o *WIREParticipant) HasBookEntrySecuritiesTransferStatus() bool`

HasBookEntrySecuritiesTransferStatus returns a boolean if a field has been set.

### GetDate

`func (o *WIREParticipant) GetDate() string`

GetDate returns the Date field if non-nil, zero value otherwise.

### GetDateOk

`func (o *WIREParticipant) GetDateOk() (*string, bool)`

GetDateOk returns a tuple with the Date field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDate

`func (o *WIREParticipant) SetDate(v string)`

SetDate sets Date field to given value.

### HasDate

`func (o *WIREParticipant) HasDate() bool`

HasDate returns a boolean if a field has been set.

### GetCleanName

`func (o *WIREParticipant) GetCleanName() string`

GetCleanName returns the CleanName field if non-nil, zero value otherwise.

### GetCleanNameOk

`func (o *WIREParticipant) GetCleanNameOk() (*string, bool)`

GetCleanNameOk returns a tuple with the CleanName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCleanName

`func (o *WIREParticipant) SetCleanName(v string)`

SetCleanName sets CleanName field to given value.

### HasCleanName

`func (o *WIREParticipant) HasCleanName() bool`

HasCleanName returns a boolean if a field has been set.

### GetLogo

`func (o *WIREParticipant) GetLogo() Logo`

GetLogo returns the Logo field if non-nil, zero value otherwise.

### GetLogoOk

`func (o *WIREParticipant) GetLogoOk() (*Logo, bool)`

GetLogoOk returns a tuple with the Logo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLogo

`func (o *WIREParticipant) SetLogo(v Logo)`

SetLogo sets Logo field to given value.

### HasLogo

`func (o *WIREParticipant) HasLogo() bool`

HasLogo returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


