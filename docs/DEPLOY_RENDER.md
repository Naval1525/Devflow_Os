# Deploy Backend (Go API) to Render

Step-by-step guide to deploy only the **backend API** to [Render](https://render.com).

---

## Prerequisites

- GitHub (or GitLab) repo with this project
- [Render](https://render.com) account (free tier works)
- PostgreSQL database (e.g. [Neon](https://neon.tech) or Render PostgreSQL)
- Values for `DATABASE_URL`, `JWT_SECRET`, and optionally `GEMINI_API_KEY`

---

## Step 1: Push your code to GitHub

Make sure your project is in a Git repo and pushed to GitHub (or GitLab). Render will build and run from this repo.

```bash
git add .
git commit -m "Prepare for Render deploy"
git push origin main
```

---

## Step 2: Create a new Web Service on Render

1. Go to [dashboard.render.com](https://dashboard.render.com) and log in.
2. Click **New +** → **Web Service**.
3. Connect your GitHub (or GitLab) account if you haven’t already.
4. Select the **repository** that contains this project (e.g. `devflowos` or your repo name).
5. Click **Connect**.

---

## Step 3: Configure the Web Service

Use these settings (Render will prompt for them):

| Field | Value |
|-------|--------|
| **Name** | `devflowos-api` (or any name you like) |
| **Region** | Choose closest to you (e.g. Oregon, Frankfurt) |
| **Branch** | `main` (or your default branch) |
| **Root Directory** | `services/api` |
| **Runtime** | **Go** |
| **Build Command** | `go build -o api ./cmd/api` |
| **Start Command** | `./api` |

- **Root Directory** must be `services/api` so Render runs all commands from the API folder (where `go.mod` lives).
- Render sets `PORT` automatically; your app already reads `PORT` from the environment.

---

## Step 4: Add environment variables

In the same Web Service, open **Environment** and add:

| Key | Value | Notes |
|-----|--------|--------|
| `DATABASE_URL` | `postgres://user:password@host/db?sslmode=require` | Your PostgreSQL URL (e.g. from Neon or Render Postgres). Use **single quotes** or “Secret file” if the URL has `&`. |
| `JWT_SECRET` | A long random string (e.g. 32+ chars) | Required. Generate with: `openssl rand -base64 32` |
| `GEMINI_API_KEY` | Your Google AI / Gemini API key | Optional; needed only for AI content generation. |
| `CORS_ORIGINS` | `https://your-frontend.onrender.com` or `*` | Optional. Comma-separated origins. Default `*` allows all. |

- For production, set `CORS_ORIGINS` to your real frontend URL(s), e.g. `https://your-app.onrender.com`.
- Mark **sensitive** values (e.g. `DATABASE_URL`, `JWT_SECRET`, `GEMINI_API_KEY`) as **Secret** so they’re hidden in the dashboard.

---

## Step 5: Deploy

1. Click **Create Web Service** (or **Save** if you’re editing an existing service).
2. Render will clone the repo, run the build command in `services/api`, then start `./api`.
3. Wait for the build and deploy to finish. The **Logs** tab shows build and runtime output.
4. Once live, your API URL will look like:  
   **`https://devflowos-api.onrender.com`** (or the name you chose).

---

## Step 6: Run database migrations

Your app does **not** run migrations on startup. Run them once against the same database you set in `DATABASE_URL`:

**Option A – From your machine (recommended)**

```bash
cd /path/to/devflowos
export DATABASE_URL='postgres://...'   # same URL as in Render
bash infra/db/run_migrations.sh
```

**Option B – From Render Shell (if available on your plan)**

1. In the Render dashboard, open your Web Service.
2. Open **Shell** (if you have it).
3. Run the same migration commands, or set `DATABASE_URL` and run `psql "$DATABASE_URL" -f ...` for each migration file (you’d need to get the migration files into that environment).

Using Option A from your laptop is usually simplest: point `DATABASE_URL` at your production DB and run `run_migrations.sh`.

---

## Step 7: Test the API

- Health-style check: open in browser or with curl:  
  `https://your-service-name.onrender.com/`  
  (Your app may not have a root route; try an actual endpoint below.)
- Signup:  
  `POST https://your-service-name.onrender.com/auth/signup`  
  Body (JSON): `{"email":"you@example.com","password":"yourpassword"}`
- Login:  
  `POST https://your-service-name.onrender.com/auth/login`  
  Body (JSON): `{"email":"you@example.com","password":"yourpassword"}`

Use the returned JWT in the `Authorization: Bearer <token>` header for protected routes.

---

## Summary checklist

- [ ] Repo connected to Render  
- [ ] Web Service created with **Root Directory** = `services/api`  
- [ ] Build: `go build -o api ./cmd/api`  
- [ ] Start: `./api`  
- [ ] `DATABASE_URL` and `JWT_SECRET` set (and optionally `GEMINI_API_KEY`, `CORS_ORIGINS`)  
- [ ] Migrations run against the production database  
- [ ] API tested (signup / login and a protected route)

---

## Troubleshooting

- **Build fails**  
  - Confirm **Root Directory** is `services/api`.  
  - Check the **Logs** tab for the exact Go error.

- **App crashes or “DATABASE_URL required”**  
  - Ensure `DATABASE_URL` is set in **Environment** and that there are no typos.  
  - If the URL contains `&`, store it as a **Secret** or wrap in quotes when pasting.

- **CORS errors from the frontend**  
  - Set `CORS_ORIGINS` to your frontend origin, e.g. `https://your-frontend.onrender.com` (no trailing slash).

- **Free tier spin-down**  
  - On the free tier, the service may sleep after inactivity. The first request after sleep can take 30–60 seconds; subsequent requests are fast.
