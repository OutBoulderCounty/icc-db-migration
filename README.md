# ICC DB Migration

## Purpose

To migrate ICC data from old MongoDB replica set to Planetscale database.

## Steps to run

1. Configure context for bastion instance

   ```sh
   # setting current context
   docker context use remote
   ```

1. Make sure you have .env file in place with full MONGODB_URI set
1. Run docker compose

   ```sh
   docker compose up --build
   ```
