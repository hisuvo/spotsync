# Implementation Plan - Parking Spot Reservation

This plan outlines the design and implementation details for the **Parking Spot Reservation** module, addressing endpoints, access control, and the critical concurrency capacity checks.

## Proposed Changes

### Database Layer

#### [MODIFY] [migrate.go](file:///d:/MISSION/Asseinment/spotsync/internal/database/migrate.go)
- Include `&reservations.Reservation{}` in GORM auto-migration to register the `reservations` table.

### Domain Layer (Reservations)

#### [MODIFY] [entity.go](file:///d:/MISSION/Asseinment/spotsync/internal/domain/reservations/entity.go)
- Update `StatusCancelled` constant value to `"cancelled"` (instead of `"canclled"`).
- Make sure `ToResponse` maps `UpdatedAt` field.

#### [MODIFY] [dto/response.go](file:///d:/MISSION/Asseinment/spotsync/internal/domain/reservations/dto/response.go)
- Add `UpdatedAt` field to `ReservationResponse`.
- Add `MyReservationResponse` and `MyReservationZoneDTO` to support the nested structure expected by the `/my-reservations` endpoint.
- Add `AdminReservationResponse`, `AdminUserDTO`, and `AdminZoneDTO` to support the admin reservations listing.

#### [MODIFY] [repository.go](file:///d:/MISSION/Asseinment/spotsync/internal/domain/reservations/repository.go)
- Modify `GetMy` to preload `"Zone"`.
- Modify `GetAll` to preload `"User"` and `"Zone"`.
- Implement status update to `"cancelled"` inside `Cancel` (instead of deleting from DB) and check authorization/existence properly.

#### [MODIFY] [service.go](file:///d:/MISSION/Asseinment/spotsync/internal/domain/reservations/service.go)
- Define the `Service` interface and its implementation `service`.
- Add mapping logic from entities to the respective response DTOs.

#### [MODIFY] [handler.go](file:///d:/MISSION/Asseinment/spotsync/internal/domain/reservations/handler.go)
- Implement `handler` struct and HTTP handlers using Echo framework:
  - `Create`: Extract user ID from JWT context claims. Bind and validate request body. Handle `ErrZoneFull` (return 400 Bad Request) and database errors.
  - `GetMy`: Get reservations of the logged-in user.
  - `Cancel`: Check permissions. If driver tries to cancel someone else's reservation, return 403 Forbidden. Otherwise update status to `"cancelled"`.
  - `GetAll`: Admin-only endpoint. List all reservations in the system.

#### [MODIFY] [register.go](file:///d:/MISSION/Asseinment/spotsync/internal/domain/reservations/register.go)
- Define `RegisterRoute` that creates repository, service, handler, and registers HTTP routes under `/api/v1/reservations`:
  - `POST /api/v1/reservations` -> Driver/Admin authenticated.
  - `GET /api/v1/reservations/my-reservations` -> Driver/Admin authenticated.
  - `DELETE /api/v1/reservations/:id` -> Driver/Admin authenticated.
  - `GET /api/v1/reservations` -> Admin authenticated only.

### Routing & Server Layer

#### [MODIFY] [http.go](file:///d:/MISSION/Asseinment/spotsync/internal/server/http.go)
- Register `reservations.RegisterRoute(e, db, cfg)` to start serving the reservation endpoints.

---

## Verification Plan

### Automated Verification
- Verify the build and compilation status of the project using `go build ./cmd/main.go`.

### Manual / API Verification
- Perform manual requests using a custom integration script (e.g. testing using curl/powershell) or launching the dev server and sending HTTP requests to verify correct behavior:
  - Authentication checks and role enforcement (Driver vs. Admin).
  - Validation of concurrency rule (over-capacity rejection).
  - Preloaded structures mapping correctness.
