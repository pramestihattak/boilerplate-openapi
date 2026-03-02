---
name: generate-tests
description: Generate unit tests (service layer) and integration tests (storage layer) for this Go project. Use when adding or completing tests for service methods or storage functions.
disable-model-invocation: true
---

# Test Generation Standards

## Workflow

1. **Check for existing tests** — look for `*_test.go` files alongside the target source file.
2. **If tests exist** — run them (`go test ./...` scoped to the package). If they fail, diagnose and fix before proceeding.
3. **If no tests exist** — create them following the standards below.
4. **Run after writing** — always run the tests after creating or modifying them. Fix any failures before finishing.

When generating unit tests or integration tests for this project, follow these standards exactly.

---

## Unit Tests (service layer)

**File:** `<package>/<function>_test.go`
**Package:** `package <domain>_test` (external test package)
**Run:** `go test ./service/...`

### Struct shape

```go
tests := []struct {
    name      string
    input     <InputType>
    mockSetup func(*mock<Package>.Mock<Interface>)
    wantOut   <OutputType>
    wantErr   error
}{...}
```

- `wantOut` and `wantErr` — always use `want*` prefix, never `expected*`
- Zero values are implicit — omit `wantOut` or `wantErr` fields when zero/nil
- `mockSetup` always takes **only the mock** as its argument; access `tc.input` via closure (Go 1.22+ loop variables are per-iteration)

### Runner

```go
for _, tc := range tests {
    t.Run(tc.name, func(t *testing.T) {
        ctrl := gomock.NewController(t)
        defer ctrl.Finish()

        store := mock<Package>.New<MockType>(ctrl)
        tc.mockSetup(store)

        svc := <pkg>.New(logrus.New(), store)
        got, err := svc.<Method>(context.Background(), tc.input)

        assert.Equal(t, tc.wantErr, err)
        assert.Equal(t, tc.wantOut, got)
    })
}
```

- Result variable is always `got`
- Always use `context.Background()` in unit tests

### Mock expectations

| Argument type | How to match |
|---|---|
| `context.Context` | always `gomock.Any()` |
| Static / known value | exact value literal |
| Dynamic value (hashed password, random token) | `gomock.Any()` + `.Do()` to assert the static fields within the call |

**Dynamic arg example:**
```go
store.EXPECT().
    Register(gomock.Any(), gomock.Any()).
    Do(func(_ context.Context, reg storageAuth.Register) {
        assert.Equal(t, "Test User", reg.FullName)   // static — assert exactly
        assert.Equal(t, "user@example.com", reg.Email)
        assert.NotEmpty(t, reg.Password)              // dynamic — just assert non-empty
        assert.NotEqual(t, "password123", reg.Password, "must be hashed")
        assert.Len(t, reg.VerificationToken, 10)      // dynamic — assert length
    }).
    Return(validUserID, nil)
```

Never use bare `gomock.Any()` without `.Do()` for arguments that carry business logic.

---

## Integration Tests (storage layer)

**File:** `storage/<domain>/<function>_integration_test.go`
**Package:** same package as source (e.g. `package auth_storage`) — required for unexported `db` access
**Build tag:** `//go:build integration` on first line
**Run:** `go test -tags integration ./storage/...`

### Setup file (`setup_integration_test.go`)

Every storage package gets one shared setup file containing:

```go
//go:build integration

package <storage_package>

var testStorage *Storage

func TestMain(m *testing.M) {
    dbConfig := util.DBConfig{
        User:     getenv("TEST_DB_USER", "postgres"),
        Host:     getenv("TEST_DB_HOST", "localhost"),
        Port:     getenv("TEST_DB_PORT", "5438"),
        DBName:   getenv("TEST_DB_NAME", "postgres"),
        Password: getenv("TEST_DB_PASSWORD", "secret"),
        SSLMode:  "disable",
    }
    connStr, _ := util.NewDBStringFromDBConfig(dbConfig)
    db, err := NewDbConn(logrus.New(), connStr)
    if err != nil {
        log.Printf("skipping integration tests: cannot connect to postgres: %v", err)
        os.Exit(0) // graceful skip — safe for CI without a DB
    }
    testStorage = &Storage{logger: logrus.New(), db: db}
    os.Exit(m.Run())
}

func getenv(key, fallback string) string { ... }
func seedUser(t *testing.T, reg Register, verified bool) string { ... }
func cleanupByEmail(t *testing.T, emails ...string) { ... }
```

### Struct shape

```go
tests := []struct {
    name    string
    seed    *<SeedType>              // optional pre-existing DB row
    input   <InputType>
    wantErr bool
    verify  func(t *testing.T, got <OutputType>) // for non-trivial output assertions
}{...}
```

### Runner

```go
for _, tc := range tests {
    t.Run(tc.name, func(t *testing.T) {
        // 1. defer cleanup FIRST — before anything that can fail
        if tc.seed != nil {
            seedUser(t, *tc.seed, false)
            defer cleanupByEmail(t, tc.seed.Email)
        }

        got, err := testStorage.<Method>(context.Background(), tc.input)

        if tc.wantErr {
            assert.Error(t, err)
            return
        }
        require.NoError(t, err)

        if tc.verify != nil {
            tc.verify(t, got)
        }
    })
}
```

- Result variable is always `got`
- verify function parameter is always `got`
- `defer cleanupByEmail` must be called **before any assertion** so cleanup runs even on failure
- Use `require.NoError` for errors that make further assertions meaningless; use `assert` for value checks
- When multiple read-only subtests share one seed row, seed once in the outer function scope and defer cleanup there

### Verify function — DB state assertions

For write operations (INSERT, UPDATE), the verify function should query the DB directly to confirm the persisted state:

```go
verify: func(t *testing.T, got string) {
    var verified bool
    var token string
    err := testStorage.db.QueryRowContext(
        context.Background(),
        "SELECT verified, verification_token FROM users WHERE email = $1",
        "user@example.com",
    ).Scan(&verified, &token)
    require.NoError(t, err)
    assert.True(t, verified)
    assert.Empty(t, token)
},
```

### What to test at the storage layer

Only test what the SQL layer enforces. Duplicate-prevention or business-rule checks that live in the service layer (e.g. `UserExist` before `Register`) are **not** storage-layer concerns and should not have integration test cases here.

---

## Assertions quick-reference

| Situation | Use |
|---|---|
| Setup / structural failure | `require.NoError(t, err)` — stops the test immediately |
| Value comparison | `assert.Equal(t, expected, got)` |
| Error expected | `assert.Error(t, err)` |
| Nil check | `assert.Nil(t, got)` / `assert.NotNil(t, got)` |
| Non-empty | `assert.NotEmpty(t, got)` |
| Length | `assert.Len(t, got, n)` |
| UUID validity | `_, err := uuid.Parse(got); assert.NoError(t, err)` |
