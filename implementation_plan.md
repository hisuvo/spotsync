# Implementation Plan - Admin Access Control and DTO Mapping

This implementation plan details the security enhancements to restrict parking zone creation and general reservations access to admin users only, as well as finishing the DTO mapping for the reservation retrieval endpoints.

## User Review Required

> [!IMPORTANT]
> The role check uses `claims.Email` to determine the user's role because the JWT service in `user/service.go` maps the `user.Role` to the `Email` claim field during token generation.

## Proposed Changes

### Domain Layer (Parking Zones)

---

#### [MODIFY] [register.go](file:///d:/MISSION/Asseinment/spotsync/internal/domain/parkingzones/register.go)
- Initialize `jwtService` and `authMiddleware`.
- Route `POST /zones`, `PUT /zones/:id`, and `DELETE /zones/:id` through the `authMiddleware` to enforce authentication.

#### [MODIFY] [handler.go](file:///d:/MISSION/Asseinment/spotsync/internal/domain/parkingzones/handler.go)
- Check `middleware.CheckUser(c).Email == "admin"` in the `Create`, `Update`, and `Delete` endpoints.
- Return `403 Forbidden` with a standardized error response if the user is not an admin.

### Domain Layer (Reservations)

---

#### [MODIFY] [service.go](file:///d:/MISSION/Asseinment/spotsync/internal/domain/reservations/service.go)
- Update `Service` interface to return `[]dto.MyReservationResponse` for `GetMyReservations` and `[]dto.AdminReservationResponse` for `GetAllReservations`.
- Implement mapping logic to transform the database models (`[]Reservation`) into their corresponding response DTOs.

#### [MODIFY] [handler.go](file:///d:/MISSION/Asseinment/spotsync/internal/domain/reservations/handler.go)
- In `GetAll` endpoint, check `middleware.CheckUser(c).Email == "admin"`.
- Return `403 Forbidden` if the user is not an admin.

## Verification Plan

### Automated Tests
- Build and compile the project using `go build -o spotsync.exe ./cmd/main.go`.

### Manual Verification
- We can verify the endpoints by starting the server and making test requests:
  - **Zone Creation**: Send `POST /api/v1/zones` with a driver's token and verify it returns `403 Forbidden`. Send with an admin's token and verify it works.
  - **All Reservations**: Send `GET /api/v1/reservations` with a driver's token and verify it returns `403 Forbidden`. Send with an admin's token and verify it returns the list of reservations mapped to `AdminReservationResponse` DTOs.
