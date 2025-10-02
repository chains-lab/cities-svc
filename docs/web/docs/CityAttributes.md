# CityAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CountryId** | [**uuid.UUID**](uuid.UUID.md) | country id | 
**Point** | [**Point**](Point.md) |  | 
**Status** | **string** | city status | 
**Name** | **string** | city name | 
**Icon** | Pointer to **string** | city icon uri | [optional] 
**Slug** | Pointer to **string** | city slug | [optional] 
**Timezone** | **string** | city timezone | 
**CreatedAt** | **time.Time** | creation date | 
**UpdatedAt** | **time.Time** | last update date | 

## Methods

### NewCityAttributes

`func NewCityAttributes(countryId uuid.UUID, point Point, status string, name string, timezone string, createdAt time.Time, updatedAt time.Time, ) *CityAttributes`

NewCityAttributes instantiates a new CityAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCityAttributesWithDefaults

`func NewCityAttributesWithDefaults() *CityAttributes`

NewCityAttributesWithDefaults instantiates a new CityAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCountryId

`func (o *CityAttributes) GetCountryId() uuid.UUID`

GetCountryId returns the CountryId field if non-nil, zero value otherwise.

### GetCountryIdOk

`func (o *CityAttributes) GetCountryIdOk() (*uuid.UUID, bool)`

GetCountryIdOk returns a tuple with the CountryId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCountryId

`func (o *CityAttributes) SetCountryId(v uuid.UUID)`

SetCountryId sets CountryId field to given value.


### GetPoint

`func (o *CityAttributes) GetPoint() Point`

GetPoint returns the Point field if non-nil, zero value otherwise.

### GetPointOk

`func (o *CityAttributes) GetPointOk() (*Point, bool)`

GetPointOk returns a tuple with the Point field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoint

`func (o *CityAttributes) SetPoint(v Point)`

SetPoint sets Point field to given value.


### GetStatus

`func (o *CityAttributes) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *CityAttributes) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *CityAttributes) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetName

`func (o *CityAttributes) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CityAttributes) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CityAttributes) SetName(v string)`

SetName sets Name field to given value.


### GetIcon

`func (o *CityAttributes) GetIcon() string`

GetIcon returns the Icon field if non-nil, zero value otherwise.

### GetIconOk

`func (o *CityAttributes) GetIconOk() (*string, bool)`

GetIconOk returns a tuple with the Icon field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIcon

`func (o *CityAttributes) SetIcon(v string)`

SetIcon sets Icon field to given value.

### HasIcon

`func (o *CityAttributes) HasIcon() bool`

HasIcon returns a boolean if a field has been set.

### GetSlug

`func (o *CityAttributes) GetSlug() string`

GetSlug returns the Slug field if non-nil, zero value otherwise.

### GetSlugOk

`func (o *CityAttributes) GetSlugOk() (*string, bool)`

GetSlugOk returns a tuple with the Slug field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSlug

`func (o *CityAttributes) SetSlug(v string)`

SetSlug sets Slug field to given value.

### HasSlug

`func (o *CityAttributes) HasSlug() bool`

HasSlug returns a boolean if a field has been set.

### GetTimezone

`func (o *CityAttributes) GetTimezone() string`

GetTimezone returns the Timezone field if non-nil, zero value otherwise.

### GetTimezoneOk

`func (o *CityAttributes) GetTimezoneOk() (*string, bool)`

GetTimezoneOk returns a tuple with the Timezone field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimezone

`func (o *CityAttributes) SetTimezone(v string)`

SetTimezone sets Timezone field to given value.


### GetCreatedAt

`func (o *CityAttributes) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *CityAttributes) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *CityAttributes) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetUpdatedAt

`func (o *CityAttributes) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *CityAttributes) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *CityAttributes) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


