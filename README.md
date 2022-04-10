# Delegated login mono-repo

This repository contains three apps/artifacts:
- frontend-image
- login-app
- oidc-stub

`frontend-image` is used to package/host your frontend resources, and will utilize `login-app` if it detects that a request is from a non-authenticated user.
[See documentation](frontend-image/README.md)

`login-app` handles the communication between you systems and the IDP. [See documentation](login-app/README.md)

## Run it locally
All apps er configured in [docker-compose](docker-compose.yml). All you have to do is `docker-compose up`, and the test frontend should be available at `http://localhost:8083/frontend`. 
The tests are written in typescript/Deno. Run them with the following command: `deno test --allow-env --allow-net` from the `test` folder.