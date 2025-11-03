# UpdateOwnAdminData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**uuid.UUID**](uuid.UUID.md) | user id | 
**Type** | **string** |  | 
**Attributes** | [**UpdateOwnAdminDataAttributes**](UpdateOwnAdminDataAttributes.md) |  | 

## Methods

### NewUpdateOwnAdminData

`func NewUpdateOwnAdminData(id uuid.UUID, type_ string, attributes UpdateOwnAdminDataAttributes, ) *UpdateOwnAdminData`

NewUpdateOwnAdminData instantiates a new UpdateOwnAdminData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateOwnAdminDataWithDefaults

`func NewUpdateOwnAdminDataWithDefaults() *UpdateOwnAdminData`

NewUpdateOwnAdminDataWithDefaults instantiates a new UpdateOwnAdminData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *UpdateOwnAdminData) GetId() uuid.UUID`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *UpdateOwnAdminData) GetIdOk() (*uuid.UUID, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *UpdateOwnAdminData) SetId(v uuid.UUID)`

SetId sets Id field to given value.


### GetType

`func (o *UpdateOwnAdminData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *UpdateOwnAdminData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *UpdateOwnAdminData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *UpdateOwnAdminData) GetAttributes() UpdateOwnAdminDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *UpdateOwnAdminData) GetAttributesOk() (*UpdateOwnAdminDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *UpdateOwnAdminData) SetAttributes(v UpdateOwnAdminDataAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


