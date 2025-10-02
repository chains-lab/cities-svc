# UpdateOwnActiveGovData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**uuid.UUID**](uuid.UUID.md) | user id | 
**Type** | **string** |  | 
**Attributes** | [**UpdateOwnActiveGovDataAttributes**](UpdateOwnActiveGovDataAttributes.md) |  | 

## Methods

### NewUpdateOwnActiveGovData

`func NewUpdateOwnActiveGovData(id uuid.UUID, type_ string, attributes UpdateOwnActiveGovDataAttributes, ) *UpdateOwnActiveGovData`

NewUpdateOwnActiveGovData instantiates a new UpdateOwnActiveGovData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateOwnActiveGovDataWithDefaults

`func NewUpdateOwnActiveGovDataWithDefaults() *UpdateOwnActiveGovData`

NewUpdateOwnActiveGovDataWithDefaults instantiates a new UpdateOwnActiveGovData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *UpdateOwnActiveGovData) GetId() uuid.UUID`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *UpdateOwnActiveGovData) GetIdOk() (*uuid.UUID, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *UpdateOwnActiveGovData) SetId(v uuid.UUID)`

SetId sets Id field to given value.


### GetType

`func (o *UpdateOwnActiveGovData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *UpdateOwnActiveGovData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *UpdateOwnActiveGovData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *UpdateOwnActiveGovData) GetAttributes() UpdateOwnActiveGovDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *UpdateOwnActiveGovData) GetAttributesOk() (*UpdateOwnActiveGovDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *UpdateOwnActiveGovData) SetAttributes(v UpdateOwnActiveGovDataAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


