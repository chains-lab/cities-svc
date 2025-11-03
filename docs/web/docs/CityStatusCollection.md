# CityStatusCollection

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**[]CityStatusData**](CityStatusData.md) |  | 
**Links** | [**PaginationData**](PaginationData.md) |  | 

## Methods

### NewCityStatusCollection

`func NewCityStatusCollection(data []CityStatusData, links PaginationData, ) *CityStatusCollection`

NewCityStatusCollection instantiates a new CityStatusCollection object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCityStatusCollectionWithDefaults

`func NewCityStatusCollectionWithDefaults() *CityStatusCollection`

NewCityStatusCollectionWithDefaults instantiates a new CityStatusCollection object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *CityStatusCollection) GetData() []CityStatusData`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *CityStatusCollection) GetDataOk() (*[]CityStatusData, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *CityStatusCollection) SetData(v []CityStatusData)`

SetData sets Data field to given value.


### GetLinks

`func (o *CityStatusCollection) GetLinks() PaginationData`

GetLinks returns the Links field if non-nil, zero value otherwise.

### GetLinksOk

`func (o *CityStatusCollection) GetLinksOk() (*PaginationData, bool)`

GetLinksOk returns a tuple with the Links field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinks

`func (o *CityStatusCollection) SetLinks(v PaginationData)`

SetLinks sets Links field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


