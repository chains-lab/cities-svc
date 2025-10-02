# TransferMayorDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**UserId** | [**uuid.UUID**](uuid.UUID.md) | ID of the user to be assigned as the new mayor | 

## Methods

### NewTransferMayorDataAttributes

`func NewTransferMayorDataAttributes(userId uuid.UUID, ) *TransferMayorDataAttributes`

NewTransferMayorDataAttributes instantiates a new TransferMayorDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTransferMayorDataAttributesWithDefaults

`func NewTransferMayorDataAttributesWithDefaults() *TransferMayorDataAttributes`

NewTransferMayorDataAttributesWithDefaults instantiates a new TransferMayorDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUserId

`func (o *TransferMayorDataAttributes) GetUserId() uuid.UUID`

GetUserId returns the UserId field if non-nil, zero value otherwise.

### GetUserIdOk

`func (o *TransferMayorDataAttributes) GetUserIdOk() (*uuid.UUID, bool)`

GetUserIdOk returns a tuple with the UserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId

`func (o *TransferMayorDataAttributes) SetUserId(v uuid.UUID)`

SetUserId sets UserId field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


