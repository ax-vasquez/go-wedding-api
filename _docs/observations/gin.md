# `gin` Observations

A running document on things I observed while using `gin`.

## `c.Get()` vs `c.GetString()`

**Both of these target the same keys!** The only difference is that `c.Get()` will return the corresponding value as an `any` type, whereas
`c.GetString()` will obtain the corresponding value _as a string_ (if it's able).

If the value is NOT a string (for example, if you called `c.Set()` in a unit test and passed a `uuid.UUID` object as the `"uid"` value), then
using `c.GetString()` _will not return a value_, even though `c.Get()` would. Although you can update your code to use `c.Get()` instead of `c.GetString()`, you would need to add extra code to verify that the value is a valid UUID. In either case, steps must be taken to ensure the value is passed and consumed correctly.

The proper way to set the `"uid"` key in a unit test's request `context` is as follows:
```go
someMockUserId := uuid.New()
ctx := gin.CreateTestContextOnly(w, router)
// IMPORTANT - call .String() when passing a UUID to ensure c.GetString() will return a value later
ctx.Set("uid", mockInviterId.String())
// Use whatever arguments you actually need here
req, err := http.NewRequestWithContext(ctx, "GET", "/some/api/route", nil)
```