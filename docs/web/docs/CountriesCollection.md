# CountriesCollection

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**[]CountryData**](CountryData.md) |  | 
**Links** | [**PaginationData**](PaginationData.md) |  | 

## Methods

### NewCountriesCollection

`func NewCountriesCollection(data []CountryData, links PaginationData, ) *CountriesCollection`

NewCountriesCollection instantiates a new CountriesCollection object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCountriesCollectionWithDefaults

`func NewCountriesCollectionWithDefaults() *CountriesCollection`

NewCountriesCollectionWithDefaults instantiates a new CountriesCollection object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *CountriesCollection) GetData() []CountryData`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *CountriesCollection) GetDataOk() (*[]CountryData, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *CountriesCollection) SetData(v []CountryData)`

SetData sets Data field to given value.


### GetLinks

`func (o *CountriesCollection) GetLinks() PaginationData`

GetLinks returns the Links field if non-nil, zero value otherwise.

### GetLinksOk

`func (o *CountriesCollection) GetLinksOk() (*PaginationData, bool)`

GetLinksOk returns a tuple with the Links field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinks

`func (o *CountriesCollection) SetLinks(v PaginationData)`

SetLinks sets Links field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


