# Logo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Company name | [optional] 
**Url** | Pointer to **string** | URL to the company logo | [optional] 

## Methods

### NewLogo

`func NewLogo() *Logo`

NewLogo instantiates a new Logo object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLogoWithDefaults

`func NewLogoWithDefaults() *Logo`

NewLogoWithDefaults instantiates a new Logo object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *Logo) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Logo) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Logo) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *Logo) HasName() bool`

HasName returns a boolean if a field has been set.

### GetUrl

`func (o *Logo) GetUrl() string`

GetUrl returns the Url field if non-nil, zero value otherwise.

### GetUrlOk

`func (o *Logo) GetUrlOk() (*string, bool)`

GetUrlOk returns a tuple with the Url field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrl

`func (o *Logo) SetUrl(v string)`

SetUrl sets Url field to given value.

### HasUrl

`func (o *Logo) HasUrl() bool`

HasUrl returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


