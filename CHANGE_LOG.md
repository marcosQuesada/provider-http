### Main Features included
---

- [X] Decouple Mapping Actions
- [X] Make Mapping Actions Optional
- [X] RawExtension replaces request Body
- [X] RawExtension replaces request Mapping Body
- [X] RawExtension replaces response Body
- [X] Request Resource Referencing
  - [X] Add Resource Reference Conditions (Blocking behaviour)
- [X] Upgrade Crossplane-Runtime version, support ManagementPolicies
- [X] Imported secretInjectionConfigs feature
  - PatchFrom a Secret Supports transformations (base64) 
    - [X] Workaround
    - [ ] Stable solution
- Request
  - [ ] SilentBody
  - [ ] SilentHeaders
  - [ ] SilentResponseHeaders
  - [ ] SilentResponseBody
---
### FIXES
- Request with StatusCode different than 200
  - they don't let you proceed in Deletion
  - or, takes a lot of time to confirm deletion
