# プラグインAPI仕様書
<a name="top"></a>

## インデックス
- [API仕様](#API仕様)

  - [organization.proto](#organization.proto)
      - [Organization](#ncs.protobuf.Organization)
  

  - [permission.proto](#permission.proto)
      - [Permission](#ncs.protobuf.Permission)
  

  - [resource.proto](#resource.proto)
      - [authzRequest](#ncs.protobuf.authzRequest)
      - [saveRequest](#ncs.protobuf.saveRequest)
      - [Resource](#ncs.protobuf.Resource)
  

  - [role.proto](#role.proto)
      - [Role](#ncs.protobuf.Role)
  

  - [types.proto](#types.proto)
      - [empty](#ncs.protobuf.empty)
      - [organizationEntities](#ncs.protobuf.organizationEntities)
      - [organizationEntity](#ncs.protobuf.organizationEntity)
      - [organizationKey](#ncs.protobuf.organizationKey)
      - [organizationUpdateEntity](#ncs.protobuf.organizationUpdateEntity)
      - [organizationUser](#ncs.protobuf.organizationUser)
      - [permissionCheckRequest](#ncs.protobuf.permissionCheckRequest)
      - [permissionCheckResult](#ncs.protobuf.permissionCheckResult)
      - [permissionEntities](#ncs.protobuf.permissionEntities)
      - [permissionEntity](#ncs.protobuf.permissionEntity)
      - [permissionKey](#ncs.protobuf.permissionKey)
      - [roleEntities](#ncs.protobuf.roleEntities)
      - [roleEntity](#ncs.protobuf.roleEntity)
      - [roleKey](#ncs.protobuf.roleKey)
      - [roleReleationPermissions](#ncs.protobuf.roleReleationPermissions)
      - [roleUpdateEntity](#ncs.protobuf.roleUpdateEntity)
      - [userEntity](#ncs.protobuf.userEntity)
      - [userKey](#ncs.protobuf.userKey)
      - [userRole](#ncs.protobuf.userRole)
  

  - [user.proto](#user.proto)
      - [User](#ncs.protobuf.User)
  

- [スカラー値型](#スカラー値型)

## API仕様


<a name="organization.proto"></a>
<p align="right"><a href="#top">Top</a></p>

### organization.proto


 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ncs.protobuf.Organization"></a>

#### Organization


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Create | [organizationEntity](#ncs.protobuf.organizationEntity) | [organizationEntity](#ncs.protobuf.organizationEntity) | Create is create orgnization |
| FindById | [organizationKey](#ncs.protobuf.organizationKey) | [organizationEntity](#ncs.protobuf.organizationEntity) | FindById is application id and organization id |
| FindAll | [empty](#ncs.protobuf.empty) | [organizationEntities](#ncs.protobuf.organizationEntities) | FindAll is application is return organizations |
| Update | [organizationUpdateEntity](#ncs.protobuf.organizationUpdateEntity) | [empty](#ncs.protobuf.empty) | Update is organization entity update |
| Delete | [organizationKey](#ncs.protobuf.organizationKey) | [empty](#ncs.protobuf.empty) | Delete is organization entity delete |

 <!-- end services -->



<a name="permission.proto"></a>
<p align="right"><a href="#top">Top</a></p>

### permission.proto


 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ncs.protobuf.Permission"></a>

#### Permission


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Create | [permissionEntities](#ncs.protobuf.permissionEntities) | [permissionEntities](#ncs.protobuf.permissionEntities) | Create is create permission |
| FindById | [permissionKey](#ncs.protobuf.permissionKey) | [permissionEntity](#ncs.protobuf.permissionEntity) | FindById is find by id |
| FindAll | [empty](#ncs.protobuf.empty) | [permissionEntities](#ncs.protobuf.permissionEntities) | FindAll is find by application id return permissions |
| Update | [permissionEntity](#ncs.protobuf.permissionEntity) | [empty](#ncs.protobuf.empty) | Update is permission entity update |
| Delete | [permissionKey](#ncs.protobuf.permissionKey) | [empty](#ncs.protobuf.empty) | Delete is permission entity delete |
| Check | [permissionCheckRequest](#ncs.protobuf.permissionCheckRequest) | [permissionCheckResult](#ncs.protobuf.permissionCheckResult) | Check is resource check |

 <!-- end services -->



<a name="resource.proto"></a>
<p align="right"><a href="#top">Top</a></p>

### resource.proto



<a name="ncs.protobuf.authzRequest"></a>

#### authzRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| method | [string](#string) |  |  |
| uri | [string](#string) |  |  |
| organizationName | [string](#string) |  |  |
| userKey | [string](#string) |  |  |






<a name="ncs.protobuf.saveRequest"></a>

#### saveRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| method | [string](#string) |  |  |
| uri | [string](#string) |  |  |
| permissions | [string](#string) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ncs.protobuf.Resource"></a>

#### Resource


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Save | [saveRequest](#ncs.protobuf.saveRequest) | [empty](#ncs.protobuf.empty) |  |
| Authz | [authzRequest](#ncs.protobuf.authzRequest) | [empty](#ncs.protobuf.empty) |  |

 <!-- end services -->



<a name="role.proto"></a>
<p align="right"><a href="#top">Top</a></p>

### role.proto


 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ncs.protobuf.Role"></a>

#### Role


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Create | [roleEntities](#ncs.protobuf.roleEntities) | [roleEntities](#ncs.protobuf.roleEntities) | RoleCreate is create role |
| FindById | [roleKey](#ncs.protobuf.roleKey) | [roleEntity](#ncs.protobuf.roleEntity) | FindById is find by id |
| FindAll | [empty](#ncs.protobuf.empty) | [roleEntities](#ncs.protobuf.roleEntities) | FindAll is find roles |
| Update | [roleUpdateEntity](#ncs.protobuf.roleUpdateEntity) | [empty](#ncs.protobuf.empty) | Update is role entity update |
| Delete | [roleKey](#ncs.protobuf.roleKey) | [empty](#ncs.protobuf.empty) | Delete is role entity delete |
| GetPermissions | [roleKey](#ncs.protobuf.roleKey) | [permissionEntities](#ncs.protobuf.permissionEntities) | GetPermissions is get permission to the role |
| AddPermissions | [roleReleationPermissions](#ncs.protobuf.roleReleationPermissions) | [empty](#ncs.protobuf.empty) | AddPermissions is add permission to the role |
| DeletePermission | [roleReleationPermissions](#ncs.protobuf.roleReleationPermissions) | [empty](#ncs.protobuf.empty) | DeletePermission is delete permission to the role |

 <!-- end services -->



<a name="types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

### types.proto



<a name="ncs.protobuf.empty"></a>

#### empty







<a name="ncs.protobuf.organizationEntities"></a>

#### organizationEntities



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organizations | [organizationEntity](#ncs.protobuf.organizationEntity) | repeated |  |






<a name="ncs.protobuf.organizationEntity"></a>

#### organizationEntity



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) |  |  |
| users | [userEntity](#ncs.protobuf.userEntity) | repeated |  |






<a name="ncs.protobuf.organizationKey"></a>

#### organizationKey



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="ncs.protobuf.organizationUpdateEntity"></a>

#### organizationUpdateEntity



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) |  |  |






<a name="ncs.protobuf.organizationUser"></a>

#### organizationUser



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_key | [string](#string) |  |  |
| organization_id | [string](#string) |  |  |
| organization_name | [string](#string) |  |  |
| organization_description | [string](#string) |  |  |






<a name="ncs.protobuf.permissionCheckRequest"></a>

#### permissionCheckRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| userKey | [string](#string) |  |  |
| organizationName | [string](#string) |  |  |
| permissionNames | [string](#string) | repeated |  |






<a name="ncs.protobuf.permissionCheckResult"></a>

#### permissionCheckResult



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| result | [bool](#bool) |  |  |
| message | [string](#string) |  |  |






<a name="ncs.protobuf.permissionEntities"></a>

#### permissionEntities



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| permissions | [permissionEntity](#ncs.protobuf.permissionEntity) | repeated |  |






<a name="ncs.protobuf.permissionEntity"></a>

#### permissionEntity



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) |  |  |






<a name="ncs.protobuf.permissionKey"></a>

#### permissionKey



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="ncs.protobuf.roleEntities"></a>

#### roleEntities



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| roles | [roleEntity](#ncs.protobuf.roleEntity) | repeated |  |






<a name="ncs.protobuf.roleEntity"></a>

#### roleEntity



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) |  |  |
| permissions | [permissionEntity](#ncs.protobuf.permissionEntity) | repeated |  |
| organizationUsers | [organizationUser](#ncs.protobuf.organizationUser) | repeated |  |






<a name="ncs.protobuf.roleKey"></a>

#### roleKey



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="ncs.protobuf.roleReleationPermissions"></a>

#### roleReleationPermissions



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| permissions | [permissionKey](#ncs.protobuf.permissionKey) | repeated |  |






<a name="ncs.protobuf.roleUpdateEntity"></a>

#### roleUpdateEntity



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) |  |  |






<a name="ncs.protobuf.userEntity"></a>

#### userEntity



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| organization_id | [string](#string) |  |  |
| roles | [roleEntity](#ncs.protobuf.roleEntity) | repeated |  |
| permissions | [permissionEntity](#ncs.protobuf.permissionEntity) | repeated |  |






<a name="ncs.protobuf.userKey"></a>

#### userKey



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| organization_id | [string](#string) |  |  |






<a name="ncs.protobuf.userRole"></a>

#### userRole



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [userKey](#ncs.protobuf.userKey) |  |  |
| roles | [roleKey](#ncs.protobuf.roleKey) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="user.proto"></a>
<p align="right"><a href="#top">Top</a></p>

### user.proto


 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ncs.protobuf.User"></a>

#### User


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Create | [userEntity](#ncs.protobuf.userEntity) | [empty](#ncs.protobuf.empty) | Create is create user |
| Delete | [userKey](#ncs.protobuf.userKey) | [empty](#ncs.protobuf.empty) | Delete is delete user |
| FindByKey | [userKey](#ncs.protobuf.userKey) | [userEntity](#ncs.protobuf.userEntity) | FindByKey is find organization id and user key |
| AddRole | [userRole](#ncs.protobuf.userRole) | [empty](#ncs.protobuf.empty) | AddRole is add role to user |
| DeleteRole | [userRole](#ncs.protobuf.userRole) | [empty](#ncs.protobuf.empty) | DeleteRole is add role to user |

 <!-- end services -->



## スカラー値型

| .proto Type | Notes | Go Type | C++ Type | Java Type | Python Type |
| ----------- | ----- | -------- | -------- | --------- | ----------- |
| <a name="double" /> double |  | float64 | double | double | float |
| <a name="float" /> float |  | float32 | float | float | float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int32 | int | int |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | int64 | long | int/long |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | uint32 | int | int/long |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | uint64 | long | int/long |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int32 | int | int |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | int64 | long | int/long |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | uint32 | int | int |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | uint64 | long | int/long |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int32 | int | int |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | int64 | long | int/long |
| <a name="bool" /> bool |  | bool | bool | boolean | boolean |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | string | String | str/unicode |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | []byte | string | ByteString | str |
