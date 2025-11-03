# CreateCityStatusDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | The name of the city. | 
**Description** | **string** | A brief description of the city. | 
**Accessible** | **bool** | Indicates whether the city is accessible to the public. | 
**AllowedAdmin** | **bool** | Indicates whether administrative actions are allowed in the city. | 

## Methods

### NewCreateCityStatusDataAttributes

`func NewCreateCityStatusDataAttributes(name string, description string, accessible bool, allowedAdmin bool, ) *CreateCityStatusDataAttributes`

NewCreateCityStatusDataAttributes instantiates a new CreateCityStatusDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateCityStatusDataAttributesWithDefaults

`func NewCreateCityStatusDataAttributesWithDefaults() *CreateCityStatusDataAttributes`

NewCreateCityStatusDataAttributesWithDefaults instantiates a new CreateCityStatusDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *CreateCityStatusDataAttributes) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CreateCityStatusDataAttributes) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CreateCityStatusDataAttributes) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *CreateCityStatusDataAttributes) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *CreateCityStatusDataAttributes) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *CreateCityStatusDataAttributes) SetDescription(v string)`

SetDescription sets Description field to given value.


### GetAccessible

`func (o *CreateCityStatusDataAttributes) GetAccessible() bool`

GetAccessible returns the Accessible field if non-nil, zero value otherwise.

### GetAccessibleOk

`func (o *CreateCityStatusDataAttributes) GetAccessibleOk() (*bool, bool)`

GetAccessibleOk returns a tuple with the Accessible field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccessible

`func (o *CreateCityStatusDataAttributes) SetAccessible(v bool)`

SetAccessible sets Accessible field to given value.


### GetAllowedAdmin

`func (o *CreateCityStatusDataAttributes) GetAllowedAdmin() bool`

GetAllowedAdmin returns the AllowedAdmin field if non-nil, zero value otherwise.

### GetAllowedAdminOk

`func (o *CreateCityStatusDataAttributes) GetAllowedAdminOk() (*bool, bool)`

GetAllowedAdminOk returns a tuple with the AllowedAdmin field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllowedAdmin

`func (o *CreateCityStatusDataAttributes) SetAllowedAdmin(v bool)`

SetAllowedAdmin sets AllowedAdmin field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


