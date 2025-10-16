# CityAdminAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Username** | **string** | username of the user | 
**CityId** | [**uuid.UUID**](uuid.UUID.md) | city id | 
**Role** | **string** | role of the user in this city | 
**Avatar** | Pointer to **string** | optional avatar url for the user in this city | [optional] 
**Label** | Pointer to **string** | optional label for the user in this city | [optional] 
**CreatedAt** | **time.Time** | record creation date | 
**UpdatedAt** | **time.Time** | last update date | 

## Methods

### NewCityAdminAttributes

`func NewCityAdminAttributes(username string, cityId uuid.UUID, role string, createdAt time.Time, updatedAt time.Time, ) *CityAdminAttributes`

NewCityAdminAttributes instantiates a new CityAdminAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCityAdminAttributesWithDefaults

`func NewCityAdminAttributesWithDefaults() *CityAdminAttributes`

NewCityAdminAttributesWithDefaults instantiates a new CityAdminAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUsername

`func (o *CityAdminAttributes) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *CityAdminAttributes) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *CityAdminAttributes) SetUsername(v string)`

SetUsername sets Username field to given value.


### GetCityId

`func (o *CityAdminAttributes) GetCityId() uuid.UUID`

GetCityId returns the CityId field if non-nil, zero value otherwise.

### GetCityIdOk

`func (o *CityAdminAttributes) GetCityIdOk() (*uuid.UUID, bool)`

GetCityIdOk returns a tuple with the CityId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCityId

`func (o *CityAdminAttributes) SetCityId(v uuid.UUID)`

SetCityId sets CityId field to given value.


### GetRole

`func (o *CityAdminAttributes) GetRole() string`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *CityAdminAttributes) GetRoleOk() (*string, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *CityAdminAttributes) SetRole(v string)`

SetRole sets Role field to given value.


### GetAvatar

`func (o *CityAdminAttributes) GetAvatar() string`

GetAvatar returns the Avatar field if non-nil, zero value otherwise.

### GetAvatarOk

`func (o *CityAdminAttributes) GetAvatarOk() (*string, bool)`

GetAvatarOk returns a tuple with the Avatar field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvatar

`func (o *CityAdminAttributes) SetAvatar(v string)`

SetAvatar sets Avatar field to given value.

### HasAvatar

`func (o *CityAdminAttributes) HasAvatar() bool`

HasAvatar returns a boolean if a field has been set.

### GetLabel

`func (o *CityAdminAttributes) GetLabel() string`

GetLabel returns the Label field if non-nil, zero value otherwise.

### GetLabelOk

`func (o *CityAdminAttributes) GetLabelOk() (*string, bool)`

GetLabelOk returns a tuple with the Label field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabel

`func (o *CityAdminAttributes) SetLabel(v string)`

SetLabel sets Label field to given value.

### HasLabel

`func (o *CityAdminAttributes) HasLabel() bool`

HasLabel returns a boolean if a field has been set.

### GetCreatedAt

`func (o *CityAdminAttributes) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *CityAdminAttributes) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *CityAdminAttributes) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetUpdatedAt

`func (o *CityAdminAttributes) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *CityAdminAttributes) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *CityAdminAttributes) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


