# Elasticsearch Theory
## Table of contents
## Overview
- Elasticsearch là một công cụ tìm kiếm và phân tích phân tán, có mã nguồn mở được xây dựng trên nền tảng Apache Lucene. Elasticsearch cung cấp khả năng tìm kiếm nhanh chóng, lưu trữ dữ liệu phi cấu trúc và phân tích dữ liệu theo thời gian thực.
## Các thành phần chính trong Elasticsearch
### 1. Cluster
- `Cluster` là một tập hợp các `node` (`server Elasticsearch`) làm việc cùng nhau để lưu trữ và tìm kiếm dữ liệu.
### 2. Node
- `Node` là một máy chủ Elasticsearch đơn lẻ trong một cụm `cluster`. Elasticsearch sẽ tự động phân bố dữ liệu và truy vấn đồng đều giữa các `node`.
- Các loại `node`:
    - `Master-eligible node`: Là các `node` có khả năng trở thành `master node`. `Master node` có nhiệm vụ tạo và xóa các `index`, quản lý và theo dõi các `node` có trong `cluster`, quản lí và phân phối các `shard` trong các `node`.
    - `Data node`: Là các `node` có trách nhiệm lưu trữ dữ liệu và thực hiện các thao tác liên quan đến dữ liệu.
    - `Ingest node`: Là các `node` có thể áp dụng một chuỗi các bước xử lí chuyển đổi, chuẩn hóa dữ liệu trước khi lưu dữ liệu vào Elasticsearch.
    - `Machine learning node`: Là các `node` có thể thực thi các tác vụ về `machine learning`.
    - `Coordination node`: Tất cả các `node` trong một `cluster` đều là `coordination node`. `Coordination node` có nhiệm vụ điều hướng các yêu cầu đến các `node` chứa dữ liệu cần thao tác và tổng hợp kết quả trả về cho `client`.
### 3. Index
- `Index` là một tập hợp các dữ liệu (`document`) có tính chất giống nhau. Nó giống như một bảng trong cơ sở dữ liệu quan hệ.
### 4. Document
- Dữ liệu sẽ được lưu trữ dưới dạng `document` và có cấu trúc dạng `JSON`.
### 5. Shard
- Mỗi `index` có thể được chia thành nhiều phần nhỏ, mỗi phần được gọi là một `shard`.
- Có 2 loại `shard` trong Elasticsearch
    - `Primary shard`:
        - Lưu trữ các `document` gốc. Mỗi document được lưu trữ trong một `primary shard` duy nhất.
    - `Replica shard`:
        - `Replica shard` là một bản sao hoàn chỉnh của `primary shard`.
        - Sử dụng `replica shard` giúp đảm bảo `index` có tính sẵn sàng cao, tránh trường hợp không thể xử lý truy vấn khi một hay nhiều `node` gặp sự cố.

## Các thao tác dữ liệu trong Elastichsearch
### 1. Các thao tác với index
- Tạo `index` mới
    ```json
    PUT /products
    {
        "settings": {
            "number_of_shards": 2,
            "number_of_replicas": 1
        }
    }
    ```
- Xóa `index` đã tồn tại
    ```json
    DELETE /products
    ```
### 2. Các thao tác với document
- Tạo `document` mới (`indexing document`)
    ```json
    POST /products/_doc
    {
        "name":"Coffee Maker",
        "price": 64,
        "in_stock": 100
    }
    ```
- Lấy thông tin về một `document` bằng id
    ```json
    GET /products/_doc/J3hlH5YB61kBimA1MfAz
    ```
- Cập nhật một `document`
    ```json
    POST /products/_update/J3hlH5YB61kBimA1MfAz
    {
        "doc": {
            "price": 100,
            "tags": ["electronics"]
        }
    }
    ```
- `Scripted update`
    - `Script` là một tính năng hỗ trợ chúng ta cập nhật `document` mà không cần biết trước giá trị của chúng
        ```json
        POST /<index>/_update/<document_id>
        {
            "script": {
                "source": "ctx._source.counter += params.increment",
                "params": {
                    "increment": 1
                }
            }
        }
        ```
- Thay thế một `document` bằng một `document` mới
    ```json
    PUT /products/_doc/document_id
    {
        "name":"Coffee Maker",
        "price": 64,
        "in_stock": 100
    }
    ```
- Xóa một `document`
    ```json
    DELETE /products/_doc/Fr62I5YBDS1ZHhN1dVEL
    ```

### 3. Routing trong Elasticsearch
- `Routing` trong Elasticsearch là cơ chế xác định xem `document` nên được lưu trữ và tìm kiếm trên shard nào.
- Khi một `document` được lưu trữ trong Elasticsearch, hệ thống cần quyết định `shard` nào sẽ giữ `document` đó. Quá trình này được gọi là `routing` và được thực hiện theo công thức:
    ```
    shard_num = hash(_routing) % num_primary_shards
    ```
- Khi đọc dữ liệu, sau khi được `routing` đến `replication group` chứa dữ liệu, Elasticsearch sẽ sử dụng `Adaptive Replica Selection` để chuyển truy vấn đến `shard` phù hợp nhất trong `replication group` để xử lý. `ARS` sẽ hoạt động giống như một `load balancer`.
- Khi viết dữ liệu, Elasticsearch sẽ chuyển yêu cầu đến cho `primary shard` để xử lý, sau đó, các `replica shard` sẽ đồng bộ dữ liệu với `primary shard`.

### 4. Bulk API
- `Bulk API` cho phép thực hiện nhiều thao tác với dữ liệu (`index`,`create`, `update`, `delete`) trong một yêu cầu duy nhất. API này giúp tăng hiệu suất đáng kể so với việc gửi nhiều yêu cầu riêng lẻ.
    ```json
    POST /_bulk
    { "index": { "_index": "san_pham", "_id": "1" } }
    { "ten": "Laptop Dell XPS", "gia": 30000000, "hang": "Dell", "ton_kho": 15 }
    { "create": { "_index": "san_pham", "_id": "2" } }
    { "ten": "iPhone 13", "gia": 22000000, "hang": "Apple", "ton_kho": 25 }
    { "update": { "_index": "san_pham", "_id": "3" } }
    { "doc": { "ton_kho": 10, "gia": 17500000 } }
    { "delete": { "_index": "san_pham", "_id": "4" } }
    ```
- Khi thực hiện yêu cầu, một thao tác thất bại sẽ không làm ảnh hưởng đến các thao tác khác.  Mỗi thao tác được xử lý độc lập và kết quả được trả về riêng biệt.

## Mapping và Analysis
### 1. Analysis
- `Analysis`, hay còn được gọi là `text analysis`, là quá trình phân tích các thuộc tính dạng văn bản khi lưu trữ (`index`) các `document`, biến các văn bản này thành các `token`. Các `token` này sẽ được lưu trữ sử dụng các cấu trúc dữ liệu tối ưu cho việc tìm kiếm.
- Một `analyzer` bao gồm ba thành phần chính:
    - `Character filters`: Được sử dụng để thêm, bớt hoặc thay đổi các ký tự trong văn bản.
    - `Tokenizers`: Sử dụng để chia văn bản thành các `token`.
    - `Token filters`: Nhận các `token` từ `tokenizers` và thay đổi, thêm, bớt các `token` này.
- Mặc định, Elasticsearch sẽ sử dụng `standard analyzer`.
- Elasticsearch cung cấp nhiều `character filters`, `tokenizers` và `token filters` được xây dựng sẵn. Ta có thể sử dụng chúng để xây dựng nên các `analyzer` mới.
- `Stemming` là quá trình cắt giảm một từ về dạng gốc (`stem`) của nó bằng cách loại bỏ các hậu tố và tiền tố. Mục đích của `stemming` là nhóm các biến thể của từ có cùng ngữ nghĩa để cải thiện kết quả tìm kiếm. Ví dụ:
    - "running", "runs", "ran" → "run"
    - "hotels", "hotel's" → "hotel"
    - "searching", "searched", "searches" → "search"
- `Stop words` là những từ xuất hiện rất phổ biến trong ngôn ngữ tự nhiên nhưng mang ít ý nghĩa ngữ nghĩa và thường không hữu ích cho việc tìm kiếm và đánh giá `relevance score`. Việc loại bỏ `stop words` giúp giảm kích thước `index`, cải thiện hiệu suất và tăng độ chính xác cho việc tìm kiếm. Ví dụ, stop words tiếng Anh bao gồm: "a", "an", "and", "are", "as", "at", "be", "but", "by",...

### 2. Mapping
- Là quá trình định nghĩa cấu trúc của một `document`. `Mapping` sẽ định nghĩa kiểu dữ liệu của từng thuộc tính trong `document`, các cấu hình và `analyzer` cho mỗi thuộc tính.
- Có hai loại `mapping` trong Elasticsearch
    - `Explicit Mapping`: Người dùng sẽ định nghĩa `mapping` cho `index`, chỉ rõ kiểu dữ liệu cho mỗi thuộc tính.
    - `Dynamic Mapping`: Với các thuộc tính mới (chưa được định nghĩa hoặc chưa tồn tại), Elasticsearch sẽ tự động phát hiện và xác định kiểu dữ liệu dựa trên các giá trị đầu vào.
- Các kiểu dữ liệu trong Elasticsearch
    - `long`, `integer`, `short`, `byte`: Các kiểu số nguyên khác nhau.
    - `double`, `float`, `half_float`: Các kiểu số thực.
    - `boolean`: Giá trị `true`/`false`.
    - `date`: Dùng cho ngày tháng, có thể được định dạng theo nhiều cách khác nhau.
    - `text`: Dùng cho giá trị dạng văn bản, được phân tích và lưu trữ tối ưu cho `full-text search`.
    - `keyword`: Dùng cho các giá trị chính xác, không được phân tích, thường dùng cho lọc, sắp xếp và tổng hợp.
    - `object`: Sử dụng để chứa một đối tượng dạng `JSON`. Sử dụng từ khóa `properties` để `mapping`.


### 3. Explicit Mapping
- Thêm `mapping` khi tạo `index`
    ```json
    PUT /reviews
    {
        "mappings": {
            "properties": {
                "rating" : {"type" : "float"},
                "content" : {"type": "text"},
                "product_id" : {"type" : "integer"},
                "author" : {
                    "type": "nested",
                    "properties": {
                        "name" : {"type" : "text"},
                        "email" : {"type" : "keyword"}
                    }
                }
            }
        }
    }
    ```
- Xem `mapping` của một `index`
    ```json
    GET /reviews/_mapping
    ```
- Thêm `mapping` (định nghĩa kiểu dữ liệu cho các thuộc tính mới) cho một `index` đã tồn tại
    ```json
    PUT /reviews/_mapping
    {
        "properties": {
            "created_at": {"type" : "date"}
        }
    }
    ```
### 4. Dynamic mapping
- Khi ta tạo một `document` chứa thuộc tính mới chưa biết đến, Elasticsearch sẽ tự động `mapping` kiểu dữ liệu dựa trên giá trị của thuộc tính đó.

- Để Elasticsearch có thể phát hiện dữ liệu dạng số khi đầu vào là văn bản, ta cần dùng
    ```json
    PUT /computers
    {
        "mappings": {
            "numeric_detection": true
        }
    }
    ```

- Theo mặc định, `date detection` luôn được bật và sẽ tự động `mapping` các giá trị văn bản có dạng ngày sang `date`

- Tắt `date detection`
    ```json
    PUT /computers
    {
        "mappings": {
            "date_detection": false
        }
    }
    ```
- Thay đổi format của `date detection`
    ```json
    PUT /computers
    {
        "mappings": {
            "dynamic_date_formats": ["dd-MM-yyyy"]
        }
    }
    ```
- Elasticsearch sẽ cung cấp 3 cấu hình cho `dynamic mapping`
    ```json
    PUT my_index
    {
        "mappings": {
            "dynamic": dynamic_mode
        }
    }
    ```
    - `True`: Elasticsearch sẽ tự động thêm `mapping` cho các thuộc tính mới.
    - `False`: Các thuộc tính mới vẫn sẽ được lưu trữ trong `_source` nhưng sẽ không được đánh chỉ mục (không thể tìm kiếm).
    - `Strict`: Elasticsearch sẽ báo lỗi nếu thuộc tính mới được thêm vào mà chưa định nghĩa trước.

## Searching
### 1. Introduction
- Elasticsearch cung cấp 2 cách tìm kiếm
    - `URI search`: Truy vấn tìm kiếm được viết trực tiếp trong `query parameter` và câu truy vấn được viết bằng ngôn ngữ truy vấn của `apache lucene`.
        ```json
        GET /ten_index/_search?q=truong:gia_tri
        ```
    - `Query DSL`: Truy vấn được viết trong `body` của `request`, cung cấp nhiều tính năng và cho phép tìm kiếm phức tạp hơn.
        ```json
        GET /products/_search
        {
            "query": {
                "match_all": {}
            }
        }
        ```
### 2. Term level queries
- Sử dụng để tìm kiếm giá trị chính xác của thuộc tính (`filtering`).
- `Term level queries` sẽ không được phân tích bởi `analyzer`, các giá trị dùng để tìm kiếm sẽ được giữ nguyên.
- `Term` sẽ tìm kiếm các tài liệu có chứa chính xác giá trị được chỉ định trong trường đã chọn.
    ```json
    GET /products/_search
    {
        "query": {
            "term": {
                "tags.keyword": {
                    "value": "Vegetable",
                    "case_insensitive": true
                }
            }
        }
    }
    ```
- `Terms` cho phép tìm kiếm sử dụng nhiều giá trị cùng một lúc và trả về `document` chứa ít nhất một giá trị trong dãy giá trị tìm kiếm.
    ```json
    GET /products/_search
    {
        "query": {
            "terms": {
                "tags.keyword": ["Soup", "Meat"]
            }
        }
    }
    ```
- Tìm kiếm theo `id`
    ```json
    GET /products/_search
    {
        "query": {
            "ids": {
                "values": ["100", "200", "300"]
            }
        }
    }
    ```
- `Range query` tìm kiếm các giá trị nằm trong một khoảng nhất định, thường dùng cho dữ liệu dạng số và ngày tháng.
    ```json
    GET /products/_search
    {
        "query": {
            "range": {
                "in_stock": {
                    "gte": 1,
                    "lte": 5
                }
            }
        }
    }
    ```

- `Prefix` tìm kiếm các `document` mà giá trị của thuộc tính bắt đầu bằng tiền tố được chỉ định.
    ```json
    GET /products/_search
    {
        "query": {
            "prefix": {
                "name.keyword": {
                    "value": "Past"
                }
            }
        }
    }
    ```
- `Wildcard` tìm kiếm các `document` dựa trên `pattern` sử dụng `wildcard`
    - `?`: Khớp với 1 ký tự bất kỳ.
    - `*`: Khớp với 0 hoặc nhiều ký tự bất kỳ.
    ```json
    GET /products/_search
    {
        "query": {
            "wildcard": {
                "tags.keyword": {
                    "value": "Bee*"
                }
            }
        }
    }
    ```
- `Exist` sẽ tìm kiếm các `document` có tồn tại `indexed value` cho một thuộc tính nào đó.
    ```json
    GET /products/_search
    {
        "query": {
            "exists": {
                "field": "tags.keyword"
            }
        }
    }
    ```
### 3. Full text queries
- Khác với `term level queries` tìm kiếm theo giá trị chính xác, `full text queries` cho phép chúng ta tìm kiếm các `document`có thuộc tính chứa cụm từ cần tìm kiếm hoặc đồng nghĩa hay liên quan về ngữ nghĩa với cụm từ cần tìm kiếm.
- `Full text search` sẽ phân tích văn bản đầu vào trước khi thực hiện tìm kiếm sử dụng các `analyzer`.
- `Match query` sẽ trả về các `document` chứa một hay nhiều `term`.
    ```json
    GET /products/_search
    {
        "query": {
            "match": {
                "name": {
                    "query": "pasta chicken",
                }
            }
        }
    }
    ```
- Sử dụng `fuzzy search` trong truy vấn `match` để tìm kiếm gần đúng.
    ```json
    GET /your_index/_search
    {
        "query": {
            "match": {
                "field_name": {
                    "query": "giá trị tìm kiếm",
                    "fuzziness": "AUTO"
                }
            }
        }
    }
    ```

- `Multi-match query` cho phép ta tìm kiếm trên nhiều thuộc tính cùng một lúc. `Document` sẽ được trả về nếu nó thỏa mãn một trong các thuộc tính.
    ```json
    GET /products/_search
    {
        "query": {
            "multi_match": {
                "query": "vegetable broth",
                "fields": ["name", "description"],
                "tie_breaker": 0.3
            }
        }
    }
    ```
- `Match query` sẽ trả về các `document` chứa các `term` mà không quan tâm đến vị trí và các `term` không cần liều với nhau. Để tìm kiếm các `document` chứa tất cả các `term` mà các `term` này đảm bảo vị trí so với văn bản đầu vào và liều nhau, ta dùng `match phrase query`.
    ```json
    GET /products/_search
    {
        "query": {
            "match_phrase": {
                "name": "mango juice"
            }
        }
    }
### 4. Compound queries
- `Leaf query` là một truy vấn độc lập tìm kiếm một giá trị trong một trường cụ thể.
- `Compound query` là truy vấn bao bọc nhiều `leaf queries` hoặc các `compound query` khác và kết hợp chúng theo logic.

### 5. Bool queries
- `Bool queries` là truy vấn cho phép chúng ta kết hợp nhiều điều kiện tìm kiếm với các toán tử logic.
- Các thành phần trong `bool query`
    - `Must`: Tất cả các điều kiện tìm kiếm trong `must` đều phải được thỏa mãn và sẽ tính toán `relevance score` của các thuộc tính tìm kiếm và đóng góp vào `relevance score` của `document` trả về.
    - `Filter`: Giống như `must`, tất cả các điều kiện tìm kiếm đều phải thỏa mãn nhưng sẽ không tính toán và đóng góp `relevance score`.
    - `Must not`: Tất cả các điều kiện tìm kiếm đều không được xuất hiện trong `document` và không ảnh hưởng đến `relevance score`. 
    - `Should`: Các điều kiện trong `should` không bắt buộc phải đáp ứng nhưng nếu đáp ứng sẽ làm tăng `relevance score`. 

- `Relevance score` sẽ là tổng `relevance score` của tất cả các điều kiện tìm kiếm.
- Ví dụ
    ```json
    GET /products/_search
    {
        "query": {
            "bool": {
                "filter": [
                    {
                        "range": {
                            "in_stock": {
                            "lte": 100
                            }
                        }
                    },
                    {
                        "term": {
                            "tags.keyword": "Beer"
                        }
                    }
                ],
                "should": [
                    { "match": { "name": "Beer" } },
                    { "match": { "description": "Beer" } }
                ]
            }
        }
    }
    ```

## Aggregation
### 1. Introduction
- `Aggregation` cho phép ta thực hiện các phân tích phức tạp và tổng hợp dữ liệu trên các `document` được lưu trữ trong Elasticsearch. 

### 2. Metric aggregations
- Metric Aggregations tính toán các số liệu thống kê (metrics) từ các giá trị của các trường trong `document`. Đây là những phép tính toán đơn giản nhưng rất hữu ích.

- Các loại Metric Aggregations phổ biến:

    - `Sum`: Tính tổng các giá trị
    - `Avg`: Tính giá trị trung bình
    - `Min/Max`: Tìm giá trị nhỏ nhất/lớn nhất
    - `Stats`: Trả về các thống kê cơ bản (min, max, sum, avg, count)
    - `Extended Stats`: Mở rộng từ stats, bao gồm thêm phương sai, độ lệch chuẩn
    - `Cardinality`: Đếm số lượng giá trị duy nhất
    - `Value Count`: Đếm số lượng giá trị không phải null

    ```json
    GET /orders/_search
    {
        "size": 0,
        "aggs": {
            "total_sales": {
                "sum": {
                    "field": "total_amount"
                }
            },
            "avg_sale": {
                "avg": {
                    "field": "total_amount"
                }
            },
            "min_sale": {
                "min": {
                    "field": "total_amount"
                }
            },
            "max_sale": {
                "max": {
                    "field": "total_amount"
                }
            }
        }
    }
    ```

### 3. Bucket aggregations
- `Bucket aggregations` sẽ phân nhóm các dữ liệu thành các `bucket` (nhóm) dựa trên các tiêu chí nhất định. Mỗi `bucket` sẽ chứa một tập hợp các tài liệu phù hợp với tiêu chí của `bucket` đó.
- `Terms Aggregation` tạo ra một bucket cho mỗi giá trị duy nhất của trường được chỉ định.
    ```json
    GET /orders/_search
    {
        "size": 0,
        "aggs": {
            "status_terms": {
                "terms": {
                    "field": "status",
                    "size": 20,
                    "missing": "N/A",
                    "min_doc_count": 0,
                    "order": {
                        "_key": "asc"
                    }
                }
            }
        }
    }
    ```
- Ta có thể sử dụng `nested aggregation` để áp dụng các `aggregation` lên các `bucket` được tạo ra bởi `bucket aggregation`.
    ```json
    GET /orders/_search
    {
        "size": 0,
        "query": {
            "range": {
                "total_amount": {
                    "gte": 100
                }
            }
        },
        "aggs": {
            "status_terms": {
                "terms": {
                    "field": "status"
                },
                "aggs": {
                    "status_stats": {
                        "stats": {
                            "field": "total_amount"
                        }
                    }
                }
            }
        }
    }
    ```
- `Filters Aggregation` tạo ra các bucket dựa trên một hoặc nhiều `filter`.
    ```json
    GET /recipes/_search
    {
        "size": 0,
        "aggs": {
            "my_filter": {
                "filters": {
                    "filters": {
                        "pasta": {
                            "match": {
                                "title": "pasta"
                            }
                        },
                        "spaghetti": {
                            "match": {
                                "title": "spaghetti"
                            }
                        }
                    }
                },
                "aggs": {
                    "avg_rating": {
                        "avg": {
                            "field": "ratings"
                        }
                    }
                }
            }
        }
    }
    ```
- Ta có thể chia `bucket` dựa trên `range` của thuộc tính.
    ```json
    GET /orders/_search
    {
        "size": 0,
        "aggs": {
            "amount_distribution": {
                "range": {
                    "field": "total_amount",
                    "ranges": [
                        {
                            "to": 50
                        },
                        {
                            "from": 50,
                            "to": 100 //excluded
                        },
                        {
                            "from": 100
                        }
                    ]
                }
            }
        }
    }
    ```

- `Histogram Aggregation` tạo ra các `bucket` với khoảng cách đều nhau dựa trên các giá trị số hoặc ngày tháng
    ```json
    GET /orders/_search
    {
        "size": 0,
        "aggs": {
            "amount_distribution": {
                "histogram": {
                    "field": "total_amount",
                    "interval": 25
                }
            }
        }
    }
    ```
- `Missing Aggregation` tạo một bucket cho các tài liệu không có giá trị cho trường được chỉ định.
    ```json
    GET /orders/_search
    {
        "size": 0,
        "aggs": {
            "orders_without_status": {
                "missing": {
                    "field": "status"
                }
            }
        }
    }
    ```
- `Nested Aggregation` cho phép thực hiện `aggregation` trên các trường `nested objects`.
    ```json
    GET /departments/_search
    {
        "size": 0,
        "aggs": {
            "employees": {
                "nested": {
                    "path": "employees"
                },
                "aggs": {
                    "minimun_age": {
                        "min": {
                            "field": "employees.age"
                        }
                    }
                }
            }
        }
    }
    ```
