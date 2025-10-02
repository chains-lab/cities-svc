# CreateInviteDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**UserId** | [**uuid.UUID**](uuid.UUID.md) | user id | 
**Role** | **string** | role of the user in this city | 
**Label** | **string** | label for the user in this city | 

## Methods

### NewCreateInviteDataAttributes

`func NewCreateInviteDataAttributes(userId uuid.UUID, role string, label string, ) *CreateInviteDataAttributes`

NewCreateInviteDataAttributes instantiates a new CreateInviteDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateInviteDataAttributesWithDefaults

`func NewCreateInviteDataAttributesWithDefaults() *CreateInviteDataAttributes`

NewCreateInviteDataAttributesWithDefaults instantiates a new CreateInviteDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUserId

`func (o *CreateInviteDataAttributes) GetUserId() uuid.UUID`

GetUserId returns the UserId field if non-nil, zero value otherwise.

### GetUserIdOk

`func (o *CreateInviteDataAttributes) GetUserIdOk() (*uuid.UUID, bool)`

GetUserIdOk returns a tuple with the UserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId

`func (o *CreateInviteDataAttributes) SetUserId(v uuid.UUID)`

SetUserId sets UserId field to given value.


### GetRole

`func (o *CreateInviteDataAttributes) GetRole() string`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *CreateInviteDataAttributes) GetRoleOk() (*string, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *CreateInviteDataAttributes) SetRole(v string)`

SetRole sets Role field to given value.


### GetLabel

`func (o *CreateInviteDataAttributes) GetLabel() string`

GetLabel returns the Label field if non-nil, zero value otherwise.

### GetLabelOk

`func (o *CreateInviteDataAttributes) GetLabelOk() (*string, bool)`

GetLabelOk returns a tuple with the Label field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabel

`func (o *CreateInviteDataAttributes) SetLabel(v string)`

SetLabel sets Label field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


