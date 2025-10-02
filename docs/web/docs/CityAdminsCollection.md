# CityAdminsCollection

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**[]CityAdminData**](CityAdminData.md) |  | 
**Links** | [**PaginationData**](PaginationData.md) |  | 

## Methods

### NewCityAdminsCollection

`func NewCityAdminsCollection(data []CityAdminData, links PaginationData, ) *CityAdminsCollection`

NewCityAdminsCollection instantiates a new CityAdminsCollection object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCityAdminsCollectionWithDefaults

`func NewCityAdminsCollectionWithDefaults() *CityAdminsCollection`

NewCityAdminsCollectionWithDefaults instantiates a new CityAdminsCollection object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *CityAdminsCollection) GetData() []CityAdminData`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *CityAdminsCollection) GetDataOk() (*[]CityAdminData, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *CityAdminsCollection) SetData(v []CityAdminData)`

SetData sets Data field to given value.


### GetLinks

`func (o *CityAdminsCollection) GetLinks() PaginationData`

GetLinks returns the Links field if non-nil, zero value otherwise.

### GetLinksOk

`func (o *CityAdminsCollection) GetLinksOk() (*PaginationData, bool)`

GetLinksOk returns a tuple with the Links field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinks

`func (o *CityAdminsCollection) SetLinks(v PaginationData)`

SetLinks sets Links field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


