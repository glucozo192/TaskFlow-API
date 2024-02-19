# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [common.proto](#common-proto)
    - [User](#pb-User)
  
- [product.proto](#product-proto)
    - [DeleteProductRequest](#pb-DeleteProductRequest)
    - [DeleteProductResponse](#pb-DeleteProductResponse)
    - [InsertProductRequest](#pb-InsertProductRequest)
    - [InsertProductResponse](#pb-InsertProductResponse)
    - [ListAllProductsRequest](#pb-ListAllProductsRequest)
    - [ListAllProductsResponse](#pb-ListAllProductsResponse)
    - [ListProductsRequest](#pb-ListProductsRequest)
    - [ListProductsResponse](#pb-ListProductsResponse)
    - [Product](#pb-Product)
  
    - [ProductService](#pb-ProductService)
  
- [user.proto](#user-proto)
    - [GetMeRequest](#pb-GetMeRequest)
    - [GetMeResponse](#pb-GetMeResponse)
    - [GetUserByIDRequest](#pb-GetUserByIDRequest)
    - [GetUserByIDResponse](#pb-GetUserByIDResponse)
    - [LoginRequest](#pb-LoginRequest)
    - [LoginResponse](#pb-LoginResponse)
    - [RegisterRequest](#pb-RegisterRequest)
    - [RegisterResponse](#pb-RegisterResponse)
  
- [user_service.proto](#user_service-proto)
    - [UserService](#pb-UserService)
  
- [transform_options/annotations.proto](#transform_options_annotations-proto)
    - [File-level Extensions](#transform_options_annotations-proto-extensions)
    - [File-level Extensions](#transform_options_annotations-proto-extensions)
    - [File-level Extensions](#transform_options_annotations-proto-extensions)
    - [File-level Extensions](#transform_options_annotations-proto-extensions)
    - [File-level Extensions](#transform_options_annotations-proto-extensions)
    - [File-level Extensions](#transform_options_annotations-proto-extensions)
    - [File-level Extensions](#transform_options_annotations-proto-extensions)
    - [File-level Extensions](#transform_options_annotations-proto-extensions)
    - [File-level Extensions](#transform_options_annotations-proto-extensions)
  
- [Scalar Value Types](#scalar-value-types)



<a name="common-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## common.proto



<a name="pb-User"></a>

### User



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| created_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |
| email | [string](#string) |  |  |
| user_id | [string](#string) |  |  |
| first_name | [string](#string) |  |  |
| last_name | [string](#string) |  |  |
| password | [string](#string) |  |  |





 

 

 

 



<a name="product-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## product.proto



<a name="pb-DeleteProductRequest"></a>

### DeleteProductRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="pb-DeleteProductResponse"></a>

### DeleteProductResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| message | [string](#string) |  |  |






<a name="pb-InsertProductRequest"></a>

### InsertProductRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| product | [Product](#pb-Product) |  |  |






<a name="pb-InsertProductResponse"></a>

### InsertProductResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| product | [Product](#pb-Product) |  |  |






<a name="pb-ListAllProductsRequest"></a>

### ListAllProductsRequest







<a name="pb-ListAllProductsResponse"></a>

### ListAllProductsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| products | [Product](#pb-Product) | repeated |  |






<a name="pb-ListProductsRequest"></a>

### ListProductsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| product_id | [string](#string) |  |  |






<a name="pb-ListProductsResponse"></a>

### ListProductsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| product | [Product](#pb-Product) |  |  |






<a name="pb-Product"></a>

### Product



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| sku | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) |  |  |
| product_status_id | [string](#string) |  |  |
| regular_price | [int32](#int32) |  |  |
| discount_price | [int32](#int32) |  |  |
| quantity | [int32](#int32) |  |  |
| taxable | [bool](#bool) |  |  |
| created_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |





 

 

 


<a name="pb-ProductService"></a>

### ProductService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| InsertProduct | [InsertProductRequest](#pb-InsertProductRequest) | [InsertProductResponse](#pb-InsertProductResponse) |  |
| ListProducts | [ListProductsRequest](#pb-ListProductsRequest) | [ListProductsResponse](#pb-ListProductsResponse) |  |
| ListAllProducts | [ListAllProductsRequest](#pb-ListAllProductsRequest) | [ListAllProductsResponse](#pb-ListAllProductsResponse) |  |
| DeleteProduct | [DeleteProductRequest](#pb-DeleteProductRequest) | [DeleteProductResponse](#pb-DeleteProductResponse) |  |

 



<a name="user-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## user.proto



<a name="pb-GetMeRequest"></a>

### GetMeRequest







<a name="pb-GetMeResponse"></a>

### GetMeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [User](#pb-User) |  |  |






<a name="pb-GetUserByIDRequest"></a>

### GetUserByIDRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [string](#string) |  |  |






<a name="pb-GetUserByIDResponse"></a>

### GetUserByIDResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [User](#pb-User) |  |  |






<a name="pb-LoginRequest"></a>

### LoginRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| email | [string](#string) |  |  |
| password | [string](#string) |  |  |






<a name="pb-LoginResponse"></a>

### LoginResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [string](#string) |  |  |
| access_token | [string](#string) |  |  |






<a name="pb-RegisterRequest"></a>

### RegisterRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [User](#pb-User) |  |  |






<a name="pb-RegisterResponse"></a>

### RegisterResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [string](#string) |  |  |
| access_token | [string](#string) |  |  |





 

 

 

 



<a name="user_service-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## user_service.proto


 

 

 


<a name="pb-UserService"></a>

### UserService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Login | [LoginRequest](#pb-LoginRequest) | [LoginResponse](#pb-LoginResponse) |  |
| Register | [RegisterRequest](#pb-RegisterRequest) | [RegisterResponse](#pb-RegisterResponse) |  |
| GetUserByID | [GetUserByIDRequest](#pb-GetUserByIDRequest) | [GetUserByIDResponse](#pb-GetUserByIDResponse) |  |

 



<a name="transform_options_annotations-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## transform_options/annotations.proto


 

 


<a name="transform_options_annotations-proto-extensions"></a>

### File-level Extensions
| Extension | Type | Base | Number | Description |
| --------- | ---- | ---- | ------ | ----------- |
| custom | bool | .google.protobuf.FieldOptions | 5305 | If true, the custom transformer will be used for the field. |
| embed | bool | .google.protobuf.FieldOptions | 5300 | Embed is used when transformed structures should be embed into parent one. It&#39;s the same as gogoproto.embed flag, but right now I can&#39;t read gogoproto.embed option. DEPRECATED, use gogooproto.embed instead. |
| map_as | string | .google.protobuf.FieldOptions | 5304 | Contains name which will be used instead of current field name.

string street_1 = 1; -&gt; pb.go Street_1 instead Street1 |
| map_to | string | .google.protobuf.FieldOptions | 5303 | Points destination field type for OneOf fields. string one_of_to = 5302; Contains model&#39;s field name if it&#39;s different from name in messages. |
| skip | bool | .google.protobuf.FieldOptions | 5301 | If true, field will not be used in transform functions. |
| go_models_file_path | string | .google.protobuf.FileOptions | 5201 | Path to source file with Go structures which will be used as destination. |
| go_protobuf_package | string | .google.protobuf.FileOptions | 5203 | Package name with protobuf srtuctures. |
| go_repo_package | string | .google.protobuf.FileOptions | 5202 | Package name which contains model structures. |
| go_struct | string | .google.protobuf.MessageOptions | 5100 | Name of structure from repo package. |

 

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

