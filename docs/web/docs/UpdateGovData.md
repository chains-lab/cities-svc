# UpdateGovData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | government id | 
**Type** | **string** |  | 
**Attributes** | [**UpdateOwnActiveGovDataAttributes**](UpdateOwnActiveGovDataAttributes.md) |  | 

## Methods

### NewUpdateGovData

`func NewUpdateGovData(id string, type_ string, attributes UpdateOwnActiveGovDataAttributes, ) *UpdateGovData`

NewUpdateGovData instantiates a new UpdateGovData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateGovDataWithDefaults

`func NewUpdateGovDataWithDefaults() *UpdateGovData`

NewUpdateGovDataWithDefaults instantiates a new UpdateGovData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *UpdateGovData) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *UpdateGovData) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *UpdateGovData) SetId(v string)`

SetId sets Id field to given value.


### GetType

`func (o *UpdateGovData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *UpdateGovData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *UpdateGovData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *UpdateGovData) GetAttributes() UpdateOwnActiveGovDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *UpdateGovData) GetAttributesOk() (*UpdateOwnActiveGovDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *UpdateGovData) SetAttributes(v UpdateOwnActiveGovDataAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


