# CreateMayorDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CityId** | [**uuid.UUID**](uuid.UUID.md) | city id | 
**UserId** | [**uuid.UUID**](uuid.UUID.md) | user id | 
**Label** | **string** | label for the user in this city | 

## Methods

### NewCreateMayorDataAttributes

`func NewCreateMayorDataAttributes(cityId uuid.UUID, userId uuid.UUID, label string, ) *CreateMayorDataAttributes`

NewCreateMayorDataAttributes instantiates a new CreateMayorDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateMayorDataAttributesWithDefaults

`func NewCreateMayorDataAttributesWithDefaults() *CreateMayorDataAttributes`

NewCreateMayorDataAttributesWithDefaults instantiates a new CreateMayorDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCityId

`func (o *CreateMayorDataAttributes) GetCityId() uuid.UUID`

GetCityId returns the CityId field if non-nil, zero value otherwise.

### GetCityIdOk

`func (o *CreateMayorDataAttributes) GetCityIdOk() (*uuid.UUID, bool)`

GetCityIdOk returns a tuple with the CityId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCityId

`func (o *CreateMayorDataAttributes) SetCityId(v uuid.UUID)`

SetCityId sets CityId field to given value.


### GetUserId

`func (o *CreateMayorDataAttributes) GetUserId() uuid.UUID`

GetUserId returns the UserId field if non-nil, zero value otherwise.

### GetUserIdOk

`func (o *CreateMayorDataAttributes) GetUserIdOk() (*uuid.UUID, bool)`

GetUserIdOk returns a tuple with the UserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId

`func (o *CreateMayorDataAttributes) SetUserId(v uuid.UUID)`

SetUserId sets UserId field to given value.


### GetLabel

`func (o *CreateMayorDataAttributes) GetLabel() string`

GetLabel returns the Label field if non-nil, zero value otherwise.

### GetLabelOk

`func (o *CreateMayorDataAttributes) GetLabelOk() (*string, bool)`

GetLabelOk returns a tuple with the Label field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabel

`func (o *CreateMayorDataAttributes) SetLabel(v string)`

SetLabel sets Label field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


