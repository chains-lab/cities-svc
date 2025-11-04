# CityStatusDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | name of the city status | 
**Description** | **string** | description of the city status | 
**Accessible** | **bool** | whether the city status is accessible to users | 
**AllowedAdmins** | **int32** | number of allowed admins for the city status | 
**CreatedAt** | **time.Time** |  | 

## Methods

### NewCityStatusDataAttributes

`func NewCityStatusDataAttributes(name string, description string, accessible bool, allowedAdmins int32, createdAt time.Time, ) *CityStatusDataAttributes`

NewCityStatusDataAttributes instantiates a new CityStatusDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCityStatusDataAttributesWithDefaults

`func NewCityStatusDataAttributesWithDefaults() *CityStatusDataAttributes`

NewCityStatusDataAttributesWithDefaults instantiates a new CityStatusDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *CityStatusDataAttributes) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CityStatusDataAttributes) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CityStatusDataAttributes) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *CityStatusDataAttributes) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *CityStatusDataAttributes) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *CityStatusDataAttributes) SetDescription(v string)`

SetDescription sets Description field to given value.


### GetAccessible

`func (o *CityStatusDataAttributes) GetAccessible() bool`

GetAccessible returns the Accessible field if non-nil, zero value otherwise.

### GetAccessibleOk

`func (o *CityStatusDataAttributes) GetAccessibleOk() (*bool, bool)`

GetAccessibleOk returns a tuple with the Accessible field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccessible

`func (o *CityStatusDataAttributes) SetAccessible(v bool)`

SetAccessible sets Accessible field to given value.


### GetAllowedAdmins

`func (o *CityStatusDataAttributes) GetAllowedAdmins() int32`

GetAllowedAdmins returns the AllowedAdmins field if non-nil, zero value otherwise.

### GetAllowedAdminsOk

`func (o *CityStatusDataAttributes) GetAllowedAdminsOk() (*int32, bool)`

GetAllowedAdminsOk returns a tuple with the AllowedAdmins field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllowedAdmins

`func (o *CityStatusDataAttributes) SetAllowedAdmins(v int32)`

SetAllowedAdmins sets AllowedAdmins field to given value.


### GetCreatedAt

`func (o *CityStatusDataAttributes) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *CityStatusDataAttributes) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *CityStatusDataAttributes) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


