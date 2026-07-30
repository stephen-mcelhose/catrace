---
title: API & Protocol Buffer Code Review Best Practices
description: API and Protobuf design review following Google AIP guidelines, versioning, and naming conventions.
---

# API & Protocol Buffer Code Review Best Practices

## Google AIP (API Improvement Proposals) Compliance

All APIs and Protocol Buffers must follow [Google AIP guidelines](https://google.aip.dev/). Key AIPs to verify during review:

### Core AIPs

| AIP | Topic | Quick Check |
|-----|-------|-------------|
| [AIP-121](https://google.aip.dev/121) | Resource-oriented design | Resources have clear hierarchy |
| [AIP-122](https://google.aip.dev/122) | Resource names | Format: `projects/{project}/resources/{resource}` |
| [AIP-123](https://google.aip.dev/123) | Resource types | Fully qualified: `example.googleapis.com/Resource` |
| [AIP-127](https://google.aip.dev/127) | HTTP and gRPC transcoding | HTTP annotations present |
| [AIP-131](https://google.aip.dev/131) | Standard methods: Get | `GetResource(GetResourceRequest)` |
| [AIP-132](https://google.aip.dev/132) | Standard methods: List | `ListResources(ListResourcesRequest)` |
| [AIP-133](https://google.aip.dev/133) | Standard methods: Create | `CreateResource(CreateResourceRequest)` |
| [AIP-134](https://google.aip.dev/134) | Standard methods: Update | `UpdateResource(UpdateResourceRequest)` |
| [AIP-135](https://google.aip.dev/135) | Standard methods: Delete | `DeleteResource(DeleteResourceRequest)` |
| [AIP-136](https://google.aip.dev/136) | Custom methods | `:customVerb` suffix |
| [AIP-140](https://google.aip.dev/140) | Field names | snake_case, no abbreviations |
| [AIP-141](https://google.aip.dev/141) | Quantity fields | Use specific types (count, size) |
| [AIP-142](https://google.aip.dev/142) | Time fields | Use `google.protobuf.Timestamp` |
| [AIP-143](https://google.aip.dev/143) | Standardized codes | Use standard error codes |
| [AIP-151](https://google.aip.dev/151) | LRO (Long Running Operations) | For async operations |
| [AIP-154](https://google.aip.dev/154) | ETags | For optimistic concurrency |
| [AIP-158](https://google.aip.dev/158) | Pagination | `page_size`, `page_token`, `next_page_token` |
| [AIP-180](https://google.aip.dev/180) | Naming conventions | Consistent terminology |

## Protocol Buffer Best Practices

### Resource Definition (AIP-121, AIP-123)

```protobuf
// ❌ BAD: No resource annotation, unclear hierarchy
message Book {
  string id = 1;
  string title = 2;
}

// ✅ GOOD: Proper resource annotation
message Book {
  option (google.api.resource) = {
    type: "library.example.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };

  // Resource name. Format: publishers/{publisher}/books/{book}
  string name = 1 [(google.api.field_behavior) = IDENTIFIER];

  string title = 2 [(google.api.field_behavior) = REQUIRED];
  string author = 3;
  google.protobuf.Timestamp create_time = 4 [(google.api.field_behavior) = OUTPUT_ONLY];
  google.protobuf.Timestamp update_time = 5 [(google.api.field_behavior) = OUTPUT_ONLY];
}
```

### Standard Methods (AIP-131 to AIP-135)

```protobuf
service LibraryService {
  // Get a single book.
  rpc GetBook(GetBookRequest) returns (Book) {
    option (google.api.http) = {
      get: "/v1/{name=publishers/*/books/*}"
    };
    option (google.api.method_signature) = "name";
  }

  // List books.
  rpc ListBooks(ListBooksRequest) returns (ListBooksResponse) {
    option (google.api.http) = {
      get: "/v1/{parent=publishers/*}/books"
    };
    option (google.api.method_signature) = "parent";
  }

  // Create a book.
  rpc CreateBook(CreateBookRequest) returns (Book) {
    option (google.api.http) = {
      post: "/v1/{parent=publishers/*}/books"
      body: "book"
    };
    option (google.api.method_signature) = "parent,book";
  }

  // Update a book.
  rpc UpdateBook(UpdateBookRequest) returns (Book) {
    option (google.api.http) = {
      patch: "/v1/{book.name=publishers/*/books/*}"
      body: "book"
    };
    option (google.api.method_signature) = "book,update_mask";
  }

  // Delete a book.
  rpc DeleteBook(DeleteBookRequest) returns (google.protobuf.Empty) {
    option (google.api.http) = {
      delete: "/v1/{name=publishers/*/books/*}"
    };
    option (google.api.method_signature) = "name";
  }
}
```

### Request/Response Messages (AIP-131 to AIP-135)

```protobuf
// ✅ GOOD: Proper Get request
message GetBookRequest {
  // Required. The name of the book to retrieve.
  // Format: publishers/{publisher}/books/{book}
  string name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference) = {
      type: "library.example.com/Book"
    }
  ];
}

// ✅ GOOD: Proper List request with pagination (AIP-158)
message ListBooksRequest {
  // Required. The parent publisher.
  // Format: publishers/{publisher}
  string parent = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference) = {
      child_type: "library.example.com/Book"
    }
  ];

  // The maximum number of books to return.
  int32 page_size = 2;

  // A page token from a previous ListBooks call.
  string page_token = 3;
}

message ListBooksResponse {
  // The books.
  repeated Book books = 1;

  // Token for the next page.
  string next_page_token = 2;
}

// ✅ GOOD: Proper Create request
message CreateBookRequest {
  // Required. The parent publisher.
  string parent = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference) = {
      child_type: "library.example.com/Book"
    }
  ];

  // Required. The book to create.
  Book book = 2 [(google.api.field_behavior) = REQUIRED];

  // Optional. The ID to use for the book.
  string book_id = 3;
}

// ✅ GOOD: Proper Update request with field mask (AIP-134)
message UpdateBookRequest {
  // Required. The book to update.
  Book book = 1 [(google.api.field_behavior) = REQUIRED];

  // The fields to update.
  google.protobuf.FieldMask update_mask = 2;
}

// ✅ GOOD: Proper Delete request
message DeleteBookRequest {
  // Required. The name of the book to delete.
  string name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference) = {
      type: "library.example.com/Book"
    }
  ];
}
```

### Field Naming (AIP-140)

```protobuf
// ❌ BAD: Incorrect field naming
message BadExample {
  string ID = 1;           // Should be lowercase
  string userName = 2;     // Should be snake_case
  string num_items = 3;    // Abbreviation, should be item_count
  int32 ts = 4;            // Cryptic, should be timestamp or create_time
}

// ✅ GOOD: Correct field naming
message GoodExample {
  string id = 1;
  string user_name = 2;
  int32 item_count = 3;
  google.protobuf.Timestamp create_time = 4;
}
```

### Time Fields (AIP-142)

```protobuf
// ❌ BAD: String or int64 for timestamps
message BadExample {
  string created_at = 1;     // Should use Timestamp
  int64 updated_at_unix = 2; // Should use Timestamp
}

// ✅ GOOD: Use google.protobuf.Timestamp
import "google/protobuf/timestamp.proto";
import "google/protobuf/duration.proto";

message GoodExample {
  google.protobuf.Timestamp create_time = 1;
  google.protobuf.Timestamp update_time = 2;
  google.protobuf.Duration timeout = 3;  // For durations
}
```

### Custom Methods (AIP-136)

```protobuf
// ✅ GOOD: Custom method with :verb suffix
rpc ArchiveBook(ArchiveBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:archive"
    body: "*"
  };
}

rpc MoveBook(MoveBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:move"
    body: "*"
  };
}
```

### ETags for Concurrency (AIP-154)

```protobuf
message Book {
  string name = 1;
  string title = 2;

  // ETag for optimistic concurrency control.
  string etag = 99;
}

message UpdateBookRequest {
  Book book = 1;
  google.protobuf.FieldMask update_mask = 2;
}

message DeleteBookRequest {
  string name = 1;

  // Optional. If set, delete only if etag matches.
  string etag = 2;
}
```

### Long-Running Operations (AIP-151)

```protobuf
import "google/longrunning/operations.proto";

rpc ImportBooks(ImportBooksRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:import"
    body: "*"
  };
  option (google.longrunning.operation_info) = {
    response_type: "ImportBooksResponse"
    metadata_type: "ImportBooksMetadata"
  };
}
```

## API Design Checklist

### Resource Design

- [ ] Resources have clear parent-child relationships
- [ ] Resource names follow `{collection}/{resource}` pattern
- [ ] Resource types are fully qualified
- [ ] Resources have `name` field as identifier

### Standard Methods

- [ ] CRUD operations use standard method names
- [ ] Request messages follow naming convention (`{Method}{Resource}Request`)
- [ ] HTTP annotations are correct
- [ ] Method signatures defined for common parameters

### Field Design

- [ ] Field names are snake_case
- [ ] No abbreviations in field names
- [ ] Timestamps use `google.protobuf.Timestamp`
- [ ] Durations use `google.protobuf.Duration`
- [ ] Field behaviors annotated (REQUIRED, OUTPUT_ONLY, etc.)

### Pagination

- [ ] List methods support pagination
- [ ] Uses `page_size` and `page_token` request fields
- [ ] Uses `next_page_token` response field
- [ ] Default page size is reasonable

### Error Handling

- [ ] Uses standard gRPC/HTTP error codes
- [ ] Error messages are actionable
- [ ] Includes relevant error details

## Linting Tools

```bash
# buf lint (recommended)
buf lint

# api-linter (Google's AIP linter)
api-linter --set-exit-status protos/

# protolint
protolint lint protos/
```

## ReferenceSearch Queries

Use these to find AIP-compliant patterns in reference repos:

```
ReferenceSearch(query="resource annotation proto", repo="csgda-kit")
ReferenceSearch(query="ListRequest pagination", repo="csgda-kit")
ReferenceSearch(query="google.api.http annotation", repo="csgda-kit")
ReferenceSearch(query="field_behavior REQUIRED", repo="csgda-kit")
```

## Security Checklist for APIs

- [ ] Authentication required on all endpoints
- [ ] Authorization checks for resource access
- [ ] Input validation (field limits, formats)
- [ ] Rate limiting configured
- [ ] Audit logging for sensitive operations
- [ ] No sensitive data in URLs (use request body)
- [ ] TLS required
- [ ] Proper CORS configuration (if applicable)
