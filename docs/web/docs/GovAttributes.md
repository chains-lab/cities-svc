# GovAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CityId** | [**uuid.UUID**](uuid.UUID.md) | city id | 
**Role** | **string** | role of the user in this city | 
**Label** | Pointer to **string** | optional label for the user in this city | [optional] 
**CreatedAt** | **time.Time** | record creation date | 
**UpdatedAt** | **time.Time** | last update date | 

## Methods

### NewGovAttributes

`func NewGovAttributes(cityId uuid.UUID, role string, createdAt time.Time, updatedAt time.Time, ) *GovAttributes`

NewGovAttributes instantiates a new GovAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGovAttributesWithDefaults

`func NewGovAttributesWithDefaults() *GovAttributes`

NewGovAttributesWithDefaults instantiates a new GovAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCityId

`func (o *GovAttributes) GetCityId() uuid.UUID`

GetCityId returns the CityId field if non-nil, zero value otherwise.

### GetCityIdOk

`func (o *GovAttributes) GetCityIdOk() (*uuid.UUID, bool)`

GetCityIdOk returns a tuple with the CityId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCityId

`func (o *GovAttributes) SetCityId(v uuid.UUID)`

SetCityId sets CityId field to given value.


### GetRole

`func (o *GovAttributes) GetRole() string`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *GovAttributes) GetRoleOk() (*string, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *GovAttributes) SetRole(v string)`

SetRole sets Role field to given value.


### GetLabel

`func (o *GovAttributes) GetLabel() string`

GetLabel returns the Label field if non-nil, zero value otherwise.

### GetLabelOk

`func (o *GovAttributes) GetLabelOk() (*string, bool)`

GetLabelOk returns a tuple with the Label field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabel

`func (o *GovAttributes) SetLabel(v string)`

SetLabel sets Label field to given value.

### HasLabel

`func (o *GovAttributes) HasLabel() bool`

HasLabel returns a boolean if a field has been set.

### GetCreatedAt

`func (o *GovAttributes) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *GovAttributes) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *GovAttributes) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetUpdatedAt

`func (o *GovAttributes) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *GovAttributes) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *GovAttributes) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


