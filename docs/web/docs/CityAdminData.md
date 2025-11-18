# CityAdminData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | user id + city id (UUID:UUID) | 
**Type** | **string** |  | 
**Attributes** | [**CityAdminAttributes**](CityAdminAttributes.md) |  | 

## Methods

### NewCityAdminData

`func NewCityAdminData(id string, type_ string, attributes CityAdminAttributes, ) *CityAdminData`

NewCityAdminData instantiates a new CityAdminData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCityAdminDataWithDefaults

`func NewCityAdminDataWithDefaults() *CityAdminData`

NewCityAdminDataWithDefaults instantiates a new CityAdminData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *CityAdminData) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *CityAdminData) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *CityAdminData) SetId(v string)`

SetId sets Id field to given value.


### GetType

`func (o *CityAdminData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *CityAdminData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *CityAdminData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *CityAdminData) GetAttributes() CityAdminAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *CityAdminData) GetAttributesOk() (*CityAdminAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *CityAdminData) SetAttributes(v CityAdminAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


