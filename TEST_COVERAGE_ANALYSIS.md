# Test Coverage Analysis - GoGreTaGo

## Overall Statistics

| Layer | Source Files | Files With Tests | File Coverage | Test Count |
|---|---|---|---|---|
| Domain | 3 logic files | 3 | 100% | 47 |
| Application | 32 use cases | 32 | 100% | 98 |
| Infrastructure | 16 files | 12 | 75% | ~60 |
| Presentation | 26 files | 5 | **19%** | 25 |
| **Total** | **77** | **52** | **68%** | **~230** |

---

## Domain Layer

### Tested Files (3/3 logic files)

- `authorization/role_hierarchy_test.go` - All roles + edge cases (unknown, empty, case sensitivity). **Good.**
- `entities/pagination_test.go` - DefaultPagination, Skip, Take, BuildPaginationMeta. **Decent.**
- `errors/domain_errors_test.go` - All 16 constructors, HTTP status mapping, interface compliance. **Decent.**

### Bugs Found

1. **No-op assertion** in `domain_errors_test.go:213`:
   ```go
   assert.Equal(t, tt.code, tt.code) // compares value to itself -- always passes
   ```
   The `TestDomainErrors_EmbedDomainError` function (16 subtests) tests nothing. Should extract the embedded `DomainError.Code` field and compare against `tt.code`.

2. **Latent division-by-zero** in `pagination.go` - `BuildPaginationMeta` computes `float64(total) / float64(params.Limit)`. If `Limit == 0`, this produces `+Inf`/`NaN`. No guard and no test.

### Missing Tests

- `errors.As` not tested on concrete typed errors (`*UserAlreadyExistsError` etc.) -- only on plain `*DomainError`. Since concrete types embed by value, `errors.As` extraction may silently fail.
- `Page <= 0` producing negative `Skip()` offsets
- 6 error codes in `errorHTTPStatus` have no constructors: `UNAUTHORIZED`, `TOKEN_EXPIRED`, `TOKEN_INVALID`, `TOKEN_MALFORMED`, `VALIDATION_ERROR`, `RELATION_CONSTRAINT`
- Relative ordering assertion for role hierarchy (ADMIN > DRIVER > USER)
- `RoleHierarchy` map immutability guard

---

## Application Layer

### Tested Files (32/32 use cases)

Every use case has at least one test. Every business rule / domain error check is covered (30/30 branches).

### Ratings

| Rating | Use Cases |
|---|---|
| EXCELLENT | Login, Register, CreateInscription, GetUser |
| VERY GOOD | UpdateCar, CreateTrip, ListTrips, ListUsers |
| GOOD | CreateBrand, CreateCar, DeleteCar, UpdateColor, CreateDriver, DeleteInscription, DeleteTrip, FindTrips, AnonymizeUser, UpdateUser |
| ADEQUATE | DeleteBrand, DeleteCity, CreateColor, DeleteColor, GetTrip, ListTripPassengers, ListUserInscriptions |
| WEAK | ListBrands, ListCars, CreateCity, ListCities, ListColors, ListInscriptions |

### Systematic Gap: Repository Error Paths

Only **14 out of 62** possible `if err != nil` repository error branches are tested. The gap is consistent across all domains:

| Domain | Repo Error Branches Tested / Total |
|---|---|
| auth | 7/9 |
| brand | 1/4 |
| car | 0/16 |
| city | 0/4 |
| color | 0/8 |
| driver | 1/4 |
| inscription | 2/9 |
| trip | 1/10 |
| user | 3/6 |

### Weak List Use Cases

6 List use cases have only 1 test each (happy path only):
- `ListBrands`, `ListCars`, `ListCities`, `ListColors`, `ListInscriptions`, `ListTripPassengers`

Compare with `ListTrips` and `ListUsers` which correctly test: success + empty results + repository error.

---

## Infrastructure Layer

### Tested Files (12/16)

**Well tested:**
- `cache_test.go` + `cache_integration_test.go` - Disabled cache paths, Redis roundtrip, invalidation pattern
- All 9 repository integration tests - Using testcontainers with PostgreSQL 16
- `argon_password_service_test.go` - 7 tests including corrupted hash/salt error paths (exemplary)
- `jwt_service_impl_test.go` - 7 tests with table-driven expiration format tests

### Untested Files

| File | Impact |
|---|---|
| `gorm_model_repository.go` | **No test file at all** - FindAll, FindByID, FindByNameAndBrand untested |
| `resend_email_service.go` | No tests |
| `database/gorm_logger.go` | No tests |
| `database/postgres.go` | No dedicated tests (models tested indirectly) |

### Source Code Bugs Found

1. **Ignored `Count` errors** in all 6 repository `FindAll` methods:
   ```go
   r.db.WithContext(ctx).Model(&BrandModel{}).Count(&total) // error ignored
   ```
   If the count query fails, `total` silently stays 0. Affects: brand, car, city, color, trip, inscription repositories.

2. **Silent deletion failures** in `cache.go:102`:
   ```go
   c.client.Del(ctx, iter.Val()) // error from Del ignored
   ```

### Critical Missing Tests

1. **`FindByFilters` (trip repository)** - The most complex query in the codebase (dynamic SQL with subqueries joining city_trips and cities). Zero test coverage.
2. **No unique constraint violation tests** across any repository (duplicate emails, license plates, driver licenses)
3. **Car `FindAll`**, City `FindByID`, Inscription `FindAll` - Methods that exist in source but are never tested
4. **Cache TTL expiration** - TTL is configured but never verified
5. **Inscription `CountByTripRefID` ACTIVE filter** - Only counts ACTIVE inscriptions but no test creates a non-ACTIVE inscription to verify the filter works

### Minor Issues

- JWT tests use `os.Setenv`/`os.Unsetenv` which is not safe for `t.Parallel()`
- No test for `UpdateRole` with non-existent RefID (GORM returns nil error for 0 rows affected)

---

## Presentation Layer

### Tested Files (5/26) -- THE BIGGEST GAP

**What IS tested (well):**
- `controllers/pagination_test.go` - Defaults, custom values, limit capping, invalid/non-numeric input
- `middleware/auth_middleware_test.go` - Bearer token, x-auth-token, missing/invalid/nil token
- `middleware/authorization_middleware_test.go` - Role hierarchy, forbidden, missing/invalid role
- `middleware/error_handler_test.go` - All 4 error types + passthrough (best test file)
- `validators/auth_validator_test.go` - Custom password/hexcolor validators, format errors

### Untested Controllers (9 controllers, 30 handler methods)

| Controller | Untested Methods |
|---|---|
| `auth_controller.go` | Register, Login |
| `brand_controller.go` | ListBrands, CreateBrand, DeleteBrand |
| `car_controller.go` | ListCars, CreateCar, UpdateCar, PatchCar, DeleteCar |
| `city_controller.go` | ListCities, CreateCity, DeleteCity |
| `color_controller.go` | ListColors, CreateColor, UpdateColor, DeleteColor |
| `driver_controller.go` | CreateDriver |
| `inscription_controller.go` | ListInscriptions, CreateInscription, DeleteInscription, ListUserInscriptions, ListTripPassengers |
| `user_controller.go` | ListUsers, GetUser, UpdateProfile, AnonymizeMe, AnonymizeUser |
| `trip_controller.go` | ListTrips, GetTrip, FindTrip, CreateTrip, DeleteTrip |

Each handler has at least 3 untested code paths: success, JSON binding error, and use case error.

`UserController.GetUser` is particularly concerning -- it contains inline authorization logic (`if id != requestingUserID && role != "ADMIN"`) that bypasses middleware-based role checking.

### Untested Middleware (4 files)

| Middleware | Behavior | Risk |
|---|---|---|
| `rate_limiter.go` | Per-IP rate limiting, 429 responses | Security |
| `secure_headers.go` | 7 security headers (HSTS, CSP, X-Frame-Options, etc.) | Security |
| `body_limit.go` | MaxBytesReader request size enforcement | Security |
| `request_logger.go` | X-Request-Id header, request logging | Observability |

### Untested Routes (all 10 files)

No route file has tests. The middleware chains (Auth -> RequireRole -> Controller) are never tested end-to-end. For example, no test verifies that `POST /api/v1/brands` requires both authentication AND ADMIN role.

### Missing Scenarios in Existing Tests

- `auth_middleware`: Malformed Authorization header (e.g. "Basic abc"), empty Bearer token, both headers set
- `authorization_middleware`: USER accessing USER route (basic happy path), DRIVER accessing ADMIN route (should fail), variadic role matching
- `error_handler`: Multiple errors in `c.Errors`, different domain error codes mapping to different HTTP statuses
- `auth_validator`: `gt` and `eqfield` tag message paths, `GetValidator()` singleton behavior
- `pagination`: Boundary values `page=0`, `limit=0`, `limit=100`, `limit=101`

---

## Top 10 Recommendations (Prioritized)

### Critical (bugs / security)

1. **Fix no-op assertion** in `domain_errors_test.go:213`. Replace `assert.Equal(t, tt.code, tt.code)` with embedded field extraction.

2. **Add guard + test for `Limit=0`** in `BuildPaginationMeta` to prevent division by zero.

3. **Test `FindByFilters`** in trip repository -- most complex query, zero coverage.

4. **Add controller tests** starting with `AuthController` (Register/Login) and `UserController.GetUser` (inline authz).

### High Priority

5. **Test `secure_headers.go` and `rate_limiter.go`** -- security middleware that regressions would silently break.

6. **Add repository error tests** to use cases (14/62 covered). Start with CreateCar (6 untested), UpdateCar (6), CreateTrip (5).

7. **Fix ignored `Count` errors** in all 6 repository `FindAll` methods.

### Medium Priority

8. **Add `GormModelRepository` integration tests** -- only repository without a test file.

9. **Standardize List use case tests** to match `ListTrips`/`ListUsers` pattern (success + empty + error).

10. **Test `errors.As` on concrete typed errors** to verify unwrapping from `*UserAlreadyExistsError` to `*DomainError` works.
