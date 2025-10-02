# SentInviteDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CityId** | [**uuid.UUID**](uuid.UUID.md) | ID of the city the invite is for | 
**Role** | **string** | Role assigned to the invited user | 

## Methods

### NewSentInviteDataAttributes

`func NewSentInviteDataAttributes(cityId uuid.UUID, role string, ) *SentInviteDataAttributes`

NewSentInviteDataAttributes instantiates a new SentInviteDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSentInviteDataAttributesWithDefaults

`func NewSentInviteDataAttributesWithDefaults() *SentInviteDataAttributes`

NewSentInviteDataAttributesWithDefaults instantiates a new SentInviteDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCityId

`func (o *SentInviteDataAttributes) GetCityId() uuid.UUID`

GetCityId returns the CityId field if non-nil, zero value otherwise.

### GetCityIdOk

`func (o *SentInviteDataAttributes) GetCityIdOk() (*uuid.UUID, bool)`

GetCityIdOk returns a tuple with the CityId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCityId

`func (o *SentInviteDataAttributes) SetCityId(v uuid.UUID)`

SetCityId sets CityId field to given value.


### GetRole

`func (o *SentInviteDataAttributes) GetRole() string`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *SentInviteDataAttributes) GetRoleOk() (*string, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *SentInviteDataAttributes) SetRole(v string)`

SetRole sets Role field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


