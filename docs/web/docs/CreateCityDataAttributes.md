# CreateCityDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CountryId** | [**uuid.UUID**](uuid.UUID.md) | country id | 
**Name** | **string** | city name | 
**Point** | [**Point**](Point.md) |  | 
**Timezone** | **string** | city timezone | 

## Methods

### NewCreateCityDataAttributes

`func NewCreateCityDataAttributes(countryId uuid.UUID, name string, point Point, timezone string, ) *CreateCityDataAttributes`

NewCreateCityDataAttributes instantiates a new CreateCityDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateCityDataAttributesWithDefaults

`func NewCreateCityDataAttributesWithDefaults() *CreateCityDataAttributes`

NewCreateCityDataAttributesWithDefaults instantiates a new CreateCityDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCountryId

`func (o *CreateCityDataAttributes) GetCountryId() uuid.UUID`

GetCountryId returns the CountryId field if non-nil, zero value otherwise.

### GetCountryIdOk

`func (o *CreateCityDataAttributes) GetCountryIdOk() (*uuid.UUID, bool)`

GetCountryIdOk returns a tuple with the CountryId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCountryId

`func (o *CreateCityDataAttributes) SetCountryId(v uuid.UUID)`

SetCountryId sets CountryId field to given value.


### GetName

`func (o *CreateCityDataAttributes) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CreateCityDataAttributes) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CreateCityDataAttributes) SetName(v string)`

SetName sets Name field to given value.


### GetPoint

`func (o *CreateCityDataAttributes) GetPoint() Point`

GetPoint returns the Point field if non-nil, zero value otherwise.

### GetPointOk

`func (o *CreateCityDataAttributes) GetPointOk() (*Point, bool)`

GetPointOk returns a tuple with the Point field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoint

`func (o *CreateCityDataAttributes) SetPoint(v Point)`

SetPoint sets Point field to given value.


### GetTimezone

`func (o *CreateCityDataAttributes) GetTimezone() string`

GetTimezone returns the Timezone field if non-nil, zero value otherwise.

### GetTimezoneOk

`func (o *CreateCityDataAttributes) GetTimezoneOk() (*string, bool)`

GetTimezoneOk returns a tuple with the Timezone field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimezone

`func (o *CreateCityDataAttributes) SetTimezone(v string)`

SetTimezone sets Timezone field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


